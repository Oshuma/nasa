package nasa

import (
	"slices"
	"sort"
	"time"
)

const (
	// marsPhotosPerPage matches the page size of the retired Mars Rover Photos API.
	marsPhotosPerPage = 25

	// marsLatestSolLookback limits how many sols MarsRoverPhotosLatest walks back when the
	// newest sols have no full-frame photos yet.
	marsLatestSolLookback = 30
)

// RoverPhoto represents a single full-frame Mars rover image, fetched live from NASA's raw
// image feeds (Curiosity and Perseverance) or the PDS Imaging Atlas (Opportunity and Spirit).
type RoverPhoto struct {
	// ID is the feed's numeric image ID. Only Curiosity provides one; it is 0 for the
	// other rovers.
	ID int `json:"id"`

	// ImageID is the source's image identifier, available for every rover. For Opportunity
	// and Spirit it is the PDS product ID.
	ImageID string `json:"image_id"`

	Sol   int    `json:"sol"`
	Image string `json:"img_src"`

	// EarthDate is the UTC calendar date of TakenAt.
	EarthDate Date `json:"earth_date"`

	// TakenAt is when the image was captured, in UTC.
	TakenAt time.Time `json:"taken_at"`

	Camera struct {
		// ID and RoverID were database IDs in the retired API and are always 0.
		ID       int    `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		RoverID  int    `json:"rover_id"`

		// Instrument is the source's own instrument name, such as "NAV_LEFT_B". Opportunity
		// and Spirit report FRONT_HAZCAM, REAR_HAZCAM, NAVCAM, PANCAM or MI. Photos from
		// instruments that do not belong to one of the rover's RoverCameras use it as
		// their Name and FullName.
		Instrument string `json:"instrument"`
	} `json:"camera"`

	// Rover describes the rover that took the photo. ID, MaxSol, MaxDate and TotalPhotos
	// are not populated; use MarsMissionManifest for mission statistics.
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

// RoverPhotos wraps an array of RoverPhoto pointers along with the requested page number
// and the total number of matching photos across all pages. Use MarsRoverPhotos to
// populate this type.
type RoverPhotos struct {
	Photos []*RoverPhoto `json:"photos"`
	Page   int
	Total  int
}

// MarsRoverPhotos returns full-frame photos for the given MarsPhotosParams and Rover,
// fetched live (see RoverPhoto for the sources). Photos are selected by EarthDate (the UTC date
// each image was captured) when it is set, otherwise by Sol, and can be narrowed to one
// Camera. A Page of 1 or more returns that page of 25 photos; otherwise every matching
// photo is returned.
//
// Curiosity, Perseverance, Opportunity and Spirit are supported. A custom Rover returns
// ErrorRoverUnsupported.
func MarsRoverPhotos(p ParamEncoder, rover Rover) (RoverPhotos, error) {
	params, ok := p.(*MarsPhotosParams)
	if !ok {
		return RoverPhotos{}, ErrorParamsMismatch
	}

	mission, err := missionFor(rover)
	if err != nil {
		return RoverPhotos{}, err
	}

	instruments, err := mission.instrumentsFor(rover, params.Camera)
	if err != nil {
		return RoverPhotos{}, err
	}

	var photos []*RoverPhoto
	if !params.EarthDate.IsZero() {
		photos, err = mission.photosOnDate(rover, params.EarthDate, instruments)
	} else {
		photos, err = mission.photosOnSol(rover, params.Sol, instruments)
	}
	if err != nil {
		return RoverPhotos{}, err
	}

	page := 1
	if params.Page > 1 {
		page = params.Page
	}

	return RoverPhotos{
		Photos: paginate(photos, params.Page),
		Page:   page,
		Total:  len(photos),
	}, nil
}

// MarsRoverPhotosLatest returns the full-frame photos from the most recent sol that has
// any, fetched live (see RoverPhoto for the sources). The Camera and Page fields of
// MarsPhotosParams are applied to that sol; Sol and EarthDate are ignored.
//
// Curiosity, Perseverance, Opportunity and Spirit are supported. A custom Rover returns
// ErrorRoverUnsupported.
func MarsRoverPhotosLatest(p ParamEncoder, rover Rover) ([]*RoverPhoto, error) {
	params, ok := p.(*MarsPhotosParams)
	if !ok {
		return []*RoverPhoto{}, ErrorParamsMismatch
	}

	mission, err := missionFor(rover)
	if err != nil {
		return []*RoverPhoto{}, err
	}

	instruments, err := mission.instrumentsFor(rover, params.Camera)
	if err != nil {
		return []*RoverPhoto{}, err
	}

	_, photos, err := mission.latestPhotos(rover)
	if err != nil {
		return []*RoverPhoto{}, err
	}

	if len(instruments) > 0 {
		filtered := []*RoverPhoto{}
		for _, photo := range photos {
			if slices.Contains(instruments, photo.Camera.Instrument) {
				filtered = append(filtered, photo)
			}
		}
		photos = filtered
	}

	return paginate(photos, params.Page), nil
}

// MissionManifest represents rover mission details. MarsMissionManifest fills in a partial
// manifest: Photos is always empty.
type MissionManifest struct {
	Name        string `json:"name"`
	LandingDate Date   `json:"landing_date"`
	LaunchDate  Date   `json:"launch_date"`
	Status      string `json:"status"`

	// MaxSol and MaxDate are the sol and UTC capture date of the most recent full-frame
	// photos.
	MaxSol  int  `json:"max_sol"`
	MaxDate Date `json:"max_date"`

	// TotalPhotos is the image count reported by the rover's source. For Perseverance,
	// Opportunity and Spirit it counts full frames only; for Curiosity it also counts
	// thumbnails and subframes.
	TotalPhotos int `json:"total_photos"`

	// Photos is not populated, since building it requires one request per sol.
	Photos []ManifestPhoto `json:"photos"`
}

// ManifestPhoto stores info about mission manifest photos.
type ManifestPhoto struct {
	Sol         int      `json:"sol"`
	EarthDate   Date     `json:"earth_date"`
	TotalPhotos int      `json:"total_photos"`
	Cameras     []string `json:"cameras"`
}

// MarsMissionManifest returns a partial mission manifest for the rover, fetched live (see
// RoverPhoto for the sources). See MissionManifest for which fields are populated.
//
// Curiosity, Perseverance, Opportunity and Spirit are supported. A custom Rover returns
// ErrorRoverUnsupported.
func MarsMissionManifest(p ParamEncoder, rover Rover) (MissionManifest, error) {
	if _, ok := p.(*MarsPhotosParams); !ok {
		return MissionManifest{}, ErrorParamsMismatch
	}

	mission, err := missionFor(rover)
	if err != nil {
		return MissionManifest{}, err
	}

	maxSol, photos, err := mission.latestPhotos(rover)
	if err != nil {
		return MissionManifest{}, err
	}

	total, err := mission.source.totalPhotos()
	if err != nil {
		return MissionManifest{}, err
	}

	manifest := MissionManifest{
		Name:        rover.Name,
		LandingDate: Date{rover.LandingDate},
		LaunchDate:  Date{rover.LaunchDate},
		Status:      rover.Status,
		MaxSol:      maxSol,
		TotalPhotos: total,
	}
	for _, photo := range photos {
		if photo.EarthDate.After(manifest.MaxDate.Time) {
			manifest.MaxDate = photo.EarthDate
		}
	}

	return manifest, nil
}

func (m marsMission) photosOnSol(rover Rover, sol int, instruments []string) ([]*RoverPhoto, error) {
	sourcePhotos, err := m.source.fetchSol(sol, instruments)
	if err != nil {
		return nil, err
	}
	return m.roverPhotos(rover, sourcePhotos), nil
}

// photosOnDate returns the photos captured on the UTC calendar day of date. Sources that are
// not a dateSource cannot filter by capture time, so every sol overlapping that day is
// fetched and filtered here.
func (m marsMission) photosOnDate(rover Rover, date time.Time, instruments []string) ([]*RoverPhoto, error) {
	day := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	if ds, ok := m.source.(dateSource); ok {
		sourcePhotos, err := ds.fetchDate(day, instruments)
		if err != nil {
			return nil, err
		}
		return m.roverPhotos(rover, sourcePhotos), nil
	}

	sourcePhotos := []sourcePhoto{}
	for _, sol := range m.solsOn(day) {
		photos, err := m.source.fetchSol(sol, instruments)
		if err != nil {
			return nil, err
		}
		for _, photo := range photos {
			if !photo.takenAt.Before(day) && photo.takenAt.Before(day.AddDate(0, 0, 1)) {
				sourcePhotos = append(sourcePhotos, photo)
			}
		}
	}

	return m.roverPhotos(rover, sourcePhotos), nil
}

// latestPhotos returns the most recent sol with full-frame photos, and those photos. The
// newest sols often hold only thumbnails until full frames are downlinked, so this walks
// back up to marsLatestSolLookback sols.
func (m marsMission) latestPhotos(rover Rover) (int, []*RoverPhoto, error) {
	latest, err := m.source.latestSol()
	if err != nil {
		return 0, nil, err
	}

	for sol := latest; sol >= 0 && sol > latest-marsLatestSolLookback; sol-- {
		photos, err := m.photosOnSol(rover, sol, nil)
		if err != nil {
			return 0, nil, err
		}
		if len(photos) > 0 {
			return sol, photos, nil
		}
	}

	return latest, []*RoverPhoto{}, nil
}

// roverPhotos maps feed photos onto RoverPhotos, sorted by camera (in the order of
// rover.Cameras, unknown instruments last), then capture time, then image ID.
func (m marsMission) roverPhotos(rover Rover, sourcePhotos []sourcePhoto) []*RoverPhoto {
	cameraOrder := make(map[string]int, len(rover.Cameras))
	for i, c := range rover.Cameras {
		cameraOrder[c.Slug] = i
	}

	type sortable struct {
		photo *RoverPhoto
		order int
	}
	sorted := make([]sortable, 0, len(sourcePhotos))

	for _, sp := range sourcePhotos {
		photo := &RoverPhoto{
			ID:        sp.id,
			ImageID:   sp.imageID,
			Sol:       sp.sol,
			Image:     sp.url,
			EarthDate: Date{time.Date(sp.takenAt.Year(), sp.takenAt.Month(), sp.takenAt.Day(), 0, 0, 0, 0, time.UTC)},
			TakenAt:   sp.takenAt,
		}

		order := len(rover.Cameras)
		photo.Camera.Instrument = sp.instrument
		if camera, ok := m.cameraFor(rover, sp.instrument); ok {
			photo.Camera.Name = camera.Name
			photo.Camera.FullName = camera.FullName
			order = cameraOrder[camera.Slug]
		} else {
			photo.Camera.Name = sp.instrument
			photo.Camera.FullName = sp.instrument
		}

		photo.Rover.Name = rover.Name
		photo.Rover.LandingDate = Date{rover.LandingDate}
		photo.Rover.LaunchDate = Date{rover.LaunchDate}
		photo.Rover.Status = rover.Status
		for _, c := range rover.Cameras {
			photo.Rover.Cameras = append(photo.Rover.Cameras, struct {
				Name     string `json:"name"`
				FullName string `json:"full_name"`
			}{c.Name, c.FullName})
		}

		sorted = append(sorted, sortable{photo, order})
	}

	sort.SliceStable(sorted, func(i, j int) bool {
		a, b := sorted[i], sorted[j]
		if a.order != b.order {
			return a.order < b.order
		}
		if !a.photo.TakenAt.Equal(b.photo.TakenAt) {
			return a.photo.TakenAt.Before(b.photo.TakenAt)
		}
		return a.photo.ImageID < b.photo.ImageID
	})

	photos := make([]*RoverPhoto, len(sorted))
	for i, s := range sorted {
		photos[i] = s.photo
	}
	return photos
}

// paginate returns the requested page of marsPhotosPerPage photos, or every photo when page
// is less than 1.
func paginate(photos []*RoverPhoto, page int) []*RoverPhoto {
	if page < 1 {
		return photos
	}

	start := (page - 1) * marsPhotosPerPage
	if start >= len(photos) {
		return []*RoverPhoto{}
	}
	return photos[start:min(start+marsPhotosPerPage, len(photos))]
}

func hasCamera(rover Rover, camera RoverCamera) bool {
	for _, c := range rover.Cameras {
		if c == camera {
			return true
		}
	}
	return false
}
