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

	// ErrorNoQuery is returned if there is no search query provided.
	ErrorNoQuery = errors.New("must provide a search query")

	// ErrorParamsMismatch is returned when the wrong type of ParamEncoder is used.
	ErrorParamsMismatch = errors.New("wrong param type passed")
)

// ErrorRoverCameraMissing is returned if the rover does not have the camera available.
type ErrorRoverCameraMissing struct {
	rover  Rover
	camera RoverCamera
}

func (e *ErrorRoverCameraMissing) Error() string {
	return fmt.Sprintf("rover %s does not have %s camera", e.rover.Name, e.camera.Name)
}
