package cli

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
	"github.com/suzuki-shunsuke/yodoc/pkg/controller/run"
	"github.com/suzuki-shunsuke/yodoc/pkg/render"
	"github.com/urfave/cli/v3"
)

type runCommand struct {
	logger      *slog.Logger
	logLevelVar *slog.LevelVar
}

func (rc *runCommand) command() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "Generate documents",
		UsageText: "yodoc run",
		Description: `Generate documents.

$ yodoc run
`,
		Action: rc.action,
	}
}

func (rc *runCommand) action(ctx context.Context, cmd *cli.Command) error {
	fs := afero.NewOsFs()
	configReader := config.NewReader(fs)
	renderer := render.NewRenderer(fs)
	finder := config.NewFinder(fs)
	ctrl := run.NewController(fs, finder, configReader, renderer)
	if err := slogutil.SetLevel(rc.logLevelVar, cmd.String("log-level")); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	return ctrl.Run(ctx, rc.logger, &run.Param{ //nolint:wrapcheck
		ConfigFilePath: cmd.String("config"),
		Files:          cmd.Args().Slice(),
	})
}
