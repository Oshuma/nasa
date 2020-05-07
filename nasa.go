package nasa

import (
	"io/ioutil"
	"net/http"
)

// Version is the package version.
const Version = "0.1.1"

// LatLon represents latitude/longitude coordinates.
type LatLon struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Quaternions represents quadrant coordinates.
type Quaternions struct {
	Q0 float64 `json:"q0"`
	Q1 float64 `json:"q1"`
	Q2 float64 `json:"q2"`
	Q3 float64 `json:"q3"`
}

// XYZ represents coordinates in 3D space.
type XYZ struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func getContent(url string, p ParamEncoder) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if p != nil {
		query, err := p.Encode()
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = query
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return content, nil
}
