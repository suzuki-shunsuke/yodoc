package render

import (
	"context"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
)

func Funcs(ctx context.Context, fs afero.Fs, tasks map[string]*config.Task, src string) map[string]any {
	return map[string]any{
		// Remove Command for security
		// "Command": NewCommand(ctx, nil, dir, envs).Run,
		"Task": NewTask(ctx, tasks).Run,
		"Read": NewRead(ctx, fs, filepath.Dir(src)).Run,
	}
}
