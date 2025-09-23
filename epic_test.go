package nasa

import (
	"strings"
	"testing"

	"time"
)

func TestBuildNaturalURLs(t *testing.T) {
	full := "https://epic.gsfc.nasa.gov/archive/natural/2020/04/24/png/epic_1b_20200424002712.png"
	thumb := "https://epic.gsfc.nasa.gov/archive/natural/2020/04/24/thumbs/epic_1b_20200424002712.jpg"

	p := &EPICParams{}
	e := &EPICImage{
		Image: "epic_1b_20200424002712",
		Date:  EPICDate{Time: time.Date(2020, 4, 24, 0, 0, 0, 0, time.UTC)},
	}

	e.buildNaturalURLs(p)

	if e.URL.Natural != full {
		t.Errorf("\nexpected: %s\ngot: %s", full, e.URL.Natural)
	}

	if e.URL.Thumb.Natural != thumb {
		t.Errorf("\nexpected: %s\ngot: %s", thumb, e.URL.Thumb.Natural)
	}
}

func TestBuildEnhancedURLs(t *testing.T) {
	full := "https://epic.gsfc.nasa.gov/archive/enhanced/2020/04/24/png/epic_RGB_20200424002712.png"
	thumb := "https://epic.gsfc.nasa.gov/archive/enhanced/2020/04/24/thumbs/epic_RGB_20200424002712.jpg"

	p := &EPICParams{}
	e := &EPICImage{
		Image: "epic_1b_20200424002712",
		Date:  EPICDate{Time: time.Date(2020, 4, 24, 0, 0, 0, 0, time.UTC)},
	}

	e.buildEnhancedURLs(p)

	if e.URL.Enhanced != full {
		t.Errorf("\nexpected: %s\ngot: %s", full, e.URL.Enhanced)
	}

	if e.URL.Thumb.Enhanced != thumb {
		t.Errorf("\nexpected: %s\ngot: %s", thumb, e.URL.Thumb.Enhanced)
	}
}

func TestBuildURLsWithAPIKey(t *testing.T) {
	p := &EPICParams{APIKey: "NASA_KEY"}
	e := &EPICImage{
		Image: "epic_1b_20200424002712",
		Date:  EPICDate{Time: time.Date(2020, 4, 24, 0, 0, 0, 0, time.UTC)},
	}

	e.buildNaturalURLs(p)

	if !strings.HasSuffix(e.URL.Natural, "?api_key=NASA_KEY") {
		t.Errorf("expected api key suffix, got: %s", e.URL.Natural)
	}

	if !strings.HasSuffix(e.URL.Thumb.Natural, "?api_key=NASA_KEY") {
		t.Errorf("expected api key suffix, got: %s", e.URL.Thumb.Natural)
	}
}

func TestBuildURLsFromEnhancedImage(t *testing.T) {
	fullNatural := "https://epic.gsfc.nasa.gov/archive/natural/2020/04/24/png/epic_1b_20200424002712.png"
	fullEnhanced := "https://epic.gsfc.nasa.gov/archive/enhanced/2020/04/24/png/epic_RGB_20200424002712.png"

	img := &EPICImage{
		Identifier: "20200424002712",
		Image:      "epic_RGB_20200424002712",
		Date:       EPICDate{Time: time.Date(2020, 4, 24, 0, 0, 0, 0, time.UTC)},
	}

	img.buildNaturalURLs(&EPICParams{})
	img.buildEnhancedURLs(&EPICParams{})

	if img.URL.Natural != fullNatural {
		t.Errorf("expected natural URL %s, got %s", fullNatural, img.URL.Natural)
	}

	if img.URL.Enhanced != fullEnhanced {
		t.Errorf("expected enhanced URL %s, got %s", fullEnhanced, img.URL.Enhanced)
	}
}

func TestEPICImagesMostRecent(t *testing.T) {
	dates := []time.Time{
		time.Date(2020, 4, 24, 1, 0, 0, 0, time.UTC),
		time.Date(2020, 4, 24, 3, 0, 0, 0, time.UTC),
		time.Date(2020, 4, 25, 0, 0, 0, 0, time.UTC),
	}

	images := EPICImages{
		&EPICImage{Image: "one", Date: EPICDate{Time: dates[0]}},
		&EPICImage{Image: "two", Date: EPICDate{Time: dates[1]}},
		&EPICImage{Image: "three", Date: EPICDate{Time: dates[2]}},
	}

	latest := images.mostRecent()

	if len(latest) != 1 {
		t.Fatalf("expected a single image, got %d", len(latest))
	}

	if latest[0].Image != "three" {
		t.Errorf("expected most recent image to be 'three', got %q", latest[0].Image)
	}

	if !latest[0].Date.Time.Equal(dates[2]) {
		t.Errorf("expected date %v, got %v", dates[2], latest[0].Date.Time)
	}

	if out := EPICImages(nil).mostRecent(); len(out) != 0 {
		t.Errorf("expected empty slice for nil input, got %d", len(out))
	}

	if out := (EPICImages{nil}).mostRecent(); len(out) != 0 {
		t.Errorf("expected empty slice when all entries are nil, got %d", len(out))
	}
}

func TestImageIDForCollection(t *testing.T) {
	t.Run("enhanced from natural image", func(t *testing.T) {
		got := imageIDForCollection("epic_1b_20200424002712", "20200424002712", "enhanced")
		expected := "epic_RGB_20200424002712"
		if got != expected {
			t.Fatalf("expected %s, got %s", expected, got)
		}
	})

	t.Run("natural from enhanced image", func(t *testing.T) {
		got := imageIDForCollection("epic_RGB_20200424002712", "20200424002712", "natural")
		expected := "epic_1b_20200424002712"
		if got != expected {
			t.Fatalf("expected %s, got %s", expected, got)
		}
	})

	t.Run("fallback to identifier", func(t *testing.T) {
		got := imageIDForCollection("", "20200424002712", "enhanced")
		expected := "epic_RGB_20200424002712"
		if got != expected {
			t.Fatalf("expected %s, got %s", expected, got)
		}
	})
}
