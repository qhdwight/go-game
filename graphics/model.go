package graphics

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"unsafe"
)

type Model struct {
	verts, norms                     []float32
	vertArrayHandle, normArrayHandle uint32
}

func (model *Model) GetNorms() []float32 {
	return model.norms
}

func NewModel(verts, norms []float32) *Model {
	model := &Model{verts: verts, norms: norms}
	vertBufHandle := makeBuffer(verts)
	normBufHandle := makeBuffer(norms)
	model.vertArrayHandle = makeVertexArray()
	model.normArrayHandle = makeVertexArray()
	bindVertexArrayToBuffer(0, vertBufHandle, 0, nil)
	bindVertexArrayToBuffer(1, normBufHandle, 0, nil)
	return model
}

func makeBuffer(data []float32) uint32 {
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(data), gl.Ptr(data), gl.STATIC_DRAW)
	return buffer
}

func makeVertexArray() uint32 {
	var vertexArray uint32
	gl.GenVertexArrays(1, &vertexArray)
	gl.BindVertexArray(vertexArray)
	return vertexArray
}

func (model *Model) BindVertexArray() {
	gl.BindVertexArray(model.vertArrayHandle)
	gl.BindVertexArray(model.normArrayHandle)
}

func bindVertexArrayToBuffer(index, buffer uint32, stride int32, pointer unsafe.Pointer) {
	gl.EnableVertexAttribArray(index)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.VertexAttribPointer(index, 3, gl.FLOAT, false, stride, pointer)
}
