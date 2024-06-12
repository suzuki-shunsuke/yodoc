package render

import (
	"context"
	"errors"
	"fmt"

	"github.com/expr-lang/expr"
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

func (t *Task) run(act *config.Action) (*CommandResult, error) {
	envs := t.envs
	if len(act.Env) != 0 {
		envs = append(envs, act.GetEnv()...)
	}
	c := NewCommand(t.ctx, act.Shell, act.GetDir(), envs)
	if act.Run != "" {
		return c.Run(act.Run), nil
	}
	shell := act.Shell
	if shell == nil {
		shell = []string{"sh"}
	}
	c.Shell = shell
	result := c.Run(act.ScriptPath)
	result.Command = act.GetScript()
	return result, nil
}

func (t *Task) Check(cr *CommandResult, check *config.Check, task *config.Task) error {
	prog := check.GetExpr()
	output, err := expr.Run(prog, cr)
	if err != nil {
		return fmt.Errorf("evaluate an expression: %w", err)
	}
	b, ok := output.(bool)
	if !ok {
		return errors.New("the result of the expression isn't a boolean value")
	}
	if b {
		return nil
	}
	return errors.New("a check is false")
}

func (t *Task) Run(taskName string) (*CommandResult, error) {
	task, ok := t.tasks[taskName]
	if !ok {
		return nil, errors.New("task not found")
	}

	if task.Before != nil {
		if cr, err := t.run(task.Before); err != nil {
			return nil, err
		} else if cr.RunError != nil {
			return cr, cr.RunError
		} else if cr.ExitCode != 0 {
			return cr, errors.New("command failed")
		}
	}

	cr, err := t.run(task.Action)
	if err != nil {
		return nil, err
	}
	if cr.RunError != nil {
		return cr, cr.RunError
	}

	if task.After != nil {
		if cr, err := t.run(task.After); err != nil {
			return nil, err
		} else if cr.RunError != nil {
			return cr, cr.RunError
		} else if cr.ExitCode != 0 {
			return cr, errors.New("command failed")
		}
	}

	for _, check := range task.Checks {
		if err := t.Check(cr, check, task); err != nil {
			return nil, err
		}
	}

	return cr, nil
}
