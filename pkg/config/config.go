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

type Task struct {
	Name             string
	Shell            []string
	Run              string
	Script           string
	ScriptPath       string
	script           string
	dir              string
	Dir              string
	Env              map[string]string
	env              []string
	BeforeScript     string `yaml:"before_script"`
	AfterScript      string `yaml:"after_script"`
	beforeScript     string
	afterScript      string
	BeforeScriptPath string
	AfterScriptPath  string
	Checks           []*Check
}

type Check struct {
	Expr string
}

func (t *Task) GetScript() string {
	if t.Run != "" {
		return t.Run
	}
	return t.script
}

func (t *Task) SetEnv() {
	envs := make([]string, 0, len(t.Env))
	for k, v := range t.Env {
		envs = append(envs, k+"="+v)
	}
	t.env = envs
}

func (t *Task) GetEnv() []string {
	return t.env
}

func (t *Task) SetDir(dir string) {
	t.dir = filepath.Join(dir, t.Dir)
	t.ScriptPath = filepath.Join(dir, t.Script)
	t.BeforeScriptPath = filepath.Join(dir, t.BeforeScript)
	t.AfterScriptPath = filepath.Join(dir, t.AfterScript)
}

func (t *Task) ReadScript(fs afero.Fs) error {
	if t.Script != "" {
		if b, err := afero.ReadFile(fs, t.ScriptPath); err != nil {
			return fmt.Errorf("read a script file: %w", err)
		} else {
			t.script = string(b)
		}
	}
	if t.AfterScript != "" {
		if b, err := afero.ReadFile(fs, t.AfterScriptPath); err != nil {
			return fmt.Errorf("read an after script file: %w", err)
		} else {
			t.afterScript = string(b)
		}
	}
	if t.BeforeScript != "" {
		if b, err := afero.ReadFile(fs, t.BeforeScriptPath); err != nil {
			return fmt.Errorf("read a before script file: %w", err)
		} else {
			t.beforeScript = string(b)
		}
	}
	return nil
}
