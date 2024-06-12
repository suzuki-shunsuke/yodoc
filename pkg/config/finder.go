package config

import (
	"errors"

	"github.com/spf13/afero"
)

type Finder struct {
	fs afero.Fs
}

func NewFinder(fs afero.Fs) *Finder {
	return &Finder{
		fs: fs,
	}
}

func (f *Finder) Find() (string, error) {
	for _, filePath := range []string{"yodoc.yaml", ".yodoc.yaml"} {
		if _, err := f.fs.Stat(filePath); err == nil {
			return filePath, nil
		}
	}
	return "", errors.New("configuration file not found")
}
