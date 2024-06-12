package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

type Config struct {
	Src   string
	Dest  string
	Delim *Delim
	Tasks []*Task
}

type Delim struct {
	Left  string
	Right string
}

func (d *Delim) GetLeft() string {
	if d == nil {
		return ""
	}
	return d.Left
}

func (d *Delim) GetRight() string {
	if d == nil {
		return ""
	}
	return d.Right
}

type Action struct {
	Shell      []string
	Run        string
	script     string
	Script     string
	ScriptPath string
	dir        string
	Dir        string
	Env        map[string]string
	env        []string
}

type Task struct {
	Name   string
	Action *Action
	Before *Action
	After  *Action
	Checks []*Check
}

type Check struct {
	Expr string
}

func (a *Action) GetScript() string {
	if a.Run != "" {
		return a.Run
	}
	return a.script
}

func (t *Task) SetEnv() {
	t.Action.SetEnv()
	t.After.SetEnv()
	t.Before.SetEnv()
}

func (a *Action) SetEnv() {
	if a == nil {
		return
	}
	envs := make([]string, 0, len(a.Env))
	for k, v := range a.Env {
		envs = append(envs, k+"="+v)
	}
	a.env = envs
}

func (a *Action) GetEnv() []string {
	return a.env
}

func (t *Task) SetDir(dir string) {
	t.Action.SetDir(dir)
	t.Before.SetDir(dir)
	t.After.SetDir(dir)
}

func (a *Action) SetDir(dir string) {
	if a == nil {
		return
	}
	a.dir = filepath.Join(dir, a.Dir)
	a.ScriptPath = filepath.Join(dir, a.Script)
}

func (a *Action) GetDir() string {
	return a.dir
}

func (a *Action) ReadScript(fs afero.Fs) error {
	if a == nil || a.script == "" {
		return nil
	}
	b, err := afero.ReadFile(fs, a.ScriptPath)
	if err != nil {
		return fmt.Errorf("read a script file: %w", err)
	}
	a.script = string(b)
	return nil
}

func (t *Task) ReadScript(fs afero.Fs) error {
	if err := t.Action.ReadScript(fs); err != nil {
		return err
	}
	if err := t.Before.ReadScript(fs); err != nil {
		return err
	}
	if err := t.After.ReadScript(fs); err != nil {
		return err
	}
	return nil
}
