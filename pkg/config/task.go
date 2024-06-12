package config

import "github.com/spf13/afero"

type Task struct {
	Name   string
	Action *Action
	Before *Action
	After  *Action
	Checks []*Check
}

func (t *Task) SetDir(dir string) {
	t.Action.SetDir(dir)
	t.Before.SetDir(dir)
	t.After.SetDir(dir)
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
