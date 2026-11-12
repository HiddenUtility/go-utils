package time

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	res := Now()
	_, err := time.Parse(time.RFC3339, res)
	if err != nil {
		t.Errorf("Now() returned invalid RFC3339 format: %v, res: %s", err, res)
	}
}
