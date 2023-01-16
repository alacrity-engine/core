package render

import (
	"fmt"
	"text/template"
)

type BatchOption func(batch *Batch, params *batchParameters) error

type batchParameters struct {
	initialObjectCapacity int
}

func BatchOptionWithVertexShaderTemplate(tmpl *template.Template) BatchOption {
	return func(batch *Batch, params *batchParameters) error {
		batch.vertexShaderTemplate = tmpl
		return nil
	}
}

func BatchOptionWithFragmentShader(fragmentShader *Shader) BatchOption {
	return func(batch *Batch, params *batchParameters) error {
		batch.fragmentShader = fragmentShader
		return nil
	}
}

func BatchOptionWithShaderProgram(shaderProgram *ShaderProgram, tmpl *template.Template) BatchOption {
	return func(batch *Batch, params *batchParameters) error {
		batch.shaderProgram = shaderProgram
		batch.vertexShaderTemplate = tmpl

		return nil
	}
}

func BatchOptionWithMaxNumCanvases(maxNumCanvases int) BatchOption {
	return func(batch *Batch, params *batchParameters) error {
		batch.maxNumCanvases = maxNumCanvases
		return nil
	}
}

func BatchOptionWithInitialObjectCapacity(capacity int) BatchOption {
	return func(batch *Batch, params *batchParameters) error {
		if capacity < 0 {
			return fmt.Errorf(
				"wrong capacity value: '%d'", capacity)
		}

		params.initialObjectCapacity = capacity
		return nil
	}
}
