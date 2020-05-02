package nasa

const (
	marsRoversAPIURL = "https://api.nasa.gov/mars-photos/api/v1/rovers"
)

// Rover represents a Mars rover.
type Rover struct {
	Name    string
	Slug    string
	Cameras []RoverCamera
}

// Defines Rovers to be used in the API request.
var (
	RoverCuriosity = Rover{
		Name:    "Curiosity",
		Slug:    "curiosity",
		Cameras: []RoverCamera{RoverCameraFHAZ, RoverCameraRHAZ, RoverCameraMAST, RoverCameraCHEMCAM, RoverCameraMAHLI, RoverCameraMARDI, RoverCameraNAVCAM},
	}
	RoverOpportunity = Rover{
		Name:    "Opportunity",
		Slug:    "opportunity",
		Cameras: []RoverCamera{RoverCameraFHAZ, RoverCameraRHAZ, RoverCameraNAVCAM, RoverCameraPANCAM, RoverCameraMINITES},
	}
	RoverSpirit = Rover{
		Name:    "Spirit",
		Slug:    "spirit",
		Cameras: []RoverCamera{RoverCameraFHAZ, RoverCameraRHAZ, RoverCameraNAVCAM, RoverCameraPANCAM, RoverCameraMINITES},
	}
)

// RoverCamera represents a rover camera type.
type RoverCamera struct {
	Name     string
	FullName string
	Slug     string
}

var (
	// RoverCameraFHAZ is the Front Hazard Avoidance Camera.
	RoverCameraFHAZ = RoverCamera{
		Name:     "FHAZ",
		FullName: "Front Hazard Avoidance Camera",
		Slug:     "fhaz",
	}

	// RoverCameraRHAZ is the Rear Hazard Avoidance Camera.
	RoverCameraRHAZ = RoverCamera{
		Name:     "RHAZ",
		FullName: "Rear Hazard Avoidance Camera",
		Slug:     "rhaz",
	}

	// RoverCameraMAST is the Mast Camera.
	RoverCameraMAST = RoverCamera{
		Name:     "MAST",
		FullName: "Mast Camera",
		Slug:     "mast",
	}

	// RoverCameraCHEMCAM is the Chemistry and Camera Complex.
	RoverCameraCHEMCAM = RoverCamera{
		Name:     "CHEMCAM",
		FullName: "Chemistry and Camera Complex",
		Slug:     "chemcam",
	}

	// RoverCameraMAHLI is the Mars Hand Lens Imager.
	RoverCameraMAHLI = RoverCamera{
		Name:     "MAHLI",
		FullName: "Mars Hand Lens Imager",
		Slug:     "mahli",
	}

	// RoverCameraMARDI is the Mars Descent Imager.
	RoverCameraMARDI = RoverCamera{
		Name:     "MARDI",
		FullName: "Mars Descent Imager",
		Slug:     "mardi",
	}

	// RoverCameraNAVCAM is the Navigation Camera.
	RoverCameraNAVCAM = RoverCamera{
		Name:     "NAVCAM",
		FullName: "Navigation Camera",
		Slug:     "navcam",
	}

	// RoverCameraPANCAM is the Panoramic Camera.
	RoverCameraPANCAM = RoverCamera{
		Name:     "PANCAM",
		FullName: "Panoramic Camera",
		Slug:     "pancam",
	}

	// RoverCameraMINITES is the Miniature Thermal Emission Spectrometer (Mini-TES).
	RoverCameraMINITES = RoverCamera{
		Name:     "MINITES",
		FullName: "Miniature Thermal Emission Spectrometer (Mini-TES)",
		Slug:     "minites",
	}
)
