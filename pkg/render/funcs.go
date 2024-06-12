package render

import (
	"bytes"
	"context"
	"io"
	"os/exec"
)

type Funcs struct {
	ctx        context.Context //nolint:containedctx
	dir        string
	envs       []string
	appendEnvs []string
	environ    func() []string
	tasks      map[string]*Task
}

type Task struct {
	Name string
}

func (f *Funcs) Funcs() map[string]any {
	return map[string]any{
		"Command": f.Command,
		// "Task": f.Task,
	}
}

type CommandResult struct {
	ExitCode       int
	Stdout         func() string
	Stderr         func() string
	CombinedOutput func() string
	RunError       error
}

func (f *Funcs) env() []string {
	if f.envs == nil {
		if len(f.appendEnvs) == 0 {
			return nil
		}
		return append(f.environ(), f.appendEnvs...)
	}
	if len(f.appendEnvs) == 0 {
		return f.envs
	}
	return append(f.envs, f.appendEnvs...)
}

func (f *Funcs) Command(s string) *CommandResult {
	cmd := exec.CommandContext(f.ctx, "sh", "-c", s)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	combinedOutput := &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(stdout, combinedOutput)
	cmd.Stderr = io.MultiWriter(stderr, combinedOutput)

	if f.dir != "" {
		cmd.Dir = f.dir
	}

	cmd.Env = f.env()

	if err := cmd.Run(); err != nil {
		return &CommandResult{
			RunError: err,
		}
	}
	return &CommandResult{
		Stdout:         stdout.String,
		Stderr:         stderr.String,
		CombinedOutput: combinedOutput.String,
		ExitCode:       cmd.ProcessState.ExitCode(),
	}
}
