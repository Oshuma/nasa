package nasa

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// ParamEncoder is the interface passed to most API methods.
type ParamEncoder interface {
	// TODO: Can probably do away with this and just cast when needed.
	GetAPIKey() string
	Encode() (string, error)
}

// APIParam is used when only an APIKey is needed.
type APIParam struct {
	APIKey string
}

func (p *APIParam) GetAPIKey() string {
	return p.APIKey
}

// Encode returns a string representation for the given API type.
func (p *APIParam) Encode() (string, error) {
	v := url.Values{}

	if p.APIKey == "" {
		return "", ErrorNoAPIKey
	}
	v.Set("api_key", p.APIKey)

	return v.Encode(), nil
}

// APODParams wraps the APOD API params.
type APODParams struct {
	APIKey string
	Date   time.Time
	HD     bool
}

// GetAPIKey returns the APIKey.
func (p *APODParams) GetAPIKey() string {
	return p.APIKey
}

// Encode returns a string representation for the given API type.
func (p *APODParams) Encode() (string, error) {
	v := url.Values{}

	if p.APIKey == "" {
		return "", ErrorNoAPIKey
	}
	v.Set("api_key", p.APIKey)

	if !p.Date.IsZero() {
		v.Set("date", p.Date.Format("2006-01-02"))
	}

	if p.HD {
		v.Set("hd", "true")
	}

	return v.Encode(), nil
}

// EPICParams wraps the EPIC API params.
type EPICParams struct {
	APIKey string
	Date   time.Time
}

// GetAPIKey returns the APIKey.
func (p *EPICParams) GetAPIKey() string {
	return p.APIKey
}

// Encode returns a string representation for the given API type.
func (p *EPICParams) Encode() (string, error) {
	if p.APIKey == "" {
		return "", ErrorNoAPIKey
	}

	val := "api/natural"

	if !p.Date.IsZero() {
		val += fmt.Sprintf("/date/%s", p.Date.Format("2006-01-02"))
	}

	val += fmt.Sprintf("?api_key=%s", p.APIKey)

	return val, nil
}

// MarsPhotosParams wraps the Mars Photos API params.
type MarsPhotosParams struct {
	APIKey    string
	Sol       int
	EarthDate time.Time
	Camera    RoverCamera
	Page      int
}

// GetAPIKey returns the APIKey.
func (p *MarsPhotosParams) GetAPIKey() string {
	return p.APIKey
}

// Encode returns a string representation for the given API type.
func (p *MarsPhotosParams) Encode() (string, error) {
	v := url.Values{}

	if p.APIKey == "" {
		return "", ErrorNoAPIKey
	}
	v.Set("api_key", p.APIKey)

	if !p.EarthDate.IsZero() {
		v.Set("earth_date", p.EarthDate.Format("2006-01-02"))
	} else {
		v.Set("sol", strconv.Itoa(p.Sol))
	}

	if p.Camera.Slug != "" {
		v.Set("camera", p.Camera.Slug)
	}

	if p.Page > 0 {
		v.Set("page", strconv.Itoa(p.Page))
	}

	return v.Encode(), nil
}
