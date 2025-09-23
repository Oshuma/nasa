package nasa

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	epicBaseURL        = "https://epic.gsfc.nasa.gov"
	epicImageURLFormat = "%s/archive/%s/%s/%s/%s/%s/%s.%s"
)

// EPICImage represents an image from the Earth Polychromatic Imaging Camera. The
// helper populates archive URLs so clients can link directly to the natural or
// enhanced imagery returned by the EPIC API.
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

// EPIC fetches imagery metadata from the Earth Polychromatic Imaging Camera. Populate
// EPICParams with the desired collection and optional date; when Date is omitted the
// function returns the most recent image for the requested collection. The helper also
// annotates each image with convenient archive download URLs.
func EPIC(p ParamEncoder) (EPICImages, error) {
	params, ok := p.(*EPICParams)
	if !ok {
		return EPICImages{}, ErrorParamsMismatch
	}

	query, err := p.Encode()
	if err != nil {
		return EPICImages{}, err
	}

	url := fmt.Sprintf("%s/%s", epicBaseURL, query)
	content, err := getContent(url, nil)
	if err != nil {
		return EPICImages{}, err
	}

	images := EPICImages{}
	err = json.Unmarshal(content, &images)
	if err != nil {
		return EPICImages{}, err
	}

	images.buildURLs(params)

	if params.Date.IsZero() {
		images = images.mostRecent()
	}

	return images, nil
}

// EPICImages is an array of pointers to EPICImage returned by the EPIC helper.
type EPICImages []*EPICImage

func (images EPICImages) buildURLs(p *EPICParams) {
	for _, epic := range images {
		epic.buildNaturalURLs(p)
		epic.buildEnhancedURLs(p)
	}
}

func (images EPICImages) mostRecent() EPICImages {
	var latest *EPICImage

	for _, candidate := range images {
		if candidate == nil {
			continue
		}
		if latest == nil || candidate.Date.Time.After(latest.Date.Time) {
			latest = candidate
		}
	}

	if latest == nil {
		return EPICImages{}
	}

	return EPICImages{latest}
}

// Full:  https://epic.gsfc.nasa.gov/archive/natural/2020/04/24/png/epic_1b_20200424002712.png
// Thumb: https://epic.gsfc.nasa.gov/archive/natural/2020/04/24/thumbs/epic_1b_20200424002712.jpg
func (e *EPICImage) buildNaturalURLs(p *EPICParams) {
	imageID := imageIDForCollection(e.Image, e.Identifier, "natural")
	e.URL.Natural = buildEPICArchiveURL(p, e.Date, "natural", "png", imageID, "png")
	e.URL.Thumb.Natural = buildEPICArchiveURL(p, e.Date, "natural", "thumbs", imageID, "jpg")
}

// Full:  https://epic.gsfc.nasa.gov/archive/enhanced/2020/04/24/png/epic_RGB_20200424002712.png
// Thumb: https://epic.gsfc.nasa.gov/archive/enhanced/2020/04/24/thumbs/epic_RGB_20200424002712.jpg
func (e *EPICImage) buildEnhancedURLs(p *EPICParams) {
	enhancedID := imageIDForCollection(e.Image, e.Identifier, "enhanced")

	e.URL.Enhanced = buildEPICArchiveURL(p, e.Date, "enhanced", "png", enhancedID, "png")
	e.URL.Thumb.Enhanced = buildEPICArchiveURL(p, e.Date, "enhanced", "thumbs", enhancedID, "jpg")
}

func buildEPICArchiveURL(p *EPICParams, date EPICDate, collection, imageType, imageID, extension string) string {
	base := fmt.Sprintf(
		epicImageURLFormat,
		epicBaseURL,
		collection,
		date.Format("2006"),
		date.Format("01"),
		date.Format("02"),
		imageType,
		imageID,
		extension,
	)

	if p != nil && p.APIKey != "" {
		return fmt.Sprintf("%s?api_key=%s", base, p.APIKey)
	}

	return base
}

func imageIDForCollection(image, identifier, collection string) string {
	switch collection {
	case "natural":
		if strings.Contains(image, "_1b_") {
			return image
		}
		if strings.Contains(image, "_RGB_") {
			return strings.Replace(image, "_RGB_", "_1b_", 1)
		}
		if strings.Contains(image, "_rgb_") {
			return strings.Replace(image, "_rgb_", "_1b_", 1)
		}
		if identifier != "" {
			return fmt.Sprintf("epic_1b_%s", identifier)
		}
	case "enhanced":
		if strings.Contains(image, "_RGB_") {
			return image
		}
		if strings.Contains(image, "_rgb_") {
			return strings.Replace(image, "_rgb_", "_RGB_", 1)
		}
		if strings.Contains(image, "_1b_") {
			return strings.Replace(image, "_1b_", "_RGB_", 1)
		}
		if identifier != "" {
			return fmt.Sprintf("epic_RGB_%s", identifier)
		}
	}

	return image
}
