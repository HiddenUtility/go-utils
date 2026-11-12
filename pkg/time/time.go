package time

import "time"

// Now returns the current time in RFC3339 format.
func Now() string {
	return time.Now().Format(time.RFC3339)
}
