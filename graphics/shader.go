package graphics

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

const (
	shaderPath = "resources/shaders/"
)

type Shader struct {
	Handle uint32
}

type Program struct {
	Handle uint32
}

func NewProgram(shaders []*Shader) (*Program, error) {
	handle := gl.CreateProgram()
	program := &Program{Handle: handle}
	for _, shader := range shaders {
		gl.AttachShader(handle, shader.Handle)
	}
	gl.LinkProgram(handle)
	var isLinked int32
	gl.GetProgramiv(handle, gl.LINK_STATUS, &isLinked)
	if isLinked == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(handle, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(handle, logLength, nil, gl.Str(log))
		gl.DeleteProgram(handle)
		return nil, fmt.Errorf("failed to create program:\n%v", log)
	}
	for _, shader := range shaders {
		gl.DetachShader(handle, shader.Handle)
	}
	return program, nil
}

func CompileShaderFromPath(path string, shaderType uint32) (*Shader, error) {
	src, err := ioutil.ReadFile(shaderPath + path)
	if err != nil {
		return nil, err
	}
	return compileShader(string(src), shaderType)
}

func Mat64to32(mat64 mgl64.Mat4) mgl32.Mat4 {
	var mat32 mgl32.Mat4
	for i, f64 := range mat64 {
		mat32[i] = float32(f64)
	}
	return mat32
}

func Vec64to32(vec64 mgl64.Vec3) mgl32.Vec3 {
	var vec32 mgl32.Vec3
	for i, f64 := range vec64 {
		vec32[i] = float32(f64)
	}
	return vec32
}

func (program *Program) GetUniformLocation(name string) (int32, error) {
	location := gl.GetUniformLocation(program.Handle, gl.Str(name+"\x00"))
	if location == -1 {
		return -1, fmt.Errorf("could not find location for uniform: %v", name)
	}
	return location, nil
}

func SetUniformVec3(location int32, vec mgl64.Vec3) {
	vec32 := Vec64to32(vec)
	gl.Uniform3fv(location, 1, &vec32[0])
}

func SetUniformMat4(location int32, mat mgl64.Mat4) {
	mat32 := Mat64to32(mat)
	gl.UniformMatrix4fv(location, 1, false, &mat32[0])
}

func compileShader(src string, shaderType uint32) (*Shader, error) {
	handle := gl.CreateShader(shaderType)
	csrc, csrcFree := gl.Strs(src + "\x00")
	gl.ShaderSource(handle, 1, csrc, nil)
	csrcFree()
	gl.CompileShader(handle)
	var status int32
	gl.GetShaderiv(handle, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(handle, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(handle, logLength, nil, gl.Str(log))
		return nil, fmt.Errorf("failed to compile:\n%v\n%v", src, log)
	}
	return &Shader{Handle: handle}, nil
}
