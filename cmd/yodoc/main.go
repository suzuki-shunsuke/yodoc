package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/suzuki-shunsuke/yodoc/pkg/cli"
	"github.com/suzuki-shunsuke/yodoc/pkg/controller/run"
	"github.com/suzuki-shunsuke/yodoc/pkg/log"
)

var (
	version = ""
	commit  = "" //nolint:gochecknoglobals
	date    = "" //nolint:gochecknoglobals
)

func main() {
	logE := log.New(version)
	if err := core(logE); err != nil {
		logerr.WithError(logE, err).Fatal(errMsg(err))
	}
}

func errMsg(err error) string {
	ce := &run.CommandError{}
	msg := "yodoc failed"
	if errors.As(err, &ce) {
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
	}
	return msg
}

func core(logE *logrus.Entry) error {
	runner := cli.Runner{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		LDFlags: &cli.LDFlags{
			Version: version,
			Commit:  commit,
			Date:    date,
		},
		LogE: logE,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return runner.Run(ctx, os.Args...) //nolint:wrapcheck
}
