package render

import (
	"path/filepath"

	"github.com/spf13/afero"
)

func Funcs(fs afero.Fs, src string) map[string]any {
	return map[string]any{
		"Read": NewRead(fs, filepath.Dir(src)).Run,
	}
}
