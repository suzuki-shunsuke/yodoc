package run

import (
	"text/template"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
	"github.com/suzuki-shunsuke/yodoc/pkg/frontmatter"
)

type Controller struct {
	fs           afero.Fs
	configReader ConfigReader
	configFinder ConfigFinder
	renderer     Renderer
}

type ConfigFinder interface {
	Find() (string, error)
}

type ConfigReader interface {
	Read(p string, cfg *config.Config) error
}

type Renderer interface {
	RenderFile(src, dest, txt string, delim *frontmatter.Delim) error
	NewTemplate() *template.Template
	SetDelims(left, right string)
}

func NewController(fs afero.Fs, configFinder ConfigFinder, configReader ConfigReader, renderer Renderer) *Controller {
	return &Controller{
		fs:           fs,
		configReader: configReader,
		configFinder: configFinder,
		renderer:     renderer,
	}
}
