package cli

import (
	"context"
	"fmt"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/yodoc/pkg/controller/initcmd"
	"github.com/urfave/cli/v3"
)

type initCommand struct{}

func (lc *initCommand) command(logger *slogutil.Logger, flags *Flags) *cli.Command {
	return &cli.Command{
		Name:      "init",
		Usage:     "Scaffold configuration file",
		UsageText: "yodoc init",
		Description: `Scaffold configuration file.

$ yodoc init

This command generates yodoc.yaml.
If the file already exists, this command does nothing.
`,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return lc.action(ctx, logger, flags)
		},
	}
}

func (lc *initCommand) action(ctx context.Context, logger *slogutil.Logger, flags *Flags) error {
	ctrl := initcmd.NewController(afero.NewOsFs())
	if err := logger.SetLevel(flags.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	if err := logger.SetColor(flags.LogColor); err != nil {
		return fmt.Errorf("set log color: %w", err)
	}
	return ctrl.Init(ctx, logger.Logger) //nolint:wrapcheck
}
