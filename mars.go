package nasa

import (
	"encoding/json"
	"fmt"
)

const (
	marsPhotosAPIURL          = "https://api.nasa.gov/mars-photos/api/v1/rovers/%s/photos"
	marsLatestPhotosAPIURL    = "https://api.nasa.gov/mars-photos/api/v1/rovers/%s/latest_photos"
	marsPhotosManifestsAPIURL = "https://api.nasa.gov/mars-photos/api/v1/manifests/%s"
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
	params, ok := p.(*MarsPhotosParams)
	if !ok {
		return RoverPhotos{}, ErrorParamsMismatch
	}

	camera := params.Camera
	if camera.Slug != "" {
		if !hasCamera(rover, camera) {
			return RoverPhotos{}, &ErrorRoverCameraMissing{rover, camera}
		}
	}

	url := fmt.Sprintf(marsPhotosAPIURL, rover.Slug)
	content, err := getContent(url, p)
	if err != nil {
		return RoverPhotos{}, err
	}

	photos := RoverPhotos{}
	err = json.Unmarshal(content, &photos)
	if err != nil {
		return RoverPhotos{}, err
	}

	page := 1
	if params.Page > 1 {
		page = params.Page
	}
	photos.Page = page

	return photos, nil
}

type latestPhotosResponse struct {
	Photos []*RoverPhoto `json:"latest_photos"`
}

// MarsRoverPhotosLatest returns the most recent Sol for which photos exist.
func MarsRoverPhotosLatest(p ParamEncoder, rover Rover) ([]*RoverPhoto, error) {
	url := fmt.Sprintf(marsLatestPhotosAPIURL, rover.Slug)
	content, err := getContent(url, p)
	if err != nil {
		return []*RoverPhoto{}, err
	}

	r := latestPhotosResponse{}
	err = json.Unmarshal(content, &r)
	if err != nil {
		return []*RoverPhoto{}, err
	}

	return r.Photos, nil
}

// MissionManifest represents rover mission details.
type MissionManifest struct {
	Name        string `json:"name"`
	LandingDate Date   `json:"landing_date"`
	LaunchDate  Date   `json:"launch_date"`
	Status      string `json:"status"`
	MaxSol      int    `json:"max_sol"`
	MaxDate     Date   `json:"max_date"`
	TotalPhotos int    `json:"total_photos"`
	Photos      []struct {
		Sol         int      `json:"sol"`
		EarthDate   Date     `json:"earth_date"`
		TotalPhotos int      `json:"total_photos"`
		Cameras     []string `json:"cameras"`
	} `json:"photos"`
}

type manifestResponse struct {
	Manifest MissionManifest `json:"photo_manifest"`
}

// MarsMissionManifest returns the rover mission details.
func MarsMissionManifest(p ParamEncoder, rover Rover) (MissionManifest, error) {
	url := fmt.Sprintf(marsPhotosManifestsAPIURL, rover.Slug)
	content, err := getContent(url, p)
	if err != nil {
		return MissionManifest{}, err
	}

	r := manifestResponse{}
	err = json.Unmarshal(content, &r)
	if err != nil {
		return MissionManifest{}, err
	}

	return r.Manifest, nil
}

func hasCamera(rover Rover, camera RoverCamera) bool {
	for _, c := range rover.Cameras {
		if c == camera {
			return true
		}
	}
	return false
}
