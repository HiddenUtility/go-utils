package sql_logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestSqlLogLevel_String(t *testing.T) {
	tests := []struct {
		level    SqlLogLevel
		expected string
	}{
		{SqlLogDebug, "DEBUG"},
		{SqlLogInfo, "INFO"},
		{SqlLogWarning, "WARNING"},
		{SqlLogError, "ERROR"},
		{SqlLogCritical, "CRITICAL"},
		{SqlLogLevel(0), "UNKNOWN"},
		{SqlLogLevel(10), "UNKNOWN"},
	}
	for _, tt := range tests {
		if got := tt.level.String(); got != tt.expected {
			t.Errorf("SqlLogLevel.String() = %v, want %v", got, tt.expected)
		}
	}
}

func TestSqlLogLevel_JSON(t *testing.T) {
	tests := []struct {
		input    string
		expected SqlLogLevel
		wantErr  bool
	}{
		{`"DEBUG"`, SqlLogDebug, false},
		{`"INFO"`, SqlLogInfo, false},
		{`"WARNING"`, SqlLogWarning, false},
		{`"ERROR"`, SqlLogError, false},
		{`"CRITICAL"`, SqlLogCritical, false},
		{`"INVALID"`, 0, true},
		{`123`, 0, true},
	}
	for _, tt := range tests {
		var l SqlLogLevel
		err := json.Unmarshal([]byte(tt.input), &l)
		if (err != nil) != tt.wantErr {
			t.Errorf("UnmarshalJSON(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && l != tt.expected {
			t.Errorf("UnmarshalJSON(%s) = %v, want %v", tt.input, l, tt.expected)
		}

		// Test Marshal
		marshaled, err := json.Marshal(l)
		if err != nil {
			t.Errorf("MarshalJSON(%v) error = %v", l, err)
			continue
		}
		if !tt.wantErr && string(marshaled) != `"`+tt.expected.String()+`"` {
			t.Errorf("MarshalJSON(%v) = %s, want %s", l, string(marshaled), `"`+tt.expected.String()+`"`)
		}
	}
}

func TestSqlLogRecode_Generators(t *testing.T) {
	name := "test-name"
	msg := "test-msg"
	pid := os.Getpid()

	tests := []struct {
		label string
		fn    func(string, string) SqlLogRecode
		level SqlLogLevel
	}{
		{"Debug", Debug, SqlLogDebug},
		{"Info", Info, SqlLogInfo},
		{"Warning", Warning, SqlLogWarning},
		{"Error", Error, SqlLogError},
		{"Critical", Critical, SqlLogCritical},
	}

	for _, tt := range tests {
		rec := tt.fn(name, msg)
		if rec.Name != name {
			t.Errorf("%s() Name = %v, want %v", tt.label, rec.Name, name)
		}
		if rec.Msg != msg {
			t.Errorf("%s() Msg = %v, want %v", tt.label, rec.Msg, msg)
		}
		if rec.Level != tt.level {
			t.Errorf("%s() Level = %v, want %v", tt.label, rec.Level, tt.level)
		}
		if rec.PID != pid {
			t.Errorf("%s() PID = %v, want %v", tt.label, rec.PID, pid)
		}
		if rec.Timestamp <= 0 {
			t.Errorf("%s() Timestamp = %v, want > 0", tt.label, rec.Timestamp)
		}
	}
}

func TestSqlLogRecode_String(t *testing.T) {
	name := "test-component"
	msg := "something happened"
	rec := Info(name, msg)
	str := rec.String()

	// Check if it contains parts of the expected output
	if !strings.Contains(str, "INFO") {
		t.Errorf("String() does not contain level: %s", str)
	}
	if !strings.Contains(str, "["+name+"]") {
		t.Errorf("String() does not contain name: %s", str)
	}
	if !strings.Contains(str, msg) {
		t.Errorf("String() does not contain message: %s", str)
	}
}

func TestSqlLogger_New(t *testing.T) {
	dbPath := "test_logger.db"
	defer os.Remove(dbPath)

	logger, err := NewSqlLogger(dbPath, SqlLogInfo, 100)
	if err != nil {
		t.Fatalf("NewSqlLogger failed: %v", err)
	}
	defer logger.Close()

	if logger.dbPath != dbPath {
		t.Errorf("Expected dbPath %s, got %s", dbPath, logger.dbPath)
	}
	if logger.minLogLevel != SqlLogInfo {
		t.Errorf("Expected minLogLevel %v, got %v", SqlLogInfo, logger.minLogLevel)
	}
	if logger.maxRecNum != 100 {
		t.Errorf("Expected maxRecNum 100, got %d", logger.maxRecNum)
	}
}

func TestSqlLogger_LogFiltering(t *testing.T) {
	dbPath := "test_filter.db"
	defer os.Remove(dbPath)

	logger, err := NewSqlLogger(dbPath, SqlLogInfo, 100)
	if err != nil {
		t.Fatalf("NewSqlLogger failed: %v", err)
	}
	defer logger.Close()

	// Should be allowed
	logger.Log(SqlLogInfo, "test-name", "info message")
	if len(logger.recodes) != 1 {
		t.Errorf("Expected 1 record, got %d", len(logger.recodes))
	}

	// Should be filtered out
	logger.Log(SqlLogDebug, "test-name", "debug message")
	if len(logger.recodes) != 1 {
		t.Errorf("Expected 1 record after filtering debug, got %d", len(logger.recodes))
	}

	// Should be allowed
	logger.Log(SqlLogWarning, "test-name", "warning message")
	if len(logger.recodes) != 2 {
		t.Errorf("Expected 2 records, got %d", len(logger.recodes))
	}
}

func TestSqlLogger_OutAndSearch(t *testing.T) {
	dbPath := "test_out.db"
	defer os.Remove(dbPath)

	logger, err := NewSqlLogger(dbPath, SqlLogDebug, 10)
	if err != nil {
		t.Fatalf("NewSqlLogger failed: %v", err)
	}
	defer logger.Close()

	logger.Log(SqlLogInfo, "app", "msg1")
	logger.Log(SqlLogWarning, "db", "msg2")
	logger.Log(SqlLogDebug, "app", "msg3")

	if err := logger.Out(); err != nil {
		t.Fatalf("Out failed: %v", err)
	}

	if len(logger.recodes) != 0 {
		t.Errorf("Expected 0 recodes in buffer after Out, got %d", len(logger.recodes))
	}

	recs, err := logger.GetRecs()
	if err != nil {
		t.Fatalf("GetRecs failed: %v", err)
	}
	if len(recs) != 3 {
		t.Errorf("Expected 3 records, got %d", len(recs))
	}

	// Test GetRecsWithLimit
	recs, err = logger.GetRecsWithLimit(2)
	if err != nil {
		t.Fatalf("GetRecsWithLimit failed: %v", err)
	}
	if len(recs) != 2 {
		t.Errorf("Expected 2 records, got %d", len(recs))
	}

	// Test GetRecsWith (AND search)
	recs, err = logger.GetRecsWith("app", SqlLogInfo)
	if err != nil {
		t.Fatalf("GetRecsWith failed: %v", err)
	}
	if len(recs) != 1 {
		t.Errorf("Expected 1 record for app/INFO, got %d", len(recs))
	}
	if recs[0].Msg != "msg1" {
		t.Errorf("Expected msg1, got %s", recs[0].Msg)
	}
}

func TestSqlLogger_Reduce(t *testing.T) {
	dbPath := "test_reduce.db"
	defer os.Remove(dbPath)

	logger, err := NewSqlLogger(dbPath, SqlLogDebug, 5)
	if err != nil {
		t.Fatalf("NewSqlLogger failed: %v", err)
	}
	defer logger.Close()

	// Fill up to 10 records in DB
	for i := 0; i < 10; i++ {
		logger.Log(SqlLogInfo, "test", fmt.Sprintf("msg%d", i))
	}
	if err := logger.Out(); err != nil {
		t.Fatalf("Out failed: %v", err)
	}

	// nowRowNum should be 10
	if logger.nowRowNum != 10 {
		t.Errorf("Expected nowRowNum 10, got %d", logger.nowRowNum)
	}

	// Reduce by days (e.g., 0 days should not delete unless nowRowNum > maxRecNum)
	// Wait, the spec says: "nowRowNumがmaxRecNum以下なら削除しない"
	// Currently nowRowNum is 10, maxRecNum is 5. So it should delete.

	err = logger.Reduce(0) // days = 0
	if err != nil {
		t.Fatalf("Reduce failed: %v", err)
	}

	// After Reduce, nowRowNum should be <= 5
	if logger.nowRowNum > 5 {
		t.Errorf("Expected nowRowNum <= 5, got %d", logger.nowRowNum)
	}
}
