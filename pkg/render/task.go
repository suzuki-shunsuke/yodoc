package render

import (
	"context"
	"errors"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
	"github.com/suzuki-shunsuke/yodoc/pkg/frontmatter"
)

type Task struct {
	ctx   context.Context //nolint:containedctx
	tasks map[string]*config.Task
	fm    *frontmatter.Frontmatter
}

func NewTask(ctx context.Context, tasks map[string]*config.Task, fm *frontmatter.Frontmatter) *Task {
	return &Task{
		ctx:   ctx,
		tasks: tasks,
		fm:    fm,
	}
}

func (t *Task) run(act *config.Action) *CommandResult {
	if t == nil {
		return nil
	}
	dir := act.GetDir()
	if dir == "" && t.fm != nil {
		dir = t.fm.Dir
	}
	c := NewCommand(t.ctx, act.Shell, dir, act.GetEnv())
	if act.Run != "" {
		return c.Run(act.Run)
	}
	shell := act.Shell
	if shell == nil {
		shell = []string{"sh"}
	}
	c.shell = shell
	result := c.Run(act.ScriptPath)
	result.Command = act.GetScript()
	return result
}

func (t *Task) Check(cr *CommandResult, check *config.Check) error {
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

func (t *Task) before(before *config.Action) error {
	cr := t.run(before)
	if cr.RunError != nil {
		return cr.RunError
	}
	if cr.ExitCode != 0 {
		return errors.New("command failed")
	}
	return nil
}

func (t *Task) after(after *config.Action) error {
	cr := t.run(after)
	if cr.RunError != nil {
		return cr.RunError
	}
	if cr.ExitCode != 0 {
		return errors.New("command failed")
	}
	return nil
}

func (t *Task) Run(taskName string) (*CommandResult, error) {
	task, ok := t.tasks[taskName]
	if !ok {
		return nil, errors.New("task not found")
	}

	if task.Before != nil {
		if err := t.before(task.Before); err != nil {
			return nil, err
		}
	}

	cr := t.run(task.Action)

	if task.After != nil {
		if err := t.after(task.After); err != nil {
			return nil, err
		}
	}

	for _, check := range task.Checks {
		if err := t.Check(cr, check); err != nil {
			return nil, err
		}
	}

	return cr, nil
}
