package render

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type ShaderType uint32

const (
	ShaderTypeVertex   ShaderType = gl.VERTEX_SHADER
	ShaderTypeFragment ShaderType = gl.FRAGMENT_SHADER
)

var (
	//go:embed std-sprite-vert.glsl
	standardSpriteVertexShaderSource string
	//go:embed std-sprite-frag.glsl
	standardSpriteFragmentShaderSource string
)

type Shader struct {
	glHandler uint32
	typ       ShaderType
}

func (shader *Shader) Type() ShaderType {
	return shader.typ
}

func (shader *Shader) Delete() {
	gl.DeleteShader(shader.glHandler)
}

func NewShaderFromSource(source string, typ ShaderType) (*Shader, error) {
	shaderHandler := gl.CreateShader(uint32(typ))
	csources, free := gl.Strs(source)

	gl.ShaderSource(shaderHandler, 1, csources, nil)
	free()
	gl.CompileShader(shaderHandler)

	var status int32
	gl.GetShaderiv(shaderHandler, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderHandler, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderHandler, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return &Shader{
		glHandler: shaderHandler,
		typ:       typ,
	}, nil
}

func NewStandardSpriteShader(typ ShaderType) (*Shader, error) {
	switch typ {
	case ShaderTypeVertex:
		return NewShaderFromSource(standardSpriteVertexShaderSource, typ)

	case ShaderTypeFragment:
		return NewShaderFromSource(standardSpriteFragmentShaderSource, typ)

	default:
		return nil, fmt.Errorf("incorrect shader type: '%v'", typ)
	}
}
