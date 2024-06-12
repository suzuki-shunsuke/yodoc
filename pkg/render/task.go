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
	c := NewCommand(t.ctx, task.Shell, t.dir, t.envs)
	return c.Run(task.Run), nil
}
