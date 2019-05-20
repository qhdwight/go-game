package graphics

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Shader struct {
	Handle uint32
}

type Program struct {
	Handle uint32
}

func NewProgram(shaders ...*Shader) (*Program, error) {
	handle := gl.CreateProgram()
	program := Program{Handle: handle}
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
	return &program, nil
}

func CompileShaderFromSource(path string, shaderType uint32) (*Shader, error) {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return compileShader(string(src), shaderType)
}

func GetUniformLocation(name string, prog uint32) (int32, error) {
	location := gl.GetUniformLocation(prog, gl.Str(name+"\x00"))
	if location == -1 {
		return -1, fmt.Errorf("could not find location for uniform: %v", name)
	}
	return location, nil
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
