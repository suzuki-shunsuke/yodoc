package cli

import (
	"context"
	"fmt"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
	"github.com/suzuki-shunsuke/yodoc/pkg/controller/run"
	"github.com/suzuki-shunsuke/yodoc/pkg/render"
	"github.com/urfave/cli/v3"
)

type runCommand struct{}

func (rc *runCommand) command(logger *slogutil.Logger) *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "Generate documents",
		UsageText: "yodoc run",
		Description: `Generate documents.

$ yodoc run
`,
		Action: urfave.Action(rc.action, logger),
	}
}

func (rc *runCommand) action(ctx context.Context, cmd *cli.Command, logger *slogutil.Logger) error {
	fs := afero.NewOsFs()
	configReader := config.NewReader(fs)
	renderer := render.NewRenderer(fs)
	finder := config.NewFinder(fs)
	ctrl := run.NewController(fs, finder, configReader, renderer)
	if err := logger.SetLevel(cmd.String("log-level")); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	if err := logger.SetColor(cmd.String("log-color")); err != nil {
		return fmt.Errorf("set log color: %w", err)
	}
	return ctrl.Run(ctx, logger.Logger, &run.Param{ //nolint:wrapcheck
		ConfigFilePath: cmd.String("config"),
		Files:          cmd.Args().Slice(),
	})
}
