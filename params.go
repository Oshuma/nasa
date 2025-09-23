package nasa

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ParamEncoder is the interface passed to most API methods.
type ParamEncoder interface {
	Encode() (string, error)
}

// APIParam is used when only an API key is needed. Most NASA endpoints require at least
// this parameter.
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

// APODParams wraps the APOD API params used by APOD and APODList. Populate APIKey and any
// optional filters such as Date or StartDate/EndDate.
type APODParams struct {
	APIKey    string
	Date      time.Time
	HD        bool
	StartDate time.Time
	EndDate   time.Time
	Count     int
	Thumbs    bool
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

	if !p.StartDate.IsZero() {
		v.Set("start_date", p.StartDate.Format("2006-01-02"))
	}

	if !p.EndDate.IsZero() {
		v.Set("end_date", p.EndDate.Format("2006-01-02"))
	}

	if p.Count > 0 {
		v.Set("count", strconv.Itoa(p.Count))
	}

	if p.HD {
		v.Set("hd", "true")
	}

	if p.Thumbs {
		v.Set("thumbs", "true")
	}

	return v.Encode(), nil
}

// EPICParams wraps the EPIC API params. Provide the desired collection (natural or
// enhanced), optional Date, and API key if required.
type EPICParams struct {
	APIKey     string
	Date       time.Time
	Collection string
}

// Encode returns a string representation for the given API type.
func (p *EPICParams) Encode() (string, error) {
	collection := strings.Trim(p.Collection, "/")
	if collection == "" {
		collection = "natural"
	}

	path := fmt.Sprintf("api/%s", collection)

	if !p.Date.IsZero() {
		path += fmt.Sprintf("/date/%s", p.Date.Format("2006-01-02"))
	}

	if p.APIKey == "" {
		return path, nil
	}

	return fmt.Sprintf("%s?api_key=%s", path, p.APIKey), nil
}

// MarsPhotosParams wraps the Mars Photos API params. Use it with MarsRoverPhotos or
// MarsRoverPhotosLatest to supply an API key, sol/earth date, camera, and pagination.
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

// MediaParams wraps the Image and Video Library (media) search parameters. At least one
// non-pagination filter, such as Query, NasaID, or Keywords, must be specified before a
// search can be performed.
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
	PageSize         int
	Photographer     string
	SecondaryCreator string
	Title            string
	YearStart        string
	YearEnd          string
}

// Encode returns a string representation for the given API type.
func (p *MediaParams) Encode() (string, error) {
	v := url.Values{}

	hasSearchParam := false

	if p.Query != "" {
		hasSearchParam = true
		v.Set("q", p.Query)
	}

	if p.Center != "" {
		hasSearchParam = true
		v.Set("center", p.Center)
	}

	if p.Description != "" {
		hasSearchParam = true
		v.Set("description", p.Description)
	}

	if p.Description508 != "" {
		hasSearchParam = true
		v.Set("description_508", p.Description508)
	}

	if p.Keywords != "" {
		hasSearchParam = true
		v.Set("keywords", p.Keywords)
	}

	if p.Location != "" {
		hasSearchParam = true
		v.Set("location", p.Location)
	}

	if p.MediaType != "" {
		hasSearchParam = true
		v.Set("media_type", p.MediaType)
	}

	if p.NasaID != "" {
		hasSearchParam = true
		v.Set("nasa_id", p.NasaID)
	}

	if p.Page > 0 {
		v.Set("page", strconv.Itoa(p.Page))
	}

	if p.PageSize > 0 {
		v.Set("page_size", strconv.Itoa(p.PageSize))
	}

	if p.Photographer != "" {
		hasSearchParam = true
		v.Set("photographer", p.Photographer)
	}

	if p.SecondaryCreator != "" {
		hasSearchParam = true
		v.Set("secondary_creator", p.SecondaryCreator)
	}

	if p.Title != "" {
		hasSearchParam = true
		v.Set("title", p.Title)
	}

	if p.YearStart != "" {
		hasSearchParam = true
		v.Set("year_start", p.YearStart)
	}

	if p.YearEnd != "" {
		hasSearchParam = true
		v.Set("year_end", p.YearEnd)
	}

	if !hasSearchParam {
		return "", ErrorNoQuery
	}

	return v.Encode(), nil
}
