package entities

import (
	graphics2 "biomquest/graphics"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

type Transform struct {
	Pos        mgl64.Vec3
	Yaw, Pitch float64
}

type WorldObject struct {
}

type Entity struct {
	Transform   Transform
	WorldObject WorldObject
}

type VisualEntity struct {
	Entity    Entity
	Model     *graphics2.Model
	Transform Transform
}

type Player struct {
	VisualEntity VisualEntity
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
	transform := &entity.Transform
	transform.Yaw = math.Mod(transform.Yaw+mouseX, math.Pi*2)
	transform.Pitch = mgl64.Clamp(transform.Pitch+mouseY, -math.Pi/2, math.Pi/2)
}
