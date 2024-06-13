package render

import (
	"context"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
	"github.com/suzuki-shunsuke/yodoc/pkg/frontmatter"
)

func Funcs(ctx context.Context, fs afero.Fs, tasks map[string]*config.Task, src string, fm *frontmatter.Frontmatter) map[string]any {
	return map[string]any{
		// Remove Command for security
		// "Command": NewCommand(ctx, nil, dir, envs).Run,
		"Task": NewTask(ctx, tasks, fm).Run,
		"Read": NewRead(ctx, fs, filepath.Dir(src)).Run,
	}
}
