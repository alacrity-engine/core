package render

import (
	"fmt"
)

type BatchOption func(batch *Batch, params *batchParameters) error

type batchParameters struct {
	vertexShader          *Shader
	fragmentShader        *Shader
	initialObjectCapacity int
}

func BatchOptionWithVertexShader(vertexShader *Shader) BatchOption {
	return func(batch *Batch, params *batchParameters) error {
		params.vertexShader = vertexShader
		return nil
	}
}

func BatchOptionWithFragmentShader(fragmentShader *Shader) BatchOption {
	return func(batch *Batch, params *batchParameters) error {
		params.fragmentShader = fragmentShader
		return nil
	}
}

func BatchOptionWithShaderProgram(shaderProgram *ShaderProgram) BatchOption {
	return func(batch *Batch, params *batchParameters) error {
		batch.shaderProgram = shaderProgram
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
