package nasa

import (
	"encoding/json"
)

const apodAPIURL = "https://api.nasa.gov/planetary/apod"

// APODImage represents an Astronomy Picture Of the Day.
type APODImage struct {
	Date           Date   `json:"date"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	HDURL          string `json:"hdurl"`
	Explanation    string `json:"explanation"`
	MediaType      string `json:"media_type"`
	Copyright      string `json:"copyright"`
	ServiceVersion string `json:"service_version"`
}

// APOD returns the Astronomy Picture Of the Day.
func APOD(p ParamEncoder) (APODImage, error) {
	content, err := getContent(apodAPIURL, p)
	if err != nil {
		return APODImage{}, err
	}

	img := APODImage{}
	err = json.Unmarshal(content, &img)
	if err != nil {
		return APODImage{}, err
	}

	return img, nil
}
