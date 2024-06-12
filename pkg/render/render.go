package render

import (
	"fmt"
	"text/template"

	"github.com/spf13/afero"
)

type Renderer struct {
	fs    afero.Fs
	funcs map[string]any
}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (r *Renderer) Render(src, dest string) error {
	srcByte, err := afero.ReadFile(r.fs, src)
	if err != nil {
		return fmt.Errorf("open a template file: %w", err)
	}

	destFile, err := r.fs.Create(dest)
	if err != nil {
		return fmt.Errorf("create a dest file: %w", err)
	}
	defer destFile.Close()

	tpl, err := template.New("_").Funcs(r.funcs).Parse(string(srcByte))
	if err != nil {
		return fmt.Errorf("parse a template: %w", err)
	}

	if err := tpl.Execute(destFile, nil); err != nil {
		return fmt.Errorf("execute a template: %w", err)
	}
	return nil
}
