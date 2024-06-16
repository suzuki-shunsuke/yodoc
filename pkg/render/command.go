package render

import (
	"bytes"
	"context"
	"io"
	"os/exec"

	"github.com/suzuki-shunsuke/yodoc/pkg/command"
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
	shell []string
	dir   string
	envs  []string
}

func NewCommand(ctx context.Context, shell []string, dir string, envs []string) *Command {
	return &Command{
		ctx:   ctx,
		shell: shell,
		dir:   dir,
		envs:  envs,
	}
}

func (c *Command) Run(s string) *CommandResult {
	if len(c.shell) == 0 {
		c.shell = []string{"sh", "-c"}
	}
	cmd := exec.CommandContext(c.ctx, c.shell[0], append(c.shell[1:], s)...) //nolint:gosec
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	combinedOutput := &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(stdout, combinedOutput)
	cmd.Stderr = io.MultiWriter(stderr, combinedOutput)
	command.SetCancel(cmd)

	if c.dir != "" {
		cmd.Dir = c.dir
	}

	cmd.Env = c.envs

	if err := cmd.Run(); err != nil {
		return &CommandResult{
			Command:        s,
			Stdout:         stdout.String(),
			Stderr:         stderr.String(),
			CombinedOutput: combinedOutput.String(),
			ExitCode:       cmd.ProcessState.ExitCode(),
			RunError:       err,
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
