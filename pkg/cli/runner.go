package cli

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/suzuki-shunsuke/yodoc/pkg/controller/run"
	"github.com/urfave/cli/v3"
)

func Run(version string) int {
	logger := slogutil.New(&slogutil.InputNew{
		Name:    "yodoc",
		Version: version,
	})
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	runner := &Runner{}
	if err := runner.Run(ctx, logger, &urfave.Env{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}); err != nil {
		slogerr.WithError(logger.Logger, err).Error(errMsg(err))
		return 1
	}
	return 0
}

type Runner struct{}

func (r *Runner) Run(ctx context.Context, logger *slogutil.Logger, env *urfave.Env) error {
	return urfave.Command(env, &cli.Command{ //nolint:wrapcheck
		Name:  "yodoc",
		Usage: "Test command results and embed them into document",
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
		Commands: []*cli.Command{
			(&initCommand{}).command(logger),
			(&runCommand{}).command(logger),
		},
	}).Run(ctx, env.Args)
}

func errMsg(err error) string {
	ce := &run.CommandError{}
	msg := "yodoc failed"
	if errors.As(err, &ce) { //nolint:nestif
		if ce.Command != "" {
			msg += "\n" + "command:\n" + ce.Command
		}
		if ce.CombinedOutput != "" {
			msg += "\n" + "output:\n" + ce.CombinedOutput
		}
		if ce.Checks != "" {
			msg += "\n" + "checks:\n" + ce.Checks
		}
		if ce.Expr != "" {
			msg += "\n" + "expr: " + ce.Expr
		}
		if ce.Start != 0 && ce.End != 0 {
			msg += "\n" + "line number: " + strconv.Itoa(ce.Start) + " ~ " + strconv.Itoa(ce.End)
		}
	}
	return msg
}
