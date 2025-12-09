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

	Files []string
}

func (rc *runCommand) command(logger *slogutil.Logger, flags *Flags) *cli.Command {
	runFlags := &RunFlags{
		Flags: flags,
	}
	return &cli.Command{
		Name:      "run",
		Usage:     "Generate documents",
		UsageText: "yodoc run",
		Description: `Generate documents.

$ yodoc run
`,
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:        "file",
				Max:         -1,
				Destination: &runFlags.Files,
			},
		},
		Action: func(ctx context.Context, _ *cli.Command) error {
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
		Files:          flags.Files,
	})
}
