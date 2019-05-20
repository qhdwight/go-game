package entities

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Entity struct {
	Pos        mgl64.Vec3
	Yaw, Pitch float64
}

type Player struct {
	Entity Entity
}

func CalcRelVecs(pitch, yaw float64) (mgl64.Vec3, mgl64.Vec3, mgl64.Vec3) {
	fwd := mgl64.Vec3{
		math.Cos(pitch) * math.Cos(yaw),
		math.Sin(pitch),
		math.Cos(pitch) * math.Sin(yaw),
	}
	right := fwd.Cross(mgl64.Vec3{0, 1, 0})
	up := right.Cross(fwd)
	return fwd, right, up
}

func (entity *Entity) AddInput(mouseX, mouseY float64) {
	entity.Yaw = math.Mod(entity.Yaw+mouseX, math.Pi*2)
	entity.Pitch = mgl64.Clamp(entity.Pitch+mouseY, -math.Pi/2, math.Pi/2)
}
