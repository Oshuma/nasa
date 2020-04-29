package nasa

import (
	"errors"
)

var (
	ErrorNoAPIKey = errors.New("no API key provided; get one at https://api.nasa.gov")
)
