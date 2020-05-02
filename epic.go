package nasa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	epicAPIURL         = "https://api.nasa.gov/EPIC"
	epicImageURLFormat = "%s/archive/%s/%s/%s/%s/%s/%s.%s?api_key=%s"
)

// EPICImage represents an image from the Earth Polychromatic Imaging Camera.
type EPICImage struct {
	Date       EPICDate `json:"date"`
	Identifier string   `json:"identifier"`
	Caption    string   `json:"caption"`
	Image      string   `json:"image"`
	Version    string   `json:"version"`
	Coords     struct {
		Centroid LatLon      `json:"centroid_coordinates"`
		Dscovr   XYZ         `json:"dscovr_j2000_position"`
		Lunar    XYZ         `json:"lunar_j2000_position"`
		Sun      XYZ         `json:"sun_j2000_position"`
		Attitude Quaternions `json:"attitude_quaternions"`
	} `json:"coords"`
	URL struct {
		Natural  string
		Enhanced string
		Thumb    struct {
			Natural  string
			Enhanced string
		} `json:"-"`
	} `json:"-"`
}

// EPIC gets a response from the Earth Polychromatic Imaging Camera.
func EPIC(p ParamEncoder) (EPICImages, error) {
	query, err := p.Encode()
	if err != nil {
		return EPICImages{}, err
	}

	url := fmt.Sprintf("%s/%s", epicAPIURL, query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return EPICImages{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return EPICImages{}, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return EPICImages{}, err
	}
	defer resp.Body.Close()

	images := EPICImages{}
	err = json.Unmarshal(content, &images)
	if err != nil {
		return EPICImages{}, err
	}

	images.buildURLs(p)

	return images, nil
}

// EPICImages is an array of pointers to EPICImage.
type EPICImages []*EPICImage

func (images EPICImages) buildURLs(p ParamEncoder) {
	for _, epic := range images {
		epic.buildNaturalURLs(p)
		epic.buildEnhancedURLs(p)
	}
}

// Full:  https://api.nasa.gov/EPIC/archive/natural/2020/04/24/png/epic_1b_20200424002712.png?api_key=DEMO_KEY
// Thumb: https://api.nasa.gov/EPIC/archive/natural/2020/04/24/thumbs/epic_1b_20200424002712.jpg?api_key=DEMO_KEY
func (e *EPICImage) buildNaturalURLs(p ParamEncoder) {
	e.URL.Natural = fmt.Sprintf(epicImageURLFormat,
		epicAPIURL,
		"natural",
		e.Date.Format("2006"), // Year
		e.Date.Format("01"),   // Month
		e.Date.Format("02"),   // Day
		"png",
		e.Image,
		"png",
		p.GetAPIKey(),
	)

	e.URL.Thumb.Natural = fmt.Sprintf(epicImageURLFormat,
		epicAPIURL,
		"natural",
		e.Date.Format("2006"), // Year
		e.Date.Format("01"),   // Month
		e.Date.Format("02"),   // Day
		"thumbs",
		e.Image,
		"jpg",
		p.GetAPIKey(),
	)
}

// Full:  https://api.nasa.gov/EPIC/archive/enhanced/2020/04/24/png/epic_RGB_20200424002712.png?api_key=DEMO_KEY
// Thumb: https://api.nasa.gov/EPIC/archive/enhanced/2020/04/24/thumbs/epic_RGB_20200424002712.jpg?api_key=DEMO_KEY
func (e *EPICImage) buildEnhancedURLs(p ParamEncoder) {
	enhancedID := strings.Replace(e.Image, "_1b_", "_RGB_", 1)

	e.URL.Enhanced = fmt.Sprintf(epicImageURLFormat,
		epicAPIURL,
		"enhanced",
		e.Date.Format("2006"), // Year
		e.Date.Format("01"),   // Month
		e.Date.Format("02"),   // Day
		"png",
		enhancedID,
		"png",
		p.GetAPIKey(),
	)

	e.URL.Thumb.Enhanced = fmt.Sprintf(epicImageURLFormat,
		epicAPIURL,
		"enhanced",
		e.Date.Format("2006"), // Year
		e.Date.Format("01"),   // Month
		e.Date.Format("02"),   // Day
		"thumbs",
		enhancedID,
		"jpg",
		p.GetAPIKey(),
	)
}
