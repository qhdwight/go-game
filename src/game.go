package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"io/ioutil"
	"runtime"
	"strings"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(640, 480, "test", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(0)
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
	vertShader, err := compileShader(string(vertShaderSource), gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragShader, err := compileShader(string(fragShaderSource), gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertShader)
	gl.AttachShader(prog, fragShader)
	gl.LinkProgram(prog)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangles)/3))
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil
}

func readShaderSource(path string) (string, error) {
	source, err := ioutil.ReadFile(path)
	//if source[len(source) - 1] != 0 {
	//	source = append(source, 0)
	//}
	return string(source), err
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
