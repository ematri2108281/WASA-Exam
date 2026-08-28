// Package globaltime provides a replaceable clock for the application.
//
// Using Now() and Since() instead of time.Now() / time.Since() throughout
// the codebase allows tests to inject a deterministic fixed time via FixedTime,
// making time-dependent logic fully reproducible without real sleeps.
package globaltime

import "time"

const isoLayout = time.RFC3339

// FixedTime, when set to a non-zero value, is returned by Now() instead of
// the real wall clock. Useful in unit tests.
var FixedTime time.Time

// Now returns the current wall-clock time, or FixedTime if it has been set.
func Now() time.Time {
	if !FixedTime.IsZero() {
		return FixedTime
	}
	return time.Now()
}

// Since returns the elapsed duration since tm, using the overridable Now().
func Since(tm time.Time) time.Duration {
	return Now().Sub(tm)
}

// FormatISO formats t as an ISO 8601 / RFC 3339 string (e.g. "2006-01-02T15:04:05Z07:00").
func FormatISO(t time.Time) string {
	return t.Format(isoLayout)
}

// ParseISO parses an ISO 8601 / RFC 3339 string and returns the corresponding time.Time.
func ParseISO(s string) (time.Time, error) {
	return time.Parse(isoLayout, s)
}
