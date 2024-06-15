package run

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
	"github.com/suzuki-shunsuke/yodoc/pkg/frontmatter"
	"github.com/suzuki-shunsuke/yodoc/pkg/osfile"
	"github.com/suzuki-shunsuke/yodoc/pkg/parser"
	"github.com/suzuki-shunsuke/yodoc/pkg/render"
	"gopkg.in/yaml.v3"
)

type Param struct {
	ConfigFilePath string
}

func (c *Controller) Run(ctx context.Context, _ *logrus.Entry, param *Param) error {
	// read config
	cfg := &config.Config{}
	cfgPath, err := c.readConfig(param.ConfigFilePath, cfg)
	if err != nil {
		return err
	}

	src := filepath.Join(filepath.Dir(cfgPath), cfg.Src)
	dest := filepath.Join(filepath.Dir(cfgPath), cfg.Dest)

	renderer := render.NewRenderer(c.fs)
	if err := c.setRenderer(renderer, cfg, cfgPath); err != nil {
		return err
	}

	// create a destination directory
	if err := osfile.MkdirAll(c.fs, dest); err != nil {
		return fmt.Errorf("create a destination directory: %w", err)
	}

	// find templates
	files, err := c.findTemplates(src, cfg.Src == cfg.Dest)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := c.handleTemplate(ctx, renderer, src, dest, file, cfgPath); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) readConfig(cfgPath string, cfg *config.Config) (string, error) {
	if cfgPath == "" {
		a, err := c.configFinder.Find()
		if err != nil {
			return "", err //nolint:wrapcheck
		}
		cfgPath = a
	}
	if err := c.configReader.Read(cfgPath, cfg); err != nil {
		return "", fmt.Errorf("read a configuration file: %w", err)
	}
	return cfgPath, nil
}

func (c *Controller) setTasks(tasks map[string]*config.Task, cfg *config.Config, cfgPath string) error {
	for _, task := range cfg.Tasks {
		for _, check := range task.Checks {
			if err := check.Build(); err != nil {
				return fmt.Errorf("build a check: %w", err)
			}
		}

		for _, action := range []*config.Action{task.Action, task.Before, task.After} {
			env, err := c.renderer.GetActionEnv(action)
			if err != nil {
				return fmt.Errorf("evaluate an environment variable: %w", err)
			}
			action.SetEnv(env)
		}

		task.SetDir(filepath.Dir(cfgPath))
		if err := task.ReadScript(c.fs); err != nil {
			return err //nolint:wrapcheck
		}
		tasks[task.Name] = task
	}
	return nil
}

func (c *Controller) findTemplates(src string, isSameDir bool) ([]string, error) {
	files := []string{}
	ignoreDirs := map[string]struct{}{
		"node_modules": {},
		".git":         {},
	}
	if err := afero.Walk(c.fs, src, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if _, ok := ignoreDirs[info.Name()]; ok {
			return filepath.SkipDir
		}
		if isSameDir {
			if strings.HasSuffix(p, "_yodoc.md") || strings.HasSuffix(p, "_yodoc.mdx") {
				files = append(files, p)
			}
		} else if strings.HasSuffix(p, ".md") || strings.HasSuffix(p, ".mdx") {
			files = append(files, p)
		}
		return nil
	}); err != nil {
		return files, fmt.Errorf("walk the source directory: %w", err)
	}
	return files, nil
}

func (c *Controller) getDest(src, dest, file string, fm *frontmatter.Frontmatter) (string, error) {
	if fm != nil && fm.Dest != "" {
		return filepath.Join(filepath.Dir(file), fm.Dest), nil
	}
	if src == dest {
		if s := strings.TrimSuffix(file, "_yodoc.md"); s != file {
			return s + ".md", nil
		}
		if s := strings.TrimSuffix(file, "_yodoc.mdx"); s != file {
			return s + ".mdx", nil
		}
		return "", errors.New("the file name must end with _yodoc.md or _yodoc.mdx")
	}
	rel, err := filepath.Rel(src, file)
	if err != nil {
		return "", fmt.Errorf("get a relative path: %w", err)
	}
	dest = filepath.Join(dest, rel)
	if a, err := filepath.Rel(src, dest); err != nil {
		return "", fmt.Errorf("get a relative path: %w", err)
	} else if !strings.HasPrefix(a, "..") {
		return "", errors.New("dest must not include in source directory")
	}
	return dest, nil
}

