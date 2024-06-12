package render

import (
	"context"

	"github.com/suzuki-shunsuke/yodoc/pkg/config"
)

func Funcs(ctx context.Context, tasks map[string]*config.Task, dir string, envs []string) map[string]any {
	return map[string]any{
		"Command": NewCommand(ctx, nil, dir, envs).Run,
		"Task":    NewTask(ctx, tasks, dir, envs).Run,
	}
}
