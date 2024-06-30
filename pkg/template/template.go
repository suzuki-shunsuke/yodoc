package template

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

func RenderTemplateToString(tpl *template.Template, data any) (string, error) {
	b := &bytes.Buffer{}
	if err := tpl.Execute(b, data); err != nil {
		return "", fmt.Errorf("render a template: %w", err)
	}
	return b.String(), nil
}

func RenderTemplate(tpl *template.Template, data any, out io.Writer) error {
	if err := tpl.Execute(out, data); err != nil {
		return fmt.Errorf("render a template: %w", err)
	}
	return nil
}
