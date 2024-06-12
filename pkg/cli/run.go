package cli

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
	"github.com/suzuki-shunsuke/yodoc/pkg/controller/run"
	"github.com/suzuki-shunsuke/yodoc/pkg/log"
	"github.com/suzuki-shunsuke/yodoc/pkg/render"
	"github.com/urfave/cli/v2"
)

type runCommand struct {
	logE *logrus.Entry
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

func (rc *runCommand) action(c *cli.Context) error {
	fs := afero.NewOsFs()
	configReader := config.NewReader(fs)
	renderer := render.NewRenderer(fs)
	finder := config.NewFinder(fs)
	ctrl := run.NewController(fs, finder, configReader, renderer)
	logE := rc.logE
	log.SetLevel(c.String("log-level"), logE)
	log.SetColor(c.String("log-color"), logE)
	return ctrl.Run(c.Context, logE, &run.Param{
		ConfigFilePath: c.String("config"),
	}) //nolint:wrapcheck
}
