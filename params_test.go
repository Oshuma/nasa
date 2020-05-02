package nasa

import (
	"testing"

	"fmt"
	"time"
)

func TestEncode(t *testing.T) {
	apiKey := "NASA_KEY"

	t.Run("APODParams", func(t *testing.T) {
		t.Run("no APIKey", func(t *testing.T) {
			p := &APODParams{}

			_, err := p.Encode()
			if err != ErrorNoAPIKey {
				t.Errorf("wrong error returned: %s", err)
			}
		})

		t.Run("APIKey only", func(t *testing.T) {
			p := &APODParams{APIKey: apiKey}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("api_key=%s", apiKey)
			if out != expected {
				t.Errorf("expected: %s, got: %s", expected, out)
			}
		})

		t.Run("date", func(t *testing.T) {
			d := time.Now()
			p := &APODParams{
				APIKey: apiKey,
				Date:   d,
			}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("api_key=%s&date=%s", apiKey, d.Format("2006-01-02"))
			if out != expected {
				t.Errorf("expected: %s, got %s", expected, out)
			}
		})

		t.Run("all set", func(t *testing.T) {
			d := time.Now()

			p := &APODParams{
				APIKey: apiKey,
				Date:   d,
				HD:     true,
			}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("api_key=%s&date=%s&hd=true", apiKey, d.Format("2006-01-02"))
			if out != expected {
				t.Errorf("expected: %s, got %s", expected, out)
			}
		})
	})

	t.Run("EPICParams", func(t *testing.T) {
		t.Run("no APIKey", func(t *testing.T) {
			p := EPICParams{}

			_, err := p.Encode()
			if err != ErrorNoAPIKey {
				t.Errorf("wrong error returned: %s", err)
			}
		})

		t.Run("date", func(t *testing.T) {
			d := time.Now()
			p := &EPICParams{
				APIKey: apiKey,
				Date:   d,
			}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("api/natural/date/%s?api_key=%s", p.Date.Format("2006-01-02"), apiKey)
			if out != expected {
				t.Errorf("expected: %s, got: %s", expected, out)
			}
		})

		t.Run("no date", func(t *testing.T) {
			p := &EPICParams{APIKey: apiKey}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("api/natural?api_key=%s", apiKey)
			if out != expected {
				t.Errorf("expected: %s, got: %s", expected, out)
			}
		})
	})

	t.Run("MarsPhotosParams", func(t *testing.T) {
		t.Run("no APIKey", func(t *testing.T) {
			p := &MarsPhotosParams{}

			_, err := p.Encode()
			if err != ErrorNoAPIKey {
				t.Errorf("wrong error returned: %s", err)
			}
		})

		t.Run("defaults", func(t *testing.T) {
			p := &MarsPhotosParams{APIKey: apiKey}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("api_key=%s&sol=0", apiKey)
			if out != expected {
				t.Errorf("\nexpected: %s\ngot: %s", expected, out)
			}
		})

		t.Run("EarthDate", func(t *testing.T) {
			d := time.Now()
			p := &MarsPhotosParams{APIKey: apiKey, EarthDate: d}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			earthDate := d.Format("2006-01-02")
			expected := fmt.Sprintf("api_key=%s&earth_date=%s", apiKey, earthDate)
			if out != expected {
				t.Errorf("\nexpected: %s\ngot: %s", expected, out)
			}
		})

		t.Run("Camera", func(t *testing.T) {
			cam := RoverCameraFHAZ
			p := &MarsPhotosParams{APIKey: apiKey, Camera: cam}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("api_key=%s&camera=%s&sol=0", apiKey, cam.Slug)
			if out != expected {
				t.Errorf("\nexpected: %s\ngot: %s", expected, out)
			}
		})
	})
}

func TestGetAPIKey(t *testing.T) {
	apiKey := "NASA_KEY"

	t.Run("APODParams", func(t *testing.T) {
		p := &APODParams{APIKey: apiKey}
		out := p.GetAPIKey()
		if out != apiKey {
			t.Errorf("expected: %s, got: %s", apiKey, out)
		}
	})

	t.Run("EPICParams", func(t *testing.T) {
		p := &EPICParams{APIKey: apiKey}
		out := p.GetAPIKey()
		if out != apiKey {
			t.Errorf("expected: %s, got: %s", apiKey, out)
		}
	})

	t.Run("MarsPhotosParams", func(t *testing.T) {
		p := &MarsPhotosParams{APIKey: apiKey}
		out := p.GetAPIKey()
		if out != apiKey {
			t.Errorf("expected: %s, got: %s", apiKey, out)
		}
	})
}
