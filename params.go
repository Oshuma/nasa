package nasa

import (
	"fmt"
	"net/url"
	"time"
)

type ParamEncoder interface {
	GetAPIKey() string
	Encode() (string, error)
}

type APODParams struct {
	APIKey string
	Date   time.Time
	HD     bool
}

func (p *APODParams) GetAPIKey() string {
	return p.APIKey
}

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

type EPICParams struct {
	APIKey string
	Date   time.Time
}

func (p *EPICParams) GetAPIKey() string {
	return p.APIKey
}

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
