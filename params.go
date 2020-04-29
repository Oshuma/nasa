package nasa

import (
	"net/url"
	"time"
)

type ParamEncoder interface {
	Encode() (string, error)
}

type APODParams struct {
	APIKey string
	Date   time.Time
	HD     bool
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
