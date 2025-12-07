package cli

import (
	"context"
	"fmt"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
	"github.com/suzuki-shunsuke/yodoc/pkg/controller/run"
	"github.com/suzuki-shunsuke/yodoc/pkg/render"
	"github.com/urfave/cli/v3"
)

type runCommand struct{}

type RunFlags struct {
	*Flags

	Args []string
}

func (rc *runCommand) command(logger *slogutil.Logger, flags *Flags) *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "Generate documents",
		UsageText: "yodoc run",
		Description: `Generate documents.

$ yodoc run
`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			runFlags := &RunFlags{
				Flags: flags,
				Args:  cmd.Args().Slice(),
			}
			return rc.action(ctx, logger, runFlags)
		},
	}
}

func (rc *runCommand) action(ctx context.Context, logger *slogutil.Logger, flags *RunFlags) error {
	fs := afero.NewOsFs()
	configReader := config.NewReader(fs)
	renderer := render.NewRenderer(fs)
	finder := config.NewFinder(fs)
	ctrl := run.NewController(fs, finder, configReader, renderer)
	if err := logger.SetLevel(flags.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	if err := logger.SetColor(flags.LogColor); err != nil {
		return fmt.Errorf("set log color: %w", err)
	}
	return ctrl.Run(ctx, logger.Logger, &run.Param{ //nolint:wrapcheck
		ConfigFilePath: flags.Config,
		Files:          flags.Args,
	})
}
