package main

import (
	"entities"
	"graphics"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

const (
	width    = 640
	height   = 480
	lookSens = 0.005
	speed    = 5
	fov      = 70
	vertPath = "lit.vert"
	fragPath = "lit.frag"
)

var (
	cube = []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		-0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,

		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,

		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, 0.5, -1.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,

		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
	}
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := initGlfw(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	window, err := initWindow()
	if err != nil {
		panic(err)
	}
	if err := gl.Init(); err != nil {
		panic(err)
	}
	makeBuffer(cube)
	vertArray := makeVertexArray()
	normArray := makeVertexArray()
	bindVertexArrayToBuffer(0, 6*4, nil)
	bindVertexArrayToBuffer(1, 6*4, gl.PtrOffset(3*4))
	//defer gl.DeleteVertexArrays(1, &vao)
	//defer gl.DeleteBuffers(1, &vbo)
	vertShader, err := graphics.CompileShaderFromPath(vertPath, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	defer gl.DeleteShader(vertShader.Handle)
	fragShader, err := graphics.CompileShaderFromPath(fragPath, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	defer gl.DeleteShader(fragShader.Handle)
	prog, err := graphics.NewProgram(vertShader, fragShader)
	if err != nil {
		panic(err)
	}
	defer gl.DeleteProgram(prog.Handle)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	modelUniform, err := prog.GetUniformLocation("model")
	viewUniform, err := prog.GetUniformLocation("view")
	projUniform, err := prog.GetUniformLocation("projection")
	lightPosUniform, err := prog.GetUniformLocation("lightPos")
	viewPosUniform, err := prog.GetUniformLocation("viewPos")
	lightColorUniform, err := prog.GetUniformLocation("lightColor")
	objectColorUniform, err := prog.GetUniformLocation("objectColor")
	player := entities.Player{Entity: entities.Entity{Pos: mgl64.Vec3{4, 3, 3}}}
	cam := graphics.MakeCamera(fov, width, height)
	lastRenderTime := 0.0
	for !window.ShouldClose() {
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
		now := glfw.GetTime()
		delta := now - lastRenderTime
		lastRenderTime = now
		mouseX, mouseY := window.GetCursorPos()
		window.SetCursorPos(width/2.0, height/2.0)
		player.Entity.AddInput((mouseX-width/2.0)*lookSens, (height/2.0-mouseY)*lookSens)
		fwd, right, up := entities.CalcRelVecs(player.Entity.Pitch, player.Entity.Yaw)
		moveOpt := func(key glfw.Key, vec mgl64.Vec3) {
			if window.GetKey(key) == glfw.Press {
				player.Entity.Pos = player.Entity.Pos.Add(vec.Mul(delta * speed))
			}
		}
		moveOpt(glfw.KeyW, fwd)
		moveOpt(glfw.KeyS, fwd.Mul(-1))
		moveOpt(glfw.KeyD, right)
		moveOpt(glfw.KeyA, right.Mul(-1))
		worldUp := mgl64.Vec3{0, 1, 0}
		moveOpt(glfw.KeyLeftShift, worldUp)
		moveOpt(glfw.KeyLeftControl, worldUp.Mul(-1))
		viewMat := mgl64.LookAtV(player.Entity.Pos, player.Entity.Pos.Add(fwd), up)
		modelMat := mgl64.Ident4()
		projMat32 := graphics.Mat64to32(cam.ProjMat)
		viewMat32 := graphics.Mat64to32(viewMat)
		modelMat32 := graphics.Mat64to32(modelMat)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog.Handle)
		gl.UniformMatrix4fv(projUniform, 1, false, &projMat32[0])
		gl.UniformMatrix4fv(viewUniform, 1, false, &viewMat32[0])
		gl.UniformMatrix4fv(modelUniform, 1, false, &modelMat32[0])
		gl.Uniform3f(lightPosUniform, 1.2, 1.4, 1.2)
		gl.Uniform3f(viewPosUniform, (float32)(player.Entity.Pos.X()), (float32)(player.Entity.Pos.Y()), (float32)(player.Entity.Pos.Z()))
		gl.Uniform3f(lightColorUniform, 1, 1, 1)
		gl.Uniform3f(objectColorUniform, 1, 0.2, 0.2)
		gl.BindVertexArray(vertArray)
		gl.BindVertexArray(normArray)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cube)/3))
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func initGlfw() error {
	if err := glfw.Init(); err != nil {
		return err
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	return nil
}

func initWindow() (*glfw.Window, error) {
	window, err := glfw.CreateWindow(width, height, "test", nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()
	window.SetInputMode(glfw.StickyKeysMode, gl.TRUE)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	//glfw.SwapInterval(0)
	return window, nil
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