func (c *Controller) handleTemplate(ctx context.Context, renderer Renderer, src, dest, file, cfgPath string) error {
	b, err := afero.ReadFile(c.fs, file)
	if err != nil {
		return fmt.Errorf("read a template file: %w", err)
	}
	s := string(b)
	fm, txt, err := frontmatter.Parse(s)
	if err != nil {
		return fmt.Errorf("parse frontmatter: %w", err)
	}
	dest, err = c.getDest(src, dest, file, fm)
	if err != nil {
		return err
	}

	if fm != nil && fm.Dir != "" {
		tpl, err := renderer.NewTemplate().Parse(fm.Dir)
		if err != nil {
			return fmt.Errorf("parse front matter's dir: %w", err)
		}
		b := &bytes.Buffer{}
		if err := tpl.Execute(b, map[string]any{
			"SourceDir": filepath.Dir(file),
			"DestDir":   filepath.Dir(dest),
			"ConfigDir": filepath.Dir(cfgPath),
		}); err != nil {
			return fmt.Errorf("render front matter's dir: %w", err)
		}
		fm.Dir = b.String()
	}

	r := strings.NewReader(txt)
	p := &parser.Parser{}
	blocks, err := p.Parse(r)
	if err != nil {
		return err
	}
	texts := make([]string, 0, len(blocks))
	result := &render.CommandResult{}
	for _, block := range blocks {
		switch block.Kind {
		case parser.HiddenBlock:
			txt := strings.Join(block.Lines[1:len(block.Lines)-1], "\n")
			cmd := exec.CommandContext(ctx, "sh", "-c", txt)
			combinedOutput := &bytes.Buffer{}
			cmd.Stdout = combinedOutput
			cmd.Stderr = combinedOutput
			cmd.Dir = filepath.Join(fm.Dir, block.Dir)
			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stderr, "[ERROR] Hidden command failed", "\n", txt, "\n", combinedOutput.String())
				return err
			}
		case parser.MainBlock:
			cmd := render.NewCommand(ctx, []string{"sh", "-c"}, filepath.Join(fm.Dir, block.Dir), nil)
			s := strings.Join(block.Lines[1:len(block.Lines)-1], "\n")
			result = cmd.Run(s)
			txt := strings.Join(block.Lines, "\n")
			txt, err := c.render(ctx, renderer, file, fm, txt, result)
			if err != nil {
				return err
			}
			texts = append(texts, txt)
		case parser.CheckBlock:
			checks := struct {
				Checks []*config.Check
			}{}
			txt := strings.Join(block.Lines[1:len(block.Lines)-1], "\n")
			if err := yaml.Unmarshal([]byte(txt), &checks); err != nil {
				return err
			}
			for _, check := range checks.Checks {
				if err := c.check(check, result); err != nil {
					fmt.Fprintln(os.Stderr, "[ERROR] Check failed", "\n", result.Command, "\n", check.Expr, "\n", err.Error())
					return err
				}
			}
		case parser.OtherBlock:
			txt, err := c.render(ctx, renderer, file, fm, strings.Join(block.Lines, "\n"), result)
			if err != nil {
				return err
			}
			texts = append(texts, txt)
		case parser.OutBlock:
			txt, err := c.render(ctx, renderer, file, fm, strings.Join(block.Lines, "\n"), result)
			if err != nil {
				return err
			}
			texts = append(texts, txt)
		default:
			return fmt.Errorf("unknown block kind: %v", block.Kind)
		}
	}

	if err := afero.WriteFile(c.fs, dest, []byte(strings.Join(texts, "\n")+render.Footer), 0o644); err != nil { //nolint:mnd
		return err
	}

	return nil
}

func (c *Controller) check(check *config.Check, result *render.CommandResult) error {
	if err := check.Build(); err != nil {
		return err
	}
	prog := check.GetExpr()
	output, err := expr.Run(prog, result)
	if err != nil {
		return fmt.Errorf("evaluate an expression: %w", err)
	}
	b, ok := output.(bool)
	if !ok {
		return errors.New("the result of the expression isn't a boolean value")
	}
	if !b {
		return errors.New("a check is false")
	}
	return nil
}

func (c *Controller) render(ctx context.Context, renderer Renderer, file string, fm *frontmatter.Frontmatter, txt string, result *render.CommandResult) (string, error) {
	tpl := renderer.NewTemplate().Funcs(render.Funcs(ctx, c.fs, file, fm))
	tpl.Delims(fm.Delim.GetLeft(), fm.Delim.GetRight())
	tpl, err := tpl.Parse(txt)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, result); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *Controller) setRenderer(renderer Renderer, cfg *config.Config, cfgPath string) error {
	tasks := make(map[string]*config.Task, len(cfg.Tasks))
	if err := c.setTasks(tasks, cfg, cfgPath); err != nil {
		return err
	}

	renderer.SetDelims(cfg.Delim.GetLeft(), cfg.Delim.GetRight())
	renderer.SetTasks(tasks)
	return nil
}
