package render

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ShaderProgram struct {
	glHandler      uint32
	vertexShader   *Shader
	fragmentShader *Shader
}

func (program *ShaderProgram) Use() {
	gl.UseProgram(program.glHandler)
}

func (program *ShaderProgram) SetInt(name string, value int) {
	location := gl.GetUniformLocation(program.glHandler, gl.Str(name+"\x00"))
	gl.Uniform1i(location, int32(value))
}

func (program *ShaderProgram) SetFloat32(name string, value float32) {
	location := gl.GetUniformLocation(program.glHandler, gl.Str(name+"\x00"))
	gl.Uniform1f(location, value)
}

func (program *ShaderProgram) SetMatrix4(name string, value mgl32.Mat4) {
	location := gl.GetUniformLocation(program.glHandler, gl.Str(name+"\x00"))
	gl.UniformMatrix4fv(location, 1, false, &value[0])
}

func NewShaderProgramFromShaders(vertexShader, fragmentShader *Shader) (*ShaderProgram, error) {
	if vertexShader == nil || vertexShader.glHandler == 0 {
		return nil, fmt.Errorf("no vertex shader")
	}

	if fragmentShader == nil || fragmentShader.glHandler == 0 {
		return nil, fmt.Errorf("no fragment shader")
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader.glHandler)
	gl.AttachShader(program, fragmentShader.glHandler)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("failed to link program: %v", log)
	}

	return &ShaderProgram{
		glHandler:      program,
		vertexShader:   vertexShader,
		fragmentShader: fragmentShader,
	}, nil
}

func NewStandardSpriteShaderProgram() (*ShaderProgram, error) {
	vertexShader, err := NewStandardSpriteShader(ShaderTypeVertex)

	if err != nil {
		return nil, err
	}

	fragmentShader, err := NewStandardSpriteShader(ShaderTypeFragment)

	if err != nil {
		return nil, err
	}

	return NewShaderProgramFromShaders(vertexShader, fragmentShader)
}
