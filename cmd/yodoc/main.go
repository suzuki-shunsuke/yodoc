package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/yodoc/pkg/cli"
	"github.com/suzuki-shunsuke/yodoc/pkg/controller/run"
)

var (
	version = ""
	commit  = "" //nolint:gochecknoglobals
	date    = "" //nolint:gochecknoglobals
)

func main() {
	if code := core(); code != 0 {
		os.Exit(code)
	}
}

func core() int {
	logLevelVar := &slog.LevelVar{}
	logger := slogutil.New(&slogutil.InputNew{
		Name:    "yodoc",
		Version: version,
		Out:     os.Stderr,
		Level:   logLevelVar,
	})
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	runner := cli.Runner{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		LDFlags: &cli.LDFlags{
			Version: version,
			Commit:  commit,
			Date:    date,
		},
		Logger:      logger,
		LogLevelVar: logLevelVar,
	}
	if err := runner.Run(ctx, os.Args...); err != nil {
		slogerr.WithError(logger, err).Error(errMsg(err))
		return 1
	}
	return 0
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
