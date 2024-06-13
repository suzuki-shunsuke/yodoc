package render

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
)

type Renderer struct {
	fs           afero.Fs
	leftDelim    string
	rightDelim   string
	tasks        map[string]*config.Task
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

const header = `
<!-- This file is generated by yodoc.
https://github.com/suzuki-shunsuke/yodoc
Please don't edit this code comment because yodoc depends on this code comment.
-->
`

func (r *Renderer) SetDelims(left, right string) {
	r.leftDelim = left
	r.rightDelim = right
}

func (r *Renderer) SetTasks(tasks map[string]*config.Task) {
	r.tasks = tasks
}

func (r *Renderer) GetActionEnv(action *config.Action) ([]string, error) {
	if action == nil || len(action.Env) == 0 {
		return nil, nil
	}
	envs := make([]string, 0, len(action.Env))
	for k, v := range action.Env {
		e, err := r.renderEnv(v)
		if err != nil {
			return nil, fmt.Errorf("render an environment variable: %w", err)
		}
		envs = append(envs, k+"="+e)
	}
	return envs, nil
}

func (r *Renderer) renderEnv(v string) (string, error) {
	tpl, err := r.NewTemplateWithEnv().Parse(v)
	if err != nil {
		return "", fmt.Errorf("parse a template: %w", err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, nil); err != nil {
		return "", fmt.Errorf("evaluate an environment: %w", err)
	}
	return buf.String(), nil
}

func (r *Renderer) NewTemplate() *template.Template {
	return template.New("_").Funcs(r.funcs)
}

func (r *Renderer) NewTemplateWithEnv() *template.Template {
	return template.New("_").Funcs(r.funcsWithEnv)
}

func (r *Renderer) Render(ctx context.Context, src, dest string) error {
	srcByte, err := afero.ReadFile(r.fs, src)
	if err != nil {
		return fmt.Errorf("open a template file: %w", err)
	}

	destFile, err := r.fs.Create(dest)
	if err != nil {
		return fmt.Errorf("create a dest file: %w", err)
	}
	defer destFile.Close()

	tpl, err := r.NewTemplate().Funcs(Funcs(ctx, r.fs, r.tasks, src)).Parse(string(srcByte))
	if err != nil {
		return fmt.Errorf("parse a template: %w", err)
	}

	tpl.Delims(r.leftDelim, r.rightDelim)

	if err := tpl.Execute(destFile, nil); err != nil {
		return fmt.Errorf("execute a template: %w", err)
	}

	if _, err := destFile.WriteString(header); err != nil {
		return fmt.Errorf("write a dest file: %w", err)
	}

	return nil
}
