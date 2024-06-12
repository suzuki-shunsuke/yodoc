package render

import (
	"context"

	"github.com/suzuki-shunsuke/yodoc/pkg/config"
)

func Funcs(ctx context.Context, tasks map[string]*config.Task) map[string]any {
	return map[string]any{
		// Remove Command for security
		// "Command": NewCommand(ctx, nil, dir, envs).Run,
		"Task": NewTask(ctx, tasks).Run,
	}
}
