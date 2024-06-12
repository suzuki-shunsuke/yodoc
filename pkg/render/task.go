package render

import (
	"context"
	"errors"

	"github.com/suzuki-shunsuke/yodoc/pkg/config"
)

type Task struct {
	ctx   context.Context //nolint:containedctx
	tasks map[string]*config.Task
	dir   string
	envs  []string
}

func NewTask(ctx context.Context, tasks map[string]*config.Task, dir string, envs []string) *Task {
	return &Task{
		ctx:   ctx,
		tasks: tasks,
		dir:   dir,
		envs:  envs,
	}
}

func (t *Task) Run(taskName string) (*CommandResult, error) {
	task, ok := t.tasks[taskName]
	if !ok {
		return nil, errors.New("task not found")
	}
	envs := t.envs
	if len(task.Action.Env) != 0 {
		envs = append(envs, task.Action.GetEnv()...)
	}
	c := NewCommand(t.ctx, task.Action.Shell, task.Action.GetDir(), envs)
	if task.Action.Run != "" {
		return c.Run(task.Action.Run), nil
	}
	shell := task.Action.Shell
	if shell == nil {
		shell = []string{"sh"}
	}
	c.Shell = shell
	result := c.Run(task.Action.ScriptPath)
	result.Command = task.Action.GetScript()
	return result, nil
}
