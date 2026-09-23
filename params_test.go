package nasa

import (
	"testing"

	"fmt"
	"net/url"
	"time"
)

func TestEncode(t *testing.T) {
	apiKey := "NASA_KEY"

	t.Run("APIParam", func(t *testing.T) {
		t.Run("no APIKey", func(t *testing.T) {
			p := &APIParam{}

			_, err := p.Encode()
			if err != ErrorNoAPIKey {
				t.Errorf("wrong error returned: %s", err)
			}
		})

		t.Run("encoded", func(t *testing.T) {
			p := APIParam{APIKey: apiKey}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("api_key=%s", apiKey)
			if out != expected {
				t.Errorf("expected: %s, got: %s", expected, out)
			}
		})
	})

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
				Thumbs: true,
			}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("api_key=%s&date=%s&hd=true&thumbs=true", apiKey, d.Format("2006-01-02"))
			if out != expected {
				t.Errorf("expected: %s, got %s", expected, out)
			}
		})

		t.Run("date range", func(t *testing.T) {
			start := time.Now().AddDate(0, 0, -5)
			end := time.Now()
			p := &APODParams{
				APIKey:    apiKey,
				StartDate: start,
				EndDate:   end,
			}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("api_key=%s&end_date=%s&start_date=%s", apiKey, end.Format("2006-01-02"), start.Format("2006-01-02"))
			if out != expected {
				t.Errorf("expected: %s, got %s", expected, out)
			}
		})

		t.Run("count", func(t *testing.T) {
			p := &APODParams{
				APIKey: apiKey,
				Count:  3,
			}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("api_key=%s&count=3", apiKey)
			if out != expected {
				t.Errorf("expected: %s, got %s", expected, out)
			}
		})
	})

	t.Run("EPICParams", func(t *testing.T) {
		t.Run("defaults", func(t *testing.T) {
			p := EPICParams{}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := "api/natural"
			if out != expected {
				t.Errorf("expected: %s, got: %s", expected, out)
			}
		})

		t.Run("date with API key", func(t *testing.T) {
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

		t.Run("escapes API key", func(t *testing.T) {
			p := &EPICParams{APIKey: "a&b=c"}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := "api/natural?api_key=a%26b%3Dc"
			if out != expected {
				t.Errorf("expected: %s, got: %s", expected, out)
			}
		})

		t.Run("custom collection", func(t *testing.T) {
			p := &EPICParams{Collection: "enhanced"}

			out, err := p.Encode()
			if err != nil {
				t.Error(err)
			}

			expected := "api/enhanced"
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

	t.Run("MediaParams", func(t *testing.T) {
		t.Run("require search parameter", func(t *testing.T) {
			p := &MediaParams{}

			_, err := p.Encode()
			if err != ErrorNoQuery {
				t.Errorf("wrong error returned: %s", err)
			}
		})

		t.Run("allows non-query search parameters", func(t *testing.T) {
			p := &MediaParams{NasaID: "as11"}

			out, err := p.Encode()
			if err != nil {
				t.Fatal(err)
			}

			values, err := url.ParseQuery(out)
			if err != nil {
				t.Fatal(err)
			}

			if got := values.Get("nasa_id"); got != "as11" {
				t.Errorf("expected nasa_id to be encoded, got: %s", got)
			}
		})

		t.Run("encodes all supported parameters", func(t *testing.T) {
			p := &MediaParams{
				Query:            "apollo 11",
				Center:           "JSC",
				Description:      "moon landing",
				Description508:   "accessible",
				Keywords:         "apollo,moon",
				Location:         "Moon",
				MediaType:        "image,video",
				NasaID:           "as11-40-5874",
				Page:             3,
				PageSize:         50,
				Photographer:     "Neil Armstrong",
				SecondaryCreator: "Buzz Aldrin",
				Title:            "Apollo 11",
				YearStart:        "1968",
				YearEnd:          "1969",
			}

			out, err := p.Encode()
			if err != nil {
				t.Fatal(err)
			}

			values, err := url.ParseQuery(out)
			if err != nil {
				t.Fatal(err)
			}

			checks := map[string]string{
				"q":                 "apollo 11",
				"center":            "JSC",
				"description":       "moon landing",
				"description_508":   "accessible",
				"keywords":          "apollo,moon",
				"location":          "Moon",
				"media_type":        "image,video",
				"nasa_id":           "as11-40-5874",
				"page":              "3",
				"page_size":         "50",
				"photographer":      "Neil Armstrong",
				"secondary_creator": "Buzz Aldrin",
				"title":             "Apollo 11",
				"year_start":        "1968",
				"year_end":          "1969",
			}

			for key, want := range checks {
				if got := values.Get(key); got != want {
					t.Errorf("expected %s=%s, got %s", key, want, got)
				}
			}
		})
	})
}
