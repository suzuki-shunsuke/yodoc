package cli

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

type Runner struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	LDFlags *LDFlags
	LogE    *logrus.Entry
}

type LDFlags struct {
	Version string
	Commit  string
	Date    string
}

func (r *Runner) Run(ctx context.Context, args ...string) error {
	app := cli.Command{
		Name:    "yodoc",
		Usage:   "",
		Version: r.LDFlags.Version + " (" + r.LDFlags.Commit + ")",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "log level",
				Sources: cli.EnvVars("YODOC_LOG_LEVEL"),
			},
			&cli.StringFlag{
				Name:    "log-color",
				Usage:   "Log color. One of 'auto' (default), 'always', 'never'",
				Sources: cli.EnvVars("YODOC_LOG_COLOR"),
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Configuration file path",
				Sources: cli.EnvVars("YODOC_CONFIG"),
			},
		},
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			(&versionCommand{}).command(),
			(&initCommand{
				logE: r.LogE,
			}).command(),
			(&runCommand{
				logE: r.LogE,
			}).command(),
			(&completionCommand{
				logE:   r.LogE,
				stdout: r.Stdout,
			}).command(),
		},
	}

	return app.Run(ctx, args) //nolint:wrapcheck
}
