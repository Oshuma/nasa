package nasa

import (
	"testing"

	"fmt"
	"time"
)

func TestEncode(t *testing.T) {
	t.Run("APODParams", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			p := APODParams{}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			if out != "" {
				t.Errorf("expected empty string, got: %s", out)
			}
		})

		t.Run("date", func(t *testing.T) {
			d := time.Now()
			p := APODParams{Date: d}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("date=%s", d.Format("2006-01-02"))
			if out != expected {
				t.Errorf("expected: %s, got %s", expected, out)
			}
		})

		t.Run("all set", func(t *testing.T) {
			api_key := "FOO_KEY"
			d := time.Now()

			p := APODParams{
				APIKey: api_key,
				Date:   d,
				HD:     true,
			}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("api_key=%s&date=%s&hd=true", api_key, d.Format("2006-01-02"))
			if out != expected {
				t.Errorf("expected: %s, got %s", expected, out)
			}
		})
	})
}
