package run

import (
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
	Render(src, dest string) error
}

func NewController(fs afero.Fs, configReader ConfigReader, renderer Renderer) *Controller {
	return &Controller{
		fs:           fs,
		configReader: configReader,
		renderer:     renderer,
	}
}
