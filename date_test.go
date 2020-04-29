package nasa

import (
	"testing"

	"time"
)

func TestUnmarshalJSON(t *testing.T) {
	in := []byte(`""2020-04-29""`)
	d := &Date{}
	err := d.UnmarshalJSON(in)
	if err != nil {
		t.Error(err)
	}

	if d.IsZero() {
		t.Error("date not properly parsed")
	}

	expected := time.Date(2020, 4, 29, 0, 0, 0, 0, time.UTC)
	if !d.Equal(expected) {
		t.Errorf("expected: %s, got: %s", expected, d)
	}
}
