package nasa

import (
	"testing"

	"time"
)

func TestBuildNaturalURLs(t *testing.T) {
	full := "https://api.nasa.gov/EPIC/archive/natural/2020/04/24/png/epic_1b_20200424002712.png?api_key=NASA_KEY"
	thumb := "https://api.nasa.gov/EPIC/archive/natural/2020/04/24/thumbs/epic_1b_20200424002712.jpg?api_key=NASA_KEY"

	p := &EPICParams{APIKey: "NASA_KEY"}
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
	full := "https://api.nasa.gov/EPIC/archive/enhanced/2020/04/24/png/epic_RGB_20200424002712.png?api_key=NASA_KEY"
	thumb := "https://api.nasa.gov/EPIC/archive/enhanced/2020/04/24/thumbs/epic_RGB_20200424002712.jpg?api_key=NASA_KEY"

	p := &EPICParams{APIKey: "NASA_KEY"}
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
