package render

import (
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/afero"
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
