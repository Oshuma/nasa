package nasa

import (
	"testing"
)

func TestHasCamera(t *testing.T) {
	t.Run("has camera", func(t *testing.T) {
		r := RoverCuriosity
		c := RoverCameraFHAZ
		if !hasCamera(r, c) {
			t.Errorf("rover %s should have camera %s", r.Name, c.Name)
		}
	})

	t.Run("perseverance has camera", func(t *testing.T) {
		r := RoverPerseverance
		c := RoverCameraMCZLEFT
		if !hasCamera(r, c) {
			t.Errorf("rover %s should have camera %s", r.Name, c.Name)
		}
	})

	t.Run("does not have camera", func(t *testing.T) {
		r := RoverCuriosity
		c := RoverCameraPANCAM
		if hasCamera(r, c) {
			t.Errorf("rover %s should not have camera %s", r.Name, c.Name)
		}
	})

	t.Run("perseverance missing camera", func(t *testing.T) {
		r := RoverPerseverance
		c := RoverCameraMAST
		if hasCamera(r, c) {
			t.Errorf("rover %s should not have camera %s", r.Name, c.Name)
		}
	})
}

func TestMarsParamsMismatch(t *testing.T) {
	p := &APIParam{APIKey: "key"}

	if _, err := MarsRoverPhotos(p, RoverCuriosity); err != ErrorParamsMismatch {
		t.Errorf("MarsRoverPhotos: expected ErrorParamsMismatch, got: %v", err)
	}

	if _, err := MarsRoverPhotosLatest(p, RoverCuriosity); err != ErrorParamsMismatch {
		t.Errorf("MarsRoverPhotosLatest: expected ErrorParamsMismatch, got: %v", err)
	}

	if _, err := MarsMissionManifest(p, RoverCuriosity); err != ErrorParamsMismatch {
		t.Errorf("MarsMissionManifest: expected ErrorParamsMismatch, got: %v", err)
	}
}
