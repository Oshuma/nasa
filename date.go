package nasa

import (
	"strings"
	"time"
)

// Date is a time.Time wrapper.
type Date struct {
	time.Time
}

// UnmarshalJSON unmarshals a date formatted as YYYY-MM-DD.
func (d *Date) UnmarshalJSON(b []byte) error {
	t, err := parseTime(b, "2006-01-02")
	if err != nil {
		return err
	}
	*d = Date{Time: t}
	return nil
}

// EPICDate is a time.Time wrapper used to parse the EPIC date response.
type EPICDate struct {
	time.Time
}

// UnmarshalJSON unmarshals a date formatted as YYYY-MM-DD HH-MM-SS.
func (d *EPICDate) UnmarshalJSON(b []byte) error {
	t, err := parseTime(b, "2006-01-02 15:04:05")
	if err != nil {
		return err
	}
	*d = EPICDate{Time: t}
	return nil
}

func parseTime(b []byte, format string) (time.Time, error) {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(format, s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
