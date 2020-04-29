package nasa

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const apodAPIURL = "https://api.nasa.gov/planetary/apod"

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

func APOD(p ParamEncoder) (APODImage, error) {
	req, err := http.NewRequest("GET", apodAPIURL, nil)
	if err != nil {
		return APODImage{}, err
	}

	query, err := p.Encode()
	if err != nil {
		return APODImage{}, err
	}
	req.URL.RawQuery = query

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return APODImage{}, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APODImage{}, err
	}
	defer resp.Body.Close()

	img := APODImage{}
	err = json.Unmarshal(content, &img)
	if err != nil {
		return APODImage{}, err
	}

	return img, nil
}
