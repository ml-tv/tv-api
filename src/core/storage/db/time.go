package db

import (
	"database/sql/driver"
	"time"
)

// ISO8601 is a time.Time layout for the ISO8601 format
const ISO8601 = "2006-01-02T15:04:05-0700"

// Time represents a time.Time that uses ISO8601 for json input/output
// instead of RFC3339
type Time struct {
	time.Time
}

// Now returns the current local time.
func Now() *Time {
	return &Time{Time: time.Now()}
}

// Value - Implementation of valuer for database/sql
func (t *Time) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}

	return t.Format(ISO8601), nil
}

// Scan - Implement the database/sql scanner interface
func (t *Time) Scan(value interface{}) error {
	if value != nil {
		t.Time = value.(time.Time)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(ISO8601) + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	t.Time, err = time.Parse(`"`+ISO8601+`"`, string(data))
	return
}
