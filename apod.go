package nasa

import (
	"bytes"
	"encoding/json"
)

const apodAPIURL = "https://api.nasa.gov/planetary/apod"

// APODImage represents an Astronomy Picture Of the Day and mirrors the fields returned
// by https://api.nasa.gov/planetary/apod. Most callers obtain an APODImage via APOD or
// APODList rather than instantiating the struct directly.
type APODImage struct {
	Date           Date   `json:"date"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	HDURL          string `json:"hdurl"`
	ThumbnailURL   string `json:"thumbnail_url"`
	Explanation    string `json:"explanation"`
	MediaType      string `json:"media_type"`
	Copyright      string `json:"copyright"`
	ServiceVersion string `json:"service_version"`
}

// APOD returns the Astronomy Picture Of the Day. Provide an APODParams value with at
// least an API key; optional filters such as Date control which image is fetched.
func APOD(p ParamEncoder) (APODImage, error) {
	content, err := getContent(apodAPIURL, p)
	if err != nil {
		return APODImage{}, err
	}

	images, err := parseAPODContent(content)
	if err != nil {
		return APODImage{}, err
	}

	if len(images) == 0 {
		return APODImage{}, nil
	}

	return images[0], nil
}

// APODList returns a slice of APOD images for queries that request multiple items, such as
// a date range or random count. See APODParams for the supported filters.
func APODList(p ParamEncoder) ([]APODImage, error) {
	content, err := getContent(apodAPIURL, p)
	if err != nil {
		return nil, err
	}

	return parseAPODContent(content)
}

func parseAPODContent(content []byte) ([]APODImage, error) {
	trimmed := bytes.TrimSpace(content)
	if len(trimmed) == 0 {
		return nil, nil
	}

	if trimmed[0] == '{' {
		img := APODImage{}
		if err := json.Unmarshal(content, &img); err != nil {
			return nil, err
		}
		return []APODImage{img}, nil
	}

	images := []APODImage{}
	if err := json.Unmarshal(content, &images); err != nil {
		return nil, err
	}

	return images, nil
}
