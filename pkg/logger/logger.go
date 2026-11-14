package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota + 1
	INFO
	WARNING
	ERROR
	CRITICAL
)

type LogRecode struct {
	Timestamp string
	Level     LogLevel
	Msg       string
	Pid       int
}

func NewCritical(msg string) LogRecode {
	return newRecode(CRITICAL, msg)
}

func NewError(msg string) LogRecode {
	return newRecode(ERROR, msg)
}

func NewWarning(msg string) LogRecode {
	return newRecode(WARNING, msg)
}

func NewInfo(msg string) LogRecode {
	return newRecode(INFO, msg)
}

func NewDebug(msg string) LogRecode {
	return newRecode(DEBUG, msg)
}

func newRecode(level LogLevel, msg string) LogRecode {
	return LogRecode{
		Timestamp: time.Now().Format("2006-01-02T15:04:05.000000"),
		Level:     level,
		Msg:       msg,
		Pid:       os.Getpid(),
	}
}

func (r LogRecode) ToLog() string {
	levelStr := ""
	switch r.Level {
	case CRITICAL:
		levelStr = "CRITICAL"
	case ERROR:
		levelStr = "ERROR"
	case WARNING:
		levelStr = "WARNING"
	case INFO:
		levelStr = "INFO"
	case DEBUG:
		levelStr = "DEBUG"
	}

	return fmt.Sprintf("%s [%s] [%d] %s", r.Timestamp, levelStr, r.Pid, r.Msg)
}

type Logger struct {
	Path        string
	Recodes     []LogRecode
	MaxFileNum  uint32
	MaxFileSize uint32
	Realtime    bool
	LogLevel    LogLevel
	NowSize     uint32
}

func New(path string) (*Logger, error) {
	return NewWith(path, 7, 200*1024, false, INFO)
}

func NewWith(path string, maxFileNum uint32, maxFileSize uint32, realtime bool, logLevel LogLevel) (*Logger, error) {
	if maxFileNum < 1 {
		return nil, fmt.Errorf("maxFileNum must be at least 1")
	}
	if maxFileSize < 1 {
		return nil, fmt.Errorf("maxFileSize must be at least 1")
	}

	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	var nowSize uint32
	info, err := os.Stat(path)
	if err == nil {
		nowSize = uint32(info.Size())
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	return &Logger{
		Path:        path,
		MaxFileNum:  maxFileNum,
		MaxFileSize: maxFileSize,
		Realtime:    realtime,
		LogLevel:    logLevel,
		NowSize:     nowSize,
		Recodes:     make([]LogRecode, 0),
	}, nil
}

func (l *Logger) Critical(msg string) error {
	return l.addRecord(NewCritical(msg))
}

func (l *Logger) Error(msg string) error {
	return l.addRecord(NewError(msg))
}

func (l *Logger) Warning(msg string) error {
	return l.addRecord(NewWarning(msg))
}

func (l *Logger) Info(msg string) error {
	return l.addRecord(NewInfo(msg))
}

func (l *Logger) Debug(msg string) error {
	return l.addRecord(NewDebug(msg))
}

func (l *Logger) addRecord(rec LogRecode) error {
	if rec.Level >= l.LogLevel {
		l.Recodes = append(l.Recodes, rec)
		if l.Realtime {
			return l.Out()
		}
	}
	return nil
}

func (l *Logger) WriteRecs(recs []LogRecode) error {
	if len(recs) == 0 {
		return nil
	}

	addedCount := 0
	for _, rec := range recs {
		if rec.Level >= l.LogLevel {
			l.Recodes = append(l.Recodes, rec)
			addedCount++
		}
	}

	if addedCount > 0 && l.Realtime {
		return l.Out()
	}

	return nil
}

func (l *Logger) Out() error {
	if len(l.Recodes) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, rec := range l.Recodes {
		sb.WriteString(rec.ToLog())
		sb.WriteString("\n")
	}
	logData := sb.String()
	dataLen := uint32(len(logData))

	// Open file in append mode, create if not exists
	f, err := os.OpenFile(l.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(logData); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	l.NowSize += dataLen
	l.Recodes = nil // Clear buffer after successful write

	if l.NowSize > l.MaxFileSize {
		if err := l.rotate(); err != nil {
			return fmt.Errorf("failed to rotate log file: %w", err)
		}
	}

	return nil
}

func (l *Logger) Close() error {
	if len(l.Recodes) == 0 {
		return nil
	}
	return l.Out()
}

func (l *Logger) rotate() error {
	if l.MaxFileNum <= 1 {
		if err := os.Remove(l.Path); err != nil && !os.IsNotExist(err) {
			return err
		}
		l.NowSize = 0
		return nil
	}

	oldest := fmt.Sprintf("%s.%d", l.Path, l.MaxFileNum-1)
	if err := os.Remove(oldest); err != nil && !os.IsNotExist(err) {
		return err
	}

	for i := int(l.MaxFileNum - 2); i >= 1; i-- {
		oldName := fmt.Sprintf("%s.%d", l.Path, i)
		newName := fmt.Sprintf("%s.%d", l.Path, i+1)
		if err := os.Rename(oldName, newName); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	if err := os.Rename(l.Path, l.Path+".1"); err != nil {
		return err
	}

	l.NowSize = 0
	f, err := os.Create(l.Path)
	if err != nil {
		return err
	}
	f.Close()

	return nil
}
