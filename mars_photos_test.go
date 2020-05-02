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

	t.Run("does not have camera", func(t *testing.T) {
		r := RoverCuriosity
		c := RoverCameraPANCAM
		if hasCamera(r, c) {
			t.Errorf("rover %s should not have camera %s", r.Name, c.Name)
		}
	})
}
