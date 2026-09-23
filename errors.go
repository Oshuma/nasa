package nasa

import (
	"errors"
	"fmt"
)

var (
	// ErrorNoAPIKey is returned with no API key is given.
	ErrorNoAPIKey = errors.New("no API key provided; get one at https://api.nasa.gov")

	// ErrorNoMetadata is returned if a media asset has no metadata.
	ErrorNoMetadata = errors.New("media has no metadata")

	// ErrorNoCaptions is returned if a media asset has no captions.
	ErrorNoCaptions = errors.New("media has no captions")

	// ErrorNoQuery is returned if there are no search parameters provided.
	ErrorNoQuery = errors.New("must provide at least one search parameter")

	// ErrorParamsMismatch is returned when the wrong type of ParamEncoder is used.
	ErrorParamsMismatch = errors.New("wrong param type passed")

	// ErrorRoverUnsupported is returned when photos are requested for a rover that has no
	// supported image source.
	ErrorRoverUnsupported = errors.New("rover is not supported")
)

// ErrorRoverCameraMissing is returned if the rover does not have the camera available.
type ErrorRoverCameraMissing struct {
	rover  Rover
	camera RoverCamera
}

func (e *ErrorRoverCameraMissing) Error() string {
	return fmt.Sprintf("rover %s does not have %s camera", e.rover.Name, e.camera.Name)
}

// ErrorHTTPStatus is returned when an API responds with a non-2xx status code. Body holds
// the raw response body, which often contains the API's own error message.
type ErrorHTTPStatus struct {
	StatusCode int
	Body       string
}

func (e *ErrorHTTPStatus) Error() string {
	return fmt.Sprintf("unexpected HTTP status %d: %s", e.StatusCode, e.Body)
}
