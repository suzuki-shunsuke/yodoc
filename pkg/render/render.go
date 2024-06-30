package render

import (
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/yodoc/pkg/frontmatter"
)

type Renderer struct {
	fs           afero.Fs
	leftDelim    string
	rightDelim   string
	funcs        map[string]any
	funcsWithEnv map[string]any
}

func NewRenderer(fs afero.Fs) *Renderer {
	fncs := sprig.TxtFuncMap()
	delete(fncs, "env")
	delete(fncs, "expandenv")
	delete(fncs, "getHostByName")

	fncsWithEnv := sprig.TxtFuncMap()
	delete(fncsWithEnv, "expandenv")
	delete(fncsWithEnv, "getHostByName")

	return &Renderer{
		fs:           fs,
		funcs:        fncs,
		funcsWithEnv: fncsWithEnv,
	}
}

func (r *Renderer) SetDelims(left, right string) {
	r.leftDelim = left
	r.rightDelim = right
}

func (r *Renderer) NewTemplate() *template.Template {
	return template.New("_").Funcs(r.funcs)
}

func (r *Renderer) NewTemplateWithEnv() *template.Template {
	return template.New("_").Funcs(r.funcsWithEnv)
}

func (r *Renderer) RenderFile(src, dest, txt string, delim *frontmatter.Delim) error {
	destFile, err := r.fs.Create(dest)
	if err != nil {
		return fmt.Errorf("create a dest file: %w", err)
	}
	defer destFile.Close()

	tpl := r.NewTemplate().Funcs(Funcs(r.fs, src))

	r.setDelim(tpl, delim)

	tpl, err = tpl.Parse(txt)
	if err != nil {
		return fmt.Errorf("parse a template: %w", err)
	}

	if err := tpl.Execute(destFile, nil); err != nil {
		return fmt.Errorf("execute a template: %w", err)
	}

	if _, err := destFile.WriteString(Footer); err != nil {
		return fmt.Errorf("write a dest file: %w", err)
	}

	return nil
}

func (r *Renderer) setDelim(tpl *template.Template, delim *frontmatter.Delim) {
	leftDelim := r.leftDelim
	rightDelim := r.rightDelim
	if delim != nil {
		if delim.Left != "" {
			leftDelim = delim.Left
		}
		if delim.Right != "" {
			leftDelim = delim.Right
		}
	}
	tpl.Delims(leftDelim, rightDelim)
}
