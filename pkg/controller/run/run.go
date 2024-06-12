package run

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
	"github.com/suzuki-shunsuke/yodoc/pkg/osfile"
)

func (c *Controller) Run(ctx context.Context, _ *logrus.Entry) error {
	// read config
	cfg := &config.Config{}
	if err := c.configReader.Read("yodoc.yaml", cfg); err != nil {
		return fmt.Errorf("read a configuration file: %w", err)
	}
	// create a destination directory
	if err := osfile.MkdirAll(c.fs, cfg.Dest); err != nil {
		return fmt.Errorf("create a destination directory: %w", err)
	}
	// find and read template
	files := []string{}
	if err := afero.Walk(c.fs, cfg.Src, func(p string, _ os.FileInfo, err error) error {
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
