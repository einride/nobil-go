package nobil

import (
	"bytes"
	"fmt"
	"time"
)

// A Date represents a date (year, month, day).
//
// This type does not include location information, and therefore does not
// describe a unique 24-hour timespan.
type Date struct {
	Year  int        // Year (e.g., 2014).
	Month time.Month // Month of the year (January = 1, ...).
	Day   int        // Day of the month, starting at 1.
}

// UnmarshalString sets the Date by decoding s.
func (d *Date) UnmarshalString(s string) error {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Year, d.Month, d.Day = t.Date()
	return nil
}

// String returns the date in RFC3339 full-date format.
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// MarshalJSON implements the json.Marshaler interface.
// The output is the result of d.String().
func (d *Date) MarshalJSON() ([]byte, error) {
	var result bytes.Buffer
	result.Grow(12)
	_ = result.WriteByte('"')
	_, _ = result.WriteString(d.String())
	_ = result.WriteByte('"')
	return result.Bytes(), nil
}
