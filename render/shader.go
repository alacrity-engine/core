package render

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type ShaderType uint32

const (
	ShaderTypeVertex   ShaderType = gl.VERTEX_SHADER
	ShaderTypeFragment ShaderType = gl.FRAGMENT_SHADER
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
	csources, free := gl.Strs(source + "\x00")

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

func NewStandardBatchShader(typ ShaderType, maxNumCanvases int) (*Shader, *template.Template, error) {
	if maxNumCanvases <= 0 {
		return nil, nil, fmt.Errorf(
			"wrong maxNumCanvases value: '%d'", maxNumCanvases)
	}

	switch typ {
	case ShaderTypeFragment:
		shader, err := NewShaderFromSource(batchFragmentShaderSource, typ)
		return shader, nil, err

	case ShaderTypeVertex:
		var strBuilder strings.Builder

		err := batchVertexShaderTemplate.
			Execute(&strBuilder, map[string]interface{}{
				"maxNumCanvases": maxNumCanvases,
			})

		if err != nil {
			return nil, nil, err
		}

		shaderSource := strBuilder.String()
		vertexShader, err := NewShaderFromSource(
			shaderSource, ShaderTypeVertex)

		if err != nil {
			return nil, nil, err
		}

		return vertexShader, batchVertexShaderTemplate, nil

	default:
		return nil, nil, fmt.Errorf("incorrect shader type: '%v'", typ)
	}
}

func NewBatchShaderWithTemplate(typ ShaderType, tmpl *template.Template, maxNumCanvases int) (*Shader, error) {
	if maxNumCanvases <= 0 {
		return nil, fmt.Errorf(
			"wrong maxNumCanvases value: '%d'", maxNumCanvases)
	}

	if tmpl == nil {
		return nil, fmt.Errorf("no shader template provided")
	}

	switch typ {
	case ShaderTypeFragment:
		shader, err := NewShaderFromSource(batchFragmentShaderSource, typ)
		return shader, err

	case ShaderTypeVertex:
		var strBuilder strings.Builder

		err := tmpl.Execute(&strBuilder, map[string]interface{}{
			"maxNumCanvases": maxNumCanvases,
		})

		if err != nil {
			return nil, err
		}

		shaderSource := strBuilder.String()
		vertexShader, err := NewShaderFromSource(
			shaderSource, ShaderTypeVertex)

		if err != nil {
			return nil, err
		}

		return vertexShader, nil

	default:
		return nil, fmt.Errorf("incorrect shader type: '%v'", typ)
	}
}
