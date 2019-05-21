package graphics

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Camera struct {
	fov     float64
	ProjMat mgl64.Mat4
}

func MakeCamera(fov, width, height float64) Camera {
	camera := Camera{fov: fov}
	camera.ProjMat = mgl64.Perspective(mgl64.DegToRad(fov), width/height, 0.1, 100)
	return camera
}
