package sql_logger

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type SqlLogLevel int

const (
	SqlLogDebug    SqlLogLevel = 1
	SqlLogInfo     SqlLogLevel = 2
	SqlLogWarning  SqlLogLevel = 3
	SqlLogError    SqlLogLevel = 4
	SqlLogCritical SqlLogLevel = 5
)

func (l SqlLogLevel) String() string {
	switch l {
	case SqlLogDebug:
		return "DEBUG"
	case SqlLogInfo:
		return "INFO"
	case SqlLogWarning:
		return "WARNING"
	case SqlLogError:
		return "ERROR"
	case SqlLogCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

func (l SqlLogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

func (l *SqlLogLevel) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "DEBUG":
		*l = SqlLogDebug
	case "INFO":
		*l = SqlLogInfo
	case "WARNING":
		*l = SqlLogWarning
	case "ERROR":
		*l = SqlLogError
	case "CRITICAL":
		*l = SqlLogCritical
	default:
		return fmt.Errorf("invalid log level: %s", s)
	}
	return nil
}

type SqlLogRecode struct {
	Timestamp float64     `json:"timestamp"`
	Level     SqlLogLevel `json:"level"`
	Name      string      `json:"name"`
	Msg       string      `json:"msg"`
	PID       int         `json:"pid"`
}

func (r SqlLogRecode) String() string {
	t := time.Unix(int64(r.Timestamp), int64((r.Timestamp-float64(int64(r.Timestamp)))*1e9))
	return fmt.Sprintf("%s %s [%s] %s", t.Local().Format("2006-01-02 15:04:05"), r.Level.String(), r.Name, r.Msg)
}

func Critical(name string, msg string) SqlLogRecode {
	return createRecode(SqlLogCritical, name, msg)
}

func Error(name string, msg string) SqlLogRecode {
	return createRecode(SqlLogError, name, msg)
}

func Warning(name string, msg string) SqlLogRecode {
	return createRecode(SqlLogWarning, name, msg)
}

func Info(name string, msg string) SqlLogRecode {
	return createRecode(SqlLogInfo, name, msg)
}

func Debug(name string, msg string) SqlLogRecode {
	return createRecode(SqlLogDebug, name, msg)
}

func createRecode(level SqlLogLevel, name string, msg string) SqlLogRecode {
	now := time.Now()
	return SqlLogRecode{
		Timestamp: float64(now.UnixNano()) / 1e9,
		Level:     level,
		Name:      name,
		Msg:       msg,
		PID:       os.Getpid(),
	}
}

type SqlLogger struct {
	db          *sql.DB
	mu          sync.Mutex
	maxRecNum   int
	nowRowNum   int
	recodes     []SqlLogRecode
	minLogLevel SqlLogLevel
	dbPath      string
}

func NewSqlLogger(dbPath string, minLogLevel SqlLogLevel, maxRecNum int) (*SqlLogger, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	l := &SqlLogger{
		db:          db,
		maxRecNum:   maxRecNum,
		minLogLevel: minLogLevel,
		dbPath:      dbPath,
		recodes:     make([]SqlLogRecode, 0, maxRecNum),
	}

	if err := l.Initialize(); err != nil {
		db.Close()
		return nil, err
	}

	return l, nil
}

func (l *SqlLogger) Initialize() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, err := l.db.Exec("DROP TABLE IF EXISTS recodes"); err != nil {
		return err
	}

	query := `
	CREATE TABLE recodes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp REAL,
		level INTEGER,
		name TEXT,
		msg TEXT,
		pid INTEGER
	);
	CREATE INDEX idx_recodes_timestamp ON recodes(timestamp);
	CREATE INDEX idx_recodes_name ON recodes(name);
	CREATE INDEX idx_recodes_level ON recodes(level);
	`
	if _, err := l.db.Exec(query); err != nil {
		return err
	}

	l.nowRowNum = 0
	l.recodes = l.recodes[:0]
	return nil
}

func (l *SqlLogger) Log(level SqlLogLevel, name string, msg string) {
	if level < l.minLogLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.recodes = append(l.recodes, createRecode(level, name, msg))
	l.nowRowNum++
}

