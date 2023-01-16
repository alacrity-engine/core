package render

import (
	_ "embed"
	"text/template"
)

var (
	//go:embed std-batch-vert.tmpl.glsl
	batchVertexShaderTemplateSource string

	batchVertexShaderTemplate *template.Template
)

func InitTemplates() error {
	var err error
	batchVertexShaderTemplate, err = template.
		New("batchVertexShaderTemplate").
		Parse(batchVertexShaderTemplateSource)

	if err != nil {
		return err
	}

	return nil
}
