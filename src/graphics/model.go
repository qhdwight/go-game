package graphics

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Model struct {
	verts, norms []float32
	vertArrayHandle, normArrayHandle uint32
}

func (model *Model) Bind() {
	makeBuffer(model.verts)
	makeBuffer(model.norms)
	model.vertArrayHandle = makeVertexArray()
	model.normArrayHandle = makeVertexArray()
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

func bindVertexArrayToBuffer(index uint32, stride int32, pointer unsafe.Pointer) {
	gl.VertexAttribPointer(index, 3, gl.FLOAT, false, stride, pointer)
	gl.EnableVertexAttribArray(index)
}
