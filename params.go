package nasa

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// ParamEncoder is the interface passed to most API methods.
type ParamEncoder interface {
	Encode() (string, error)
}

// APIParam is used when only an APIKey is needed.
type APIParam struct {
	APIKey string
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

// MediaParams wraps the Image and Video Library (media) params.
type MediaParams struct {
	Query            string
	Center           string
	Description      string
	Description508   string
	Keywords         string
	Location         string
	MediaType        string
	NasaID           string
	Page             int
	Photographer     string
	SecondaryCreator string
	Title            string
	YearStart        string
	YearEnd          string
}

// Encode returns a string representation for the given API type.
func (p *MediaParams) Encode() (string, error) {
	v := url.Values{}

	if p.Query == "" {
		return "", ErrorNoQuery
	}
	v.Set("q", p.Query)

	if p.Center != "" {
		v.Set("center", p.Center)
	}

	if p.Description != "" {
		v.Set("description", p.Description)
	}

	if p.Description508 != "" {
		v.Set("description_508", p.Description508)
	}

	if p.Keywords != "" {
		v.Set("keywords", p.Keywords)
	}

	if p.Location != "" {
		v.Set("location", p.Location)
	}

	if p.MediaType != "" {
		v.Set("media_type", p.MediaType)
	}

	if p.NasaID != "" {
		v.Set("nasa_id", p.NasaID)
	}

	if p.Page > 0 {
		v.Set("page", strconv.Itoa(p.Page))
	}

	if p.Photographer != "" {
		v.Set("photographer", p.Photographer)
	}

	if p.SecondaryCreator != "" {
		v.Set("secondary_creator", p.SecondaryCreator)
	}

	if p.Title != "" {
		v.Set("title", p.Title)
	}

	if p.YearStart != "" {
		v.Set("year_start", p.YearStart)
	}

	if p.YearEnd != "" {
		v.Set("year_end", p.YearEnd)
	}

	return v.Encode(), nil
}
