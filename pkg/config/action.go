package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

type Action struct {
	Shell      []string
	Run        string
	script     string
	Script     string
	ScriptPath string `yaml:"-"`
	dir        string
	Dir        string
	Env        map[string]string
	env        []string
}

func (a *Action) GetScript() string {
	if a.Run != "" {
		return a.Run
	}
	return a.script
}

func (a *Action) SetEnv(env []string) {
	if a == nil {
		return
	}
	a.env = env
}

func (a *Action) GetEnv() []string {
	if a == nil {
		return nil
	}
	return a.env
}

func (a *Action) SetDir(dir string) {
	if a == nil {
		return
	}
	a.dir = filepath.Join(dir, a.Dir)
	a.ScriptPath = filepath.Join(dir, a.Script)
}

func (a *Action) GetDir() string {
	if a == nil {
		return ""
	}
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
