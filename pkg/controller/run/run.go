package run

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
)

func (c *Controller) Run(ctx context.Context, _ *logrus.Entry) error {
	// read config
	cfg := &config.Config{}
	if err := c.configReader.Read("yodoc.yaml", cfg); err != nil {
		return fmt.Errorf("read a configuration file: %w", err)
	}
	// find and read template
	files := []string{}
	if err := afero.Walk(c.fs, cfg.Src, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(p, ".md") || strings.HasSuffix(p, ".mdx") {
			files = append(files, p)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("walk the source directory: %w", err)
	}
	for _, file := range files {
		rel, err := filepath.Rel(cfg.Src, file)
		if err != nil {
			return fmt.Errorf("get a relative path: %w", err)
		}
		dest := filepath.Join(cfg.Dest, rel)
		// render templates and update documents
		if err := c.renderer.Render(file, dest); err != nil {
			return fmt.Errorf("render a file: %w", err)
		}
	}
	return nil
}

type commander struct {
	ctx   context.Context
	tasks map[string]*Task
}

type Task struct {
	Name    string
	Command string
}

func (c *commander) command(s string) (string, error) {
	cmd := exec.CommandContext(c.ctx, "sh", "-c", s)
	combinedOutput := &bytes.Buffer{}

	cmd.Stdout = combinedOutput
	cmd.Stderr = combinedOutput

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return combinedOutput.String(), nil
}

func (c *commander) task(taskID string) (any, error) {
	task, ok := c.tasks[taskID]
	if !ok {
		return nil, errors.New("task not found")
	}
	s, err := c.command(task.Command)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"Command":        task.Command,
		"CombinedOutput": s,
	}, nil
}

func core() error {
	tasks := []*Task{
		{
			Name:    "hello",
			Command: "echo Hello",
		},
	}
	ctx := context.Background()
	cmdr := &commander{
		ctx: ctx,
	}
	cmdr.tasks = make(map[string]*Task, len(tasks))
	for _, task := range tasks {
		cmdr.tasks[task.Name] = task
	}
	const text = `
{{with Task "hello"}}	

$ {{.Command}}

{{.CombinedOutput}}

{{end}}
`
	tpl, err := template.New("_").Funcs(map[string]any{
		"Command": cmdr.command,
		"Task":    cmdr.task,
	}).Parse(text)
	if err != nil {
		return err
	}
	if err := tpl.Execute(os.Stdout, nil); err != nil {
		return err
	}
	return nil
}
