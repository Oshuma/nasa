package nasa

import (
	"errors"
)

var (
	// ErrorNoAPIKey is returned with no API key is given.
	ErrorNoAPIKey = errors.New("no API key provided; get one at https://api.nasa.gov")
)
