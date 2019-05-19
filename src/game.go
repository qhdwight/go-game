package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"io/ioutil"
	"math"
	"runtime"
	"strings"
)

const (
	width  = 640
	height = 480
)

type player struct {
	pos        mgl64.Vec3
	xRot, yRot float64
}

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	player := player{pos: mgl64.Vec3{4, 3, 0}}
	projMat := mgl64.Perspective(mgl64.DegToRad(60), width/height, 0.1, 100)
	dir := mgl64.Vec3{
		math.Cos(player.xRot) * math.Sin(player.yRot),
		math.Sin(player.xRot),
		math.Cos(player.xRot) * math.Cos(player.yRot),
	}
	right := mgl64.Vec3{
		math.Sin(player.xRot - math.Pi/2),
		0,
		math.Cos(player.xRot - math.Pi/2),
	}
	up := right.Cross(dir)
	viewMat := mgl64.LookAtV(player.pos, dir, up)
	modelMat := mgl64.Ident4()
	mvpMat := projMat.Mul4(viewMat.Mul4(modelMat))
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	//glfw.SwapInterval(0)
	window, err := glfw.CreateWindow(width, height, "test", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}
	triangles := []float32{
		0, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}
	vertShaderSource, err := readShaderSource("resources/shaders/triangle.vert")
	if err != nil {
		panic(err)
	}
	fragShaderSource, err := readShaderSource("resources/shaders/triangle.frag")
	if err != nil {
		panic(err)
	}
	vao := makeVao(triangles)
	vertShader, err := compileShader(vertShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragShader, err := compileShader(fragShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertShader)
	gl.AttachShader(prog, fragShader)
	gl.LinkProgram(prog)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	lastRenderTime := 0.0
	for !window.ShouldClose() {
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
		now := glfw.GetTime()
		delta := lastRenderTime - now
		lastRenderTime = now
		//mouseX, mouseY = window.GetCursorPos()
		//window.SetCursorPos(width / 2.0, height / 2.0)
		//player.xRot += width / 2.0 - mouseX
		//player.yRot += height / 2.0 - mouseY
		if window.GetKey(glfw.KeyUp) == glfw.Press {
			player.pos.Add(dir.Mul(delta))
		}
		if window.GetKey(glfw.KeyDown) == glfw.Press {
			player.pos.Sub(dir.Mul(delta))
		}
		if window.GetKey(glfw.KeyRight) == glfw.Press {
			player.pos.Add(right.Mul(delta))
		}
		if window.GetKey(glfw.KeyLeft) == glfw.Press {
			player.pos.Sub(right.Mul(delta))
		}
		var mvpMat32 mgl32.Mat4
		//for i, f64 := range mvpMat {
		//	mvpMat32[i] = float32(f64)
		//}
		for i := 0; i < 16; i++ {
			mvpMat32[i] = float32(mvpMat[i])

		}
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog)
		gl.UniformMatrix4fv(shaderMvp, 1, false, &mvpMat32[0])
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangles)/3))
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func getUniformLocation(name string) (location, err) {

}

func compileShader(src string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csrc, csrcFree := gl.Strs(src + "\x00")
	gl.ShaderSource(shader, 1, csrc, nil)
	csrcFree()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile:\n%v\n%v", src, log)
	}
	return shader, nil
}

func readShaderSource(path string) (string, error) {
	src, err := ioutil.ReadFile(path)
	return string(src), err
}

func makeVao(points []float32) uint32 {
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
	return vao
}
