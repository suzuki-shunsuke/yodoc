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
	if err := c.setRenderer(renderer, cfg); err != nil {
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

func (c *Controller) handleTemplate(ctx context.Context, renderer Renderer, src, dest, file, cfgPath string) error { //nolint:cyclop
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
		if err := c.setFormatterDir(renderer, fm, file, dest, cfgPath); err != nil {
			return err
		}
	}

	r := strings.NewReader(txt)
	p := &parser.Parser{}
	blocks, err := p.Parse(r)
	if err != nil {
		return fmt.Errorf("parse a template: %w", err)
	}
	texts := make([]string, 0, len(blocks))
	result := &render.CommandResult{}
	for _, block := range blocks {
		a, txt, err := c.handleBlock(ctx, renderer, fm, file, block, result)
		if err != nil {
			return err
		}
		if txt != "" {
			texts = append(texts, txt)
		}
		result = a
	}

	if err := afero.WriteFile(c.fs, dest, []byte(strings.Join(texts, "\n")+render.Footer), 0o644); err != nil { //nolint:mnd
		return fmt.Errorf("write a document: %w", err)
	}

	return nil
}

func (c *Controller) handleHiddenBlock(ctx context.Context, fm *frontmatter.Frontmatter, block *parser.Block) error {
	txt := strings.Join(block.Lines[1:len(block.Lines)-1], "\n")
	cmd := exec.CommandContext(ctx, "sh", "-c", txt)
	combinedOutput := &bytes.Buffer{}
	cmd.Stdout = combinedOutput
	cmd.Stderr = combinedOutput
	cmd.Dir = filepath.Join(fm.Dir, block.Dir)
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] Hidden command failed", "\n", txt, "\n", combinedOutput.String())
		return fmt.Errorf("execute a hidden block: %w", err)
	}
	return nil
}

func (c *Controller) handleMainBlock(ctx context.Context, renderer Renderer, fm *frontmatter.Frontmatter, file string, block *parser.Block) (*render.CommandResult, string, error) {
	cmd := render.NewCommand(ctx, []string{"sh", "-c"}, filepath.Join(fm.Dir, block.Dir), nil)
	s := strings.Join(block.Lines[1:len(block.Lines)-1], "\n")
	result := cmd.Run(s)
	txt := strings.Join(block.Lines, "\n")
	txt, err := c.render(renderer, file, fm, txt, result)
	if err != nil {
		return result, "", err
	}
	return result, txt, nil
}

func (c *Controller) handleCheckBlock(block *parser.Block, result *render.CommandResult) error {
	checks := struct {
		Checks []*config.Check
	}{}
	txt := strings.Join(block.Lines[1:len(block.Lines)-1], "\n")
	if err := yaml.Unmarshal([]byte(txt), &checks); err != nil {
		return fmt.Errorf("unmarshal a checks block as YAML: %w", err)
	}
	for _, check := range checks.Checks {
		if err := c.check(check, result); err != nil {
			fmt.Fprintln(os.Stderr, "[ERROR] Check failed", "\n", result.Command, "\n", check.Expr, "\n", err.Error())
			return err
		}
	}
	return nil
}

func (c *Controller) handleOtherBlock(renderer Renderer, fm *frontmatter.Frontmatter, file string, block *parser.Block, result *render.CommandResult) (string, error) {
	return c.render(renderer, file, fm, strings.Join(block.Lines, "\n"), result)
}

func (c *Controller) handleOutBlock(renderer Renderer, fm *frontmatter.Frontmatter, file string, block *parser.Block, result *render.CommandResult) (string, error) {
	return c.render(renderer, file, fm, strings.Join(block.Lines, "\n"), result)
}

func (c *Controller) handleBlock(ctx context.Context, renderer Renderer, fm *frontmatter.Frontmatter, file string, block *parser.Block, result *render.CommandResult) (*render.CommandResult, string, error) {
	switch block.Kind {
	case parser.HiddenBlock:
		return result, "", c.handleHiddenBlock(ctx, fm, block)
	case parser.MainBlock:
		return c.handleMainBlock(ctx, renderer, fm, file, block)
	case parser.CheckBlock:
		return result, "", c.handleCheckBlock(block, result)
	case parser.OtherBlock:
		txt, err := c.handleOtherBlock(renderer, fm, file, block, result)
		return result, txt, err
	case parser.OutBlock:
		txt, err := c.handleOutBlock(renderer, fm, file, block, result)
		return result, txt, err
	default:
		return result, "", fmt.Errorf("unknown block kind: %v", block.Kind)
	}
}

func (c *Controller) setFormatterDir(renderer Renderer, fm *frontmatter.Frontmatter, file, dest, cfgPath string) error {
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
	return nil
}

func (c *Controller) check(check *config.Check, result *render.CommandResult) error {
	if err := check.Build(); err != nil {
		return fmt.Errorf("build a check: %w", err)
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

func (c *Controller) render(renderer Renderer, file string, fm *frontmatter.Frontmatter, txt string, result *render.CommandResult) (string, error) {
	tpl := renderer.NewTemplate().Funcs(render.Funcs(c.fs, file))
	tpl.Delims(fm.Delim.GetLeft(), fm.Delim.GetRight())
	tpl, err := tpl.Parse(txt)
	if err != nil {
		return "", fmt.Errorf("parse a template: %w", err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, result); err != nil {
		return "", fmt.Errorf("render a template: %w", err)
	}
	return buf.String(), nil
}

func (c *Controller) setRenderer(renderer Renderer, cfg *config.Config) error {
	renderer.SetDelims(cfg.Delim.GetLeft(), cfg.Delim.GetRight())
	return nil
}
