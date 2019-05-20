package main

import (
	"entities"
	"graphics"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

const (
	width    = 640
	height   = 480
	lookSens = 0.005
	speed    = 5
	fov      = 70
	vertPath = "resources/shaders/basic.vert"
	fragPath = "resources/shaders/basic.frag"
)

var (
	triangles = []float32{
		-1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
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
	vbo, vao := makeVao(triangles)
	defer gl.DeleteVertexArrays(1, &vao)
	defer gl.DeleteBuffers(1, &vbo)
	vertShader, err := graphics.CompileShaderFromSource(vertPath, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	defer gl.DeleteShader(vertShader.Handle)
	fragShader, err := graphics.CompileShaderFromSource(fragPath, gl.FRAGMENT_SHADER)
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
	mvpUniform, err := graphics.GetUniformLocation("mvp", prog.Handle)
	if err != nil {
		panic(err)
	}
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
		mvpMat := cam.ProjMat.Mul4(viewMat).Mul4(modelMat)
		var mvpMat32 mgl32.Mat4
		//for i, f64 := range mvpMat {
		//	mvpMat32[i] = float32(f64)
		//}
		for i := 0; i < 16; i++ {
			mvpMat32[i] = float32(mvpMat[i])
		}
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog.Handle)
		gl.UniformMatrix4fv(mvpUniform, 1, false, &mvpMat32[0])
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangles)/3))
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

func makeVao(points []float32) (uint32, uint32) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	return vbo, vao
}
