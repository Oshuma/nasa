package nasa

import (
	"io"
	"net/http"
	"time"
)

// Version is the package version.
const Version = "0.4.0"

// LatLon represents latitude/longitude coordinates as returned by several NASA
// endpoints, including EPIC image metadata.
type LatLon struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Quaternions represents spacecraft attitude quaternions. EPIC responses embed this
// orientation information for the DSCOVR spacecraft.
type Quaternions struct {
	Q0 float64 `json:"q0"`
	Q1 float64 `json:"q1"`
	Q2 float64 `json:"q2"`
	Q3 float64 `json:"q3"`
}

// XYZ represents coordinates in 3D space. NASA services commonly use this shape for
// positional data such as DSCOVR-to-sun vectors.
type XYZ struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// httpClient is used for every API request. The timeout keeps a stalled server from
// blocking callers forever.
var httpClient = &http.Client{Timeout: 2 * time.Minute}

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

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, &ErrorHTTPStatus{StatusCode: resp.StatusCode, Body: string(content)}
	}

	return content, nil
}
