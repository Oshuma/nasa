package nasa

import (
	"strings"
	"time"
)

type APODDate struct {
	time.Time
}

func (d *APODDate) UnmarshalJSON(b []byte) error {
	t, err := parseTime(b, "2006-01-02")
	if err != nil {
		return err
	}
	*d = APODDate{Time: t}
	return nil
}

type EPICDate struct {
	time.Time
}

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