func (l *SqlLogger) Out() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.recodes) == 0 {
		return nil
	}

	tx, err := l.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO recodes (timestamp, level, name, msg, pid) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, rec := range l.recodes {
		_, err := stmt.Exec(rec.Timestamp, int(rec.Level), rec.Name, rec.Msg, rec.PID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	l.recodes = l.recodes[:0]
	return nil
}

func (l *SqlLogger) Close() error {
	_ = l.Out()
	return l.db.Close()
}

func (l *SqlLogger) Reduce(days int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.nowRowNum <= l.maxRecNum && days <= 0 {
		return nil
	}

	if days > 0 {
		cutoff := time.Now().AddDate(0, 0, -days).UnixNano() / 1e9
		_, err := l.db.Exec("DELETE FROM recodes WHERE timestamp < ?", cutoff)
		if err != nil {
			return err
		}
	}

	var dbCount int
	err := l.db.QueryRow("SELECT COUNT(*) FROM recodes").Scan(&dbCount)
	if err != nil {
		return err
	}

	l.nowRowNum = dbCount + len(l.recodes)

	if l.nowRowNum > l.maxRecNum {
		keepInDb := l.maxRecNum - len(l.recodes)
		if keepInDb < 0 {
			keepInDb = 0
		}

		if keepInDb == 0 {
			_, err = l.db.Exec("DELETE FROM recodes")
		} else {
			_, err = l.db.Exec(`
				DELETE FROM recodes WHERE id NOT IN (
					SELECT id FROM recodes ORDER BY timestamp DESC LIMIT ?
				)`, keepInDb)
		}
		if err != nil {
			return err
		}

		err = l.db.QueryRow("SELECT COUNT(*) FROM recodes").Scan(&dbCount)
		if err != nil {
			return err
		}
		l.nowRowNum = dbCount + len(l.recodes)
	}

	return nil
}

func (l *SqlLogger) GetRecs() ([]SqlLogRecode, error) {
	return l.GetRecsWithLimit(1000)
}

func (l *SqlLogger) GetRecsWithLimit(limit int) ([]SqlLogRecode, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	rows, err := l.db.Query("SELECT timestamp, level, name, msg, pid FROM recodes ORDER BY timestamp DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbRecs []SqlLogRecode
	for rows.Next() {
		var r SqlLogRecode
		var levelInt int
		if err := rows.Scan(&r.Timestamp, &levelInt, &r.Name, &r.Msg, &r.PID); err != nil {
			return nil, err
		}
		r.Level = SqlLogLevel(levelInt)
		dbRecs = append(dbRecs, r)
	}

	all := append([]SqlLogRecode{}, l.recodes...)
	all = append(all, dbRecs...)

	sort.Slice(all, func(i, j int) bool {
		return all[i].Timestamp > all[j].Timestamp
	})

	if len(all) > limit {
		all = all[:limit]
	}

	return all, nil
}

func (l *SqlLogger) GetRecsWith(name string, level SqlLogLevel) ([]SqlLogRecode, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	rows, err := l.db.Query("SELECT timestamp, level, name, msg, pid FROM recodes WHERE name = ? AND level = ? ORDER BY timestamp DESC", name, int(level))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbRecs []SqlLogRecode
	for rows.Next() {
		var r SqlLogRecode
		var levelInt int
		if err := rows.Scan(&r.Timestamp, &levelInt, &r.Name, &r.Msg, &r.PID); err != nil {
			return nil, err
		}
		r.Level = SqlLogLevel(levelInt)
		dbRecs = append(dbRecs, r)
	}

	all := append([]SqlLogRecode{}, l.recodes...)
	all = append(all, dbRecs...)

	var filtered []SqlLogRecode
	for _, r := range all {
		matchName := (name == "" || r.Name == name)
		matchLevel := (int(r.Level) == int(level))
		if matchName && matchLevel {
			filtered = append(filtered, r)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp > filtered[j].Timestamp
	})

	return filtered, nil
}
