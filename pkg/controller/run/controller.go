package run

import (
	"context"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
)

type Controller struct {
	fs           afero.Fs
	configReader ConfigReader
	renderer     Renderer
}

type ConfigReader interface {
	Read(p string, cfg *config.Config) error
}

type Renderer interface {
	Render(ctx context.Context, src, dest string) error
	SetDelims(left, right string)
	SetTasks(tasks map[string]*config.Task)
}

func NewController(fs afero.Fs, configReader ConfigReader, renderer Renderer) *Controller {
	return &Controller{
		fs:           fs,
		configReader: configReader,
		renderer:     renderer,
	}
}
