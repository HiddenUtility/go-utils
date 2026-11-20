package logger

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLogLevel(t *testing.T) {
	if CRITICAL <= ERROR {
		t.Errorf("CRITICAL should be greater than ERROR")
	}
	if ERROR <= WARNING {
		t.Errorf("ERROR should be greater than WARNING")
	}
	if WARNING <= INFO {
		t.Errorf("WARNING should be greater than INFO")
	}
	if INFO <= DEBUG {
		t.Errorf("INFO should be greater than DEBUG")
	}
}

func TestNewLogRecode(t *testing.T) {
	msg := "test message"
	pid := os.Getpid()

	tests := []struct {
		name     string
		recode   LogRecode
		expected LogLevel
	}{
		{"Critical", NewCritical(msg), CRITICAL},
		{"Error", NewError(msg), ERROR},
		{"Warning", NewWarning(msg), WARNING},
		{"Info", NewInfo(msg), INFO},
		{"Debug", NewDebug(msg), DEBUG},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.recode.Level != tt.expected {
				t.Errorf("expected level %v, got %v", tt.expected, tt.recode.Level)
			}
			if tt.recode.Msg != msg {
				t.Errorf("expected msg %s, got %s", msg, tt.recode.Msg)
			}
			if tt.recode.Pid != pid {
				t.Errorf("expected pid %d, got %d", pid, tt.recode.Pid)
			}
			// Timestamp format check: YYYY-MM-DDTHH:MM:SS.ffffff
			_, err := time.Parse("2006-01-02T15:04:05.000000", tt.recode.Timestamp)
			if err != nil {
				t.Errorf("invalid timestamp format: %v", err)
			}
		})
	}
}

func TestLogRecode_ToLog(t *testing.T) {
	msg := "test message"
	rec := NewInfo(msg)
	got := rec.ToLog()

	if !strings.Contains(got, "[INFO]") {
		t.Errorf("ToLog() = %q; want it to contain [INFO]", got)
	}
	if !strings.Contains(got, msg) {
		t.Errorf("ToLog() = %q; want it to contain %q", got, msg)
	}

	expectedPidPart := fmt.Sprintf("[%d]", rec.Pid)
	if !strings.Contains(got, expectedPidPart) {
		t.Errorf("ToLog() = %q; want it to contain %q", got, expectedPidPart)
	}
}

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := tmpDir + "/app.log"

	lg, err := New(logPath)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	if lg.Path != logPath {
		t.Errorf("expected path %s, got %s", logPath, lg.Path)
	}
	if lg.MaxFileNum != 7 {
		t.Errorf("expected default MaxFileNum 7, got %d", lg.MaxFileNum)
	}
	if lg.MaxFileSize != 200*1024 {
		t.Errorf("expected default MaxFileSize %d, got %d", 200*1024, lg.MaxFileSize)
	}
	if lg.Realtime != false {
		t.Errorf("expected default Realtime false, got %v", lg.Realtime)
	}
	if lg.LogLevel != INFO {
		t.Errorf("expected default LogLevel INFO, got %v", lg.LogLevel)
	}
	if lg.NowSize != 0 {
		t.Errorf("expected default NowSize 0, got %d", lg.NowSize)
	}
}

func TestNewWith(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := tmpDir + "/app.log"

	maxFileNum := uint32(5)
	maxFileSize := uint32(1024)
	realtime := true
	logLevel := DEBUG

	lg, err := NewWith(logPath, maxFileNum, maxFileSize, realtime, logLevel)
	if err != nil {
		t.Fatalf("failed to create logger with values: %v", err)
	}

	if lg.MaxFileNum != maxFileNum {
		t.Errorf("expected MaxFileNum %d, got %d", maxFileNum, lg.MaxFileNum)
	}
	if lg.MaxFileSize != maxFileSize {
		t.Errorf("expected MaxFileSize %d, got %d", maxFileSize, lg.MaxFileSize)
	}
	if lg.Realtime != realtime {
		t.Errorf("expected Realtime %v, got %v", realtime, lg.Realtime)
	}
	if lg.LogLevel != logLevel {
		t.Errorf("expected LogLevel %v, got %v", logLevel, lg.LogLevel)
	}
}

