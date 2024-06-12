package run

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
	"github.com/suzuki-shunsuke/yodoc/pkg/osfile"
	"github.com/suzuki-shunsuke/yodoc/pkg/render"
)

type Param struct {
	ConfigFilePath string
}

func (c *Controller) Run(ctx context.Context, _ *logrus.Entry, param *Param) error {
	// read config
	cfg := &config.Config{}
	cfgPath := param.ConfigFilePath
	if cfgPath == "" {
		a, err := c.configFinder.Find()
		if err != nil {
			return err //nolint:wrapcheck
		}
		cfgPath = a
	}
	if err := c.configReader.Read(cfgPath, cfg); err != nil {
		return fmt.Errorf("read a configuration file: %w", err)
	}

	src := filepath.Join(filepath.Dir(cfgPath), cfg.Src)
	dest := filepath.Join(filepath.Dir(cfgPath), cfg.Dest)

	tasks := make(map[string]*config.Task, len(cfg.Tasks))
	for _, task := range cfg.Tasks {
		for _, check := range task.Checks {
			if err := check.Build(); err != nil {
				return fmt.Errorf("build a check: %w", err)
			}
		}
		task.SetEnv()
		task.SetDir(filepath.Dir(cfgPath))
		if err := task.ReadScript(c.fs); err != nil {
			return err //nolint:wrapcheck
		}
		tasks[task.Name] = task
	}

	renderer := render.NewRenderer(c.fs)

	renderer.SetDelims(cfg.Delim.GetLeft(), cfg.Delim.GetRight())
	renderer.SetTasks(tasks)
	// create a destination directory
	if err := osfile.MkdirAll(c.fs, dest); err != nil {
		return fmt.Errorf("create a destination directory: %w", err)
	}
	// find and read template
	files := []string{}
	if err := afero.Walk(c.fs, src, func(p string, _ os.FileInfo, err error) error {
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
		rel, err := filepath.Rel(src, file)
		if err != nil {
			return fmt.Errorf("get a relative path: %w", err)
		}
		dest := filepath.Join(dest, rel)
		if a, err := filepath.Rel(src, dest); err != nil {
			return fmt.Errorf("get a relative path: %w", err)
		} else if !strings.HasPrefix(a, "..") {
			return errors.New("dest must not include in source directory")
		}
		// render templates and update documents
		if err := renderer.Render(ctx, file, dest); err != nil {
			return fmt.Errorf("render a file: %w", err)
		}
	}
	return nil
}
