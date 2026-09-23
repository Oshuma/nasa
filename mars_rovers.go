package nasa

// Rover represents a Mars rover. Use the predefined Rover* values when calling
// MarsRoverPhotos or MarsMissionManifest to avoid constructing your own instance.
type Rover struct {
	Name    string
	Slug    string
	Cameras []RoverCamera
}

// Defines Rovers to be used in Mars rover API requests.
var (
	RoverCuriosity = Rover{
		Name:    "Curiosity",
		Slug:    "curiosity",
		Cameras: []RoverCamera{RoverCameraFHAZ, RoverCameraRHAZ, RoverCameraMAST, RoverCameraCHEMCAM, RoverCameraMAHLI, RoverCameraMARDI, RoverCameraNAVCAM},
	}
	RoverPerseverance = Rover{
		Name: "Perseverance",
		Slug: "perseverance",
		Cameras: []RoverCamera{
			RoverCameraEDLRUCAM,
			RoverCameraEDLRDCAM,
			RoverCameraEDLDDCAM,
			RoverCameraEDLPUCAM1,
			RoverCameraEDLPUCAM2,
			RoverCameraNAVCAMLEFT,
			RoverCameraNAVCAMRIGHT,
			RoverCameraMCZLEFT,
			RoverCameraMCZRIGHT,
			RoverCameraFRONTHAZCAMLEFTA,
			RoverCameraFRONTHAZCAMRIGHTA,
			RoverCameraREARHAZCAMLEFT,
			RoverCameraREARHAZCAMRIGHT,
			RoverCameraSKYCAM,
			RoverCameraSHERLOCWATSON,
			RoverCameraSUPERCAMRMI,
			RoverCameraLCAM,
			RoverCameraMEDA,
		},
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

	// Rovers is a convenience slice for iterating through every supported rover when
	// building UI pickers or issuing batch calls.
	Rovers = []Rover{
		RoverCuriosity,
		RoverPerseverance,
		RoverOpportunity,
		RoverSpirit,
	}
)

// RoverCamera represents a rover camera type. Refer to the exported RoverCamera*
// variables for valid camera selections per rover.
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

	// RoverCameraEDLRUCAM is the Entry, Descent, and Landing Rover Up-Look Camera.
	RoverCameraEDLRUCAM = RoverCamera{
		Name:     "EDL_RUCAM",
		FullName: "Entry, Descent, and Landing Rover Up-Look Camera",
		Slug:     "edl_rucam",
	}

	// RoverCameraEDLRDCAM is the Entry, Descent, and Landing Rover Down-Look Camera.
	RoverCameraEDLRDCAM = RoverCamera{
		Name:     "EDL_RDCAM",
		FullName: "Entry, Descent, and Landing Rover Down-Look Camera",
		Slug:     "edl_rdcam",
	}

	// RoverCameraEDLDDCAM is the Entry, Descent, and Landing Descent Stage Down-Look Camera.
	RoverCameraEDLDDCAM = RoverCamera{
		Name:     "EDL_DDCAM",
		FullName: "Entry, Descent, and Landing Descent Stage Down-Look Camera",
		Slug:     "edl_ddcam",
	}

	// RoverCameraEDLPUCAM1 is the Entry, Descent, and Landing Parachute Up-Look Camera A.
	RoverCameraEDLPUCAM1 = RoverCamera{
		Name:     "EDL_PUCAM1",
		FullName: "Entry, Descent, and Landing Parachute Up-Look Camera A",
		Slug:     "edl_pucam1",
	}

	// RoverCameraEDLPUCAM2 is the Entry, Descent, and Landing Parachute Up-Look Camera B.
	RoverCameraEDLPUCAM2 = RoverCamera{
		Name:     "EDL_PUCAM2",
		FullName: "Entry, Descent, and Landing Parachute Up-Look Camera B",
		Slug:     "edl_pucam2",
	}

	// RoverCameraNAVCAMLEFT is the Perseverance Navigation Camera (Left).
	RoverCameraNAVCAMLEFT = RoverCamera{
		Name:     "NAVCAM_LEFT",
		FullName: "Navigation Camera - Left",
		Slug:     "navcam_left",
	}

	// RoverCameraNAVCAMRIGHT is the Perseverance Navigation Camera (Right).
	RoverCameraNAVCAMRIGHT = RoverCamera{
		Name:     "NAVCAM_RIGHT",
		FullName: "Navigation Camera - Right",
		Slug:     "navcam_right",
	}

	// RoverCameraMCZLEFT is the Mastcam-Z Left Camera.
	RoverCameraMCZLEFT = RoverCamera{
		Name:     "MCZ_LEFT",
		FullName: "Mastcam-Z Left Camera",
		Slug:     "mcz_left",
	}

	// RoverCameraMCZRIGHT is the Mastcam-Z Right Camera.
	RoverCameraMCZRIGHT = RoverCamera{
		Name:     "MCZ_RIGHT",
		FullName: "Mastcam-Z Right Camera",
		Slug:     "mcz_right",
	}

	// RoverCameraFRONTHAZCAMLEFTA is the Front Hazard Avoidance Camera (Left A).
	RoverCameraFRONTHAZCAMLEFTA = RoverCamera{
		Name:     "FRONT_HAZCAM_LEFT_A",
		FullName: "Front Hazard Avoidance Camera - Left",
		Slug:     "front_hazcam_left_a",
	}

	// RoverCameraFRONTHAZCAMRIGHTA is the Front Hazard Avoidance Camera (Right A).
	RoverCameraFRONTHAZCAMRIGHTA = RoverCamera{
		Name:     "FRONT_HAZCAM_RIGHT_A",
		FullName: "Front Hazard Avoidance Camera - Right",
		Slug:     "front_hazcam_right_a",
	}

	// RoverCameraREARHAZCAMLEFT is the Rear Hazard Avoidance Camera (Left).
	RoverCameraREARHAZCAMLEFT = RoverCamera{
		Name:     "REAR_HAZCAM_LEFT",
		FullName: "Rear Hazard Avoidance Camera - Left",
		Slug:     "rear_hazcam_left",
	}

	// RoverCameraREARHAZCAMRIGHT is the Rear Hazard Avoidance Camera (Right).
	RoverCameraREARHAZCAMRIGHT = RoverCamera{
		Name:     "REAR_HAZCAM_RIGHT",
		FullName: "Rear Hazard Avoidance Camera - Right",
		Slug:     "rear_hazcam_right",
	}

	// RoverCameraSKYCAM is the SkyCam instrument.
	RoverCameraSKYCAM = RoverCamera{
		Name:     "SKYCAM",
		FullName: "SkyCam",
		Slug:     "skycam",
	}

	// RoverCameraSHERLOCWATSON is the SHERLOC WATSON instrument.
	RoverCameraSHERLOCWATSON = RoverCamera{
		Name:     "SHERLOC_WATSON",
		FullName: "SHERLOC WATSON",
		Slug:     "sherloc_watson",
	}

	// RoverCameraSUPERCAMRMI is the SuperCam Remote Micro-Imager.
	RoverCameraSUPERCAMRMI = RoverCamera{
		Name:     "SUPERCAM_RMI",
		FullName: "SuperCam Remote Micro-Imager",
		Slug:     "supercam_rmi",
	}

	// RoverCameraLCAM is the Lander Vision System Camera.
	RoverCameraLCAM = RoverCamera{
		Name:     "LCAM",
		FullName: "Lander Vision System Camera",
		Slug:     "lcam",
	}

	// RoverCameraMEDA is the Mars Environmental Dynamics Analyzer.
	RoverCameraMEDA = RoverCamera{
		Name:     "MEDA",
		FullName: "Mars Environmental Dynamics Analyzer",
		Slug:     "meda",
	}

	// RoverCameras is an easily iteratable array of cameras.
	RoverCameras = []RoverCamera{
		RoverCameraFHAZ,
		RoverCameraRHAZ,
		RoverCameraMAST,
		RoverCameraCHEMCAM,
		RoverCameraMAHLI,
		RoverCameraMARDI,
		RoverCameraNAVCAM,
		RoverCameraPANCAM,
		RoverCameraMINITES,
		RoverCameraEDLRUCAM,
		RoverCameraEDLRDCAM,
		RoverCameraEDLDDCAM,
		RoverCameraEDLPUCAM1,
		RoverCameraEDLPUCAM2,
		RoverCameraNAVCAMLEFT,
		RoverCameraNAVCAMRIGHT,
		RoverCameraMCZLEFT,
		RoverCameraMCZRIGHT,
		RoverCameraFRONTHAZCAMLEFTA,
		RoverCameraFRONTHAZCAMRIGHTA,
		RoverCameraREARHAZCAMLEFT,
		RoverCameraREARHAZCAMRIGHT,
		RoverCameraSKYCAM,
		RoverCameraSHERLOCWATSON,
		RoverCameraSUPERCAMRMI,
		RoverCameraLCAM,
		RoverCameraMEDA,
	}
)