func TestNewWithValidation(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := tmpDir + "/app.log"

	tests := []struct {
		name        string
		maxFileNum  uint32
		maxFileSize uint32
		realtime    bool
		logLevel    LogLevel
		expectErr   bool
	}{
		{"Valid", 1, 1, false, INFO, false},
		{"InvalidMaxFileNum", 0, 1, false, INFO, true},
		{"InvalidMaxFileSize", 1, 0, false, INFO, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWith(logPath, tt.maxFileNum, tt.maxFileSize, tt.realtime, tt.logLevel)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

func TestNewWithFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := tmpDir + "/app.log"

	// Create file with some content
	content := "existing log content"
	err := os.WriteFile(logPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	lg, err := New(logPath)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	if lg.NowSize != uint32(len(content)) {
		t.Errorf("expected NowSize %d, got %d", len(content), lg.NowSize)
	}
}

func TestNewWithCreateDir(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := tmpDir + "/subdir/app.log"

	lg, err := New(logPath)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	if lg.Path != logPath {
		t.Errorf("expected path %s, got %s", logPath, lg.Path)
	}
	// Check if directory was created
	_, err = os.Stat(tmpDir + "/subdir")
	if err != nil {
		t.Errorf("expected directory to be created, but got error: %v", err)
	}
	_ = lg
}

func TestLogger_LoggingMethods(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := tmpDir + "/app.log"

	t.Run("Filtering", func(t *testing.T) {
		lg, _ := NewWith(logPath, 7, 1024, false, INFO)

		_ = lg.Info("info msg")
		if len(lg.Recodes) != 1 {
			t.Errorf("expected 1 record, got %d", len(lg.Recodes))
		}

		_ = lg.Debug("debug msg")
		if len(lg.Recodes) != 1 {
			t.Errorf("expected 1 record (debug filtered), got %d", len(lg.Recodes))
		}

		_ = lg.Critical("critical msg")
		if len(lg.Recodes) != 2 {
			t.Errorf("expected 2 records, got %d", len(lg.Recodes))
		}
	})

	t.Run("Realtime", func(t *testing.T) {
		lg, _ := NewWith(logPath, 7, 1024, true, INFO)
		_ = lg.Info("realtime msg")
		// Since Out() is a stub, it won't do much, but it shouldn't fail.
	})
}

func TestWriteRecs(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := tmpDir + "/app.log"

	t.Run("Basic", func(t *testing.T) {
		lg, _ := NewWith(logPath, 7, 1024, false, INFO)
		recs := []LogRecode{
			NewInfo("info 1"),
			NewDebug("debug 1"), // filtered
			NewError("error 1"),
		}
		err := lg.WriteRecs(recs)
		if err != nil {
			t.Fatalf("WriteRecs failed: %v", err)
		}
		if len(lg.Recodes) != 2 {
			t.Errorf("expected 2 records, got %d", len(lg.Recodes))
		}
	})

	t.Run("Empty", func(t *testing.T) {
		lg, _ := NewWith(logPath, 7, 1024, false, INFO)
		err := lg.WriteRecs([]LogRecode{})
		if err != nil {
			t.Fatalf("WriteRecs failed on empty: %v", err)
		}
		if len(lg.Recodes) != 0 {
			t.Errorf("expected 0 records, got %d", len(lg.Recodes))
		}
	})

	t.Run("AllFiltered", func(t *testing.T) {
		lg, _ := NewWith(logPath, 7, 1024, false, ERROR)
		recs := []LogRecode{
			NewInfo("info 1"),
			NewWarning("warning 1"),
		}
		err := lg.WriteRecs(recs)
		if err != nil {
			t.Fatalf("WriteRecs failed on all filtered: %v", err)
		}
		if len(lg.Recodes) != 0 {
			t.Errorf("expected 0 records, got %d", len(lg.Recodes))
		}
	})
}

func TestLogger_Out(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := tmpDir + "/app.log"

	t.Run("BufferOutput", func(t *testing.T) {
		lg, _ := New(logPath)
		_ = lg.Info("msg 1")
		_ = lg.Info("msg 2")

		if len(lg.Recodes) != 2 {
			t.Errorf("expected 2 records in buffer, got %d", len(lg.Recodes))
		}

		err := lg.Out()
		if err != nil {
			t.Fatalf("Out failed: %v", err)
		}

		if len(lg.Recodes) != 0 {
			t.Errorf("expected 0 records in buffer after Out, got %d", len(lg.Recodes))
		}

		content, _ := os.ReadFile(logPath)
		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		if len(lines) != 2 {
			t.Errorf("expected 2 lines in file, got %d", len(lines))
		}
		if !strings.Contains(lines[0], "msg 1") || !strings.Contains(lines[1], "msg 2") {
			t.Errorf("unexpected file content: %q", string(content))
		}
	})

	t.Run("EmptyBuffer", func(t *testing.T) {
		lg, _ := New(logPath)
		err := lg.Out()
		if err != nil {
			t.Errorf("Out failed on empty buffer: %v", err)
		}
	})
}

func TestLogger_Rotation(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := tmpDir + "/app.log"

	t.Run("RotationTrigger", func(t *testing.T) {
		// MaxFileSize = 100 bytes
		lg, _ := NewWith(logPath, 3, 100, false, INFO)

		msg := "triggering rotation with long message" // ~39 bytes
		// 3 writes should exceed 100.
		_ = lg.Info(msg)
		_ = lg.Out()
		_ = lg.Info(msg)
		_ = lg.Out()
		_ = lg.Info(msg)
		_ = lg.Out()

		// Check files
		// app.log, app.log.1 should exist
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			t.Errorf("app.log should exist")
		}
		if _, err := os.Stat(logPath + ".1"); os.IsNotExist(err) {
			t.Errorf("app.log.1 should exist")
		}
	})

	t.Run("MaxFileNumLimit", func(t *testing.T) {
		// MaxFileNum = 2. Files: app.log, app.log.1
		lg, _ := NewWith(logPath, 2, 50, false, INFO)

		msg := "rotation test" // ~13 bytes. 14 with newline.
		for i := 0; i < 10; i++ {
			_ = lg.Info(msg)
			_ = lg.Out()
		}

		// Only app.log and app.log.1 should exist. app.log.2 should NOT exist.
		if _, err := os.Stat(logPath + ".2"); !os.IsNotExist(err) {
			t.Errorf("app.log.2 should not exist")
		}
		if _, err := os.Stat(logPath + ".1"); os.IsNotExist(err) {
			t.Errorf("app.log.1 should exist")
		}
	})

	t.Run("MaxFileNumOne", func(t *testing.T) {
		// MaxFileNum = 1. No rotation files.
		tmpDir := t.TempDir()
		logPath := tmpDir + "/app.log"
		lg, _ := NewWith(logPath, 1, 50, false, INFO)

		msg := "rotation test"
		_ = lg.Info(msg)
		_ = lg.Out()
		_ = lg.Info(msg)
		_ = lg.Out()

		if _, err := os.Stat(logPath + ".1"); !os.IsNotExist(err) {
			t.Errorf("app.log.1 should not exist when MaxFileNum=1")
		}
	})
}

func TestLogger_Close(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := tmpDir + "/app.log"

	t.Run("FlushOnClose", func(t *testing.T) {
		lg, _ := New(logPath)
		_ = lg.Info("flush me")

		err := lg.Close()
		if err != nil {
			t.Fatalf("Close failed: %v", err)
		}

		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), "flush me") {
			t.Errorf("expected 'flush me' in file, got %q", string(content))
		}
	})

	t.Run("MultipleClose", func(t *testing.T) {
		lg, _ := New(logPath)
		_ = lg.Info("first")
		_ = lg.Close()
		err := lg.Close()
		if err != nil {
			t.Errorf("second Close failed: %v", err)
		}
	})
}
