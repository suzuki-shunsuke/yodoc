package render

import (
	"bytes"
	"context"
	"io"
	"os/exec"
)

type CommandResult struct {
	Command        string
	ExitCode       int
	Stdout         string
	Stderr         string
	CombinedOutput string
	RunError       error
}

type Command struct {
	ctx   context.Context //nolint:containedctx
	Shell []string
	dir   string
	envs  []string
}

func NewCommand(ctx context.Context, shell []string, dir string, envs []string) *Command {
	return &Command{
		ctx:  ctx,
		dir:  dir,
		envs: envs,
	}
}

func (c *Command) Run(s string) *CommandResult {
	if c.Shell == nil {
		c.Shell = []string{"sh", "-c"}
	}
	cmd := exec.CommandContext(c.ctx, c.Shell[0], append(c.Shell[1:], s)...)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	combinedOutput := &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(stdout, combinedOutput)
	cmd.Stderr = io.MultiWriter(stderr, combinedOutput)

	if c.dir != "" {
		cmd.Dir = c.dir
	}

	cmd.Env = c.envs

	if err := cmd.Run(); err != nil {
		return &CommandResult{
			Command:  s,
			RunError: err,
		}
	}
	return &CommandResult{
		Command:        s,
		Stdout:         stdout.String(),
		Stderr:         stderr.String(),
		CombinedOutput: combinedOutput.String(),
		ExitCode:       cmd.ProcessState.ExitCode(),
	}
}

func envs(envs, appendEnvs []string, environ func() []string) []string {
	if envs == nil {
		if len(appendEnvs) == 0 {
			return nil
		}
		return append(environ(), appendEnvs...)
	}
	if len(appendEnvs) == 0 {
		return envs
	}
	return append(envs, appendEnvs...)
}
