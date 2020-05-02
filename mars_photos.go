package nasa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	marsPhotosAPIURL          = "https://api.nasa.gov/mars-photos/api/v1/rovers/%s/photos"
	marsPhotosManifestsAPIURL = "https://api.nasa.gov/mars-photos/api/v1/manifests"
)

// Rover represents a Mars rover.
type Rover string

// Defines Rovers to be used in the API request.
var (
	RoverCuriosity   Rover = "curiosity"
	RoverOpportunity Rover = "opportunity"
	RoverSpirit      Rover = "spirit"
)

// RoverCamera represents a rover camera type.
type RoverCamera string

var (
	// RoverCameraFHAZ is the Front Hazard Avoidance Camera
	RoverCameraFHAZ RoverCamera = "fhaz"
	// RoverCameraRHAZ is the Rear Hazard Avoidance Camera
	RoverCameraRHAZ RoverCamera = "rhaz"
	// RoverCameraMAST is the Mast Camera
	RoverCameraMAST RoverCamera = "mast"
	// RoverCameraCHEMCAM is the Chemistry and Camera Complex
	RoverCameraCHEMCAM RoverCamera = "chemcam"
	// RoverCameraMAHLI is the Mars Hand Lens Imager
	RoverCameraMAHLI RoverCamera = "mahli"
	// RoverCameraMARDI is the Mars Descent Imager
	RoverCameraMARDI RoverCamera = "mardi"
	// RoverCameraNAVCAM is the Navigation Camera
	RoverCameraNAVCAM RoverCamera = "navcam"
	// RoverCameraPANCAM is the Panoramic Camera
	RoverCameraPANCAM RoverCamera = "pancam"
	// RoverCameraMINITES is the Miniature Thermal Emission Spectrometer (Mini-TES)
	RoverCameraMINITES RoverCamera = "minites"
)

// RoverPhoto represents a single photo from a rover camera.
type RoverPhoto struct {
	ID        int    `json:"id"`
	Sol       int    `json:"sol"`
	Image     string `json:"img_src"`
	EarthDate Date   `json:"earth_date"`
	Camera    struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		RoverID  int    `json:"rover_id"`
	} `json:"camera"`
	Rover struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		LandingDate Date   `json:"landing_date"`
		LaunchDate  Date   `json:"launch_date"`
		Status      string `json:"status"`
		MaxSol      int    `json:"max_sol"`
		MaxDate     Date   `json:"max_date"`
		TotalPhotos int    `json:"total_photos"`
		Cameras     []struct {
			Name     string `json:"name"`
			FullName string `json:"full_name"`
		} `json:"cameras"`
	} `json:"rover"`
}

// RoverPhotos wraps an array of pointers of RoverPhoto.
type RoverPhotos struct {
	Photos []*RoverPhoto `json:"photos"`
	Page   int
}

// MarsRoverPhotos returns photos for the given params and Rover.
func MarsRoverPhotos(p ParamEncoder, rover Rover) (RoverPhotos, error) {
	url := fmt.Sprintf(marsPhotosAPIURL, rover)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RoverPhotos{}, err
	}

	query, err := p.Encode()
	if err != nil {
		return RoverPhotos{}, err
	}
	req.URL.RawQuery = query

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return RoverPhotos{}, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RoverPhotos{}, err
	}
	defer resp.Body.Close()

	photos := RoverPhotos{}
	err = json.Unmarshal(content, &photos)
	if err != nil {
		return RoverPhotos{}, err
	}

	page := p.(*MarsPhotosParams).Page
	if page <= 0 {
		page = 1
	}
	photos.Page = page

	return photos, nil
}
