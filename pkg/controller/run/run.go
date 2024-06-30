package run

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/suzuki-shunsuke/yodoc/pkg/config"
	"github.com/suzuki-shunsuke/yodoc/pkg/expr"
	"github.com/suzuki-shunsuke/yodoc/pkg/frontmatter"
	"github.com/suzuki-shunsuke/yodoc/pkg/osfile"
	"github.com/suzuki-shunsuke/yodoc/pkg/parser"
	"github.com/suzuki-shunsuke/yodoc/pkg/render"
	"gopkg.in/yaml.v3"
)

type Param struct {
	ConfigFilePath string
	Files          []string
}

func (c *Controller) Run(ctx context.Context, _ *logrus.Entry, param *Param) error {
	// find and read config
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

	files := param.Files
	if len(files) == 0 {
		// find templates
		a, err := c.findTemplates(src, cfg.Src == cfg.Dest)
		if err != nil {
			return err
		}
		files = a
	}

	for _, file := range files {
		if err := c.handleTemplate(ctx, renderer, src, dest, file, cfgPath); err != nil {
			return logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"file": file,
			})
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
	// read and parse a template file
	f, err := c.fs.Open(file)
	if err != nil {
		return fmt.Errorf("open a template file: %w", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	// parse frontmatter
	fm, lastLine, ln, err := frontmatter.Parse(scanner)
	if err != nil {
		return fmt.Errorf("parse frontmatter: %w", err)
	}

	// get a destination file path
	dest, err = c.getDest(src, dest, file, fm)
	if err != nil {
		return err
	}

	if fm != nil && fm.Dir != "" {
		if err := c.setFrontmatterDir(renderer, fm, file, dest, cfgPath); err != nil {
			return err
		}
	}

	// parse a template file
	p := &parser.Parser{}
	blocks, err := p.Parse(ln+1, lastLine, scanner)
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

	// write a document
	if err := afero.WriteFile(c.fs, dest, []byte(strings.Join(texts, "\n")+render.Footer), 0o644); err != nil { //nolint:mnd
		return fmt.Errorf("write a document: %w", err)
	}

	return nil
}

type CommandError struct {
	Command        string
	CombinedOutput string
	Expr           string
	Checks         string
	err            error
	Start          int
	End            int
}

func NewCommandError(err error, command, combinedOutput string) *CommandError {
	return &CommandError{
		Command:        command,
		CombinedOutput: combinedOutput,
		err:            err,
	}
}

func (e *CommandError) WithExpr(exp string) *CommandError {
	e.Expr = exp
	return e
}

func (e *CommandError) WithChecks(checks string) *CommandError {
	e.Checks = checks
	return e
}

func (e *CommandError) WithLocation(start, end int) *CommandError {
	e.Start = start
	e.End = end
	return e
}

func (e *CommandError) Error() string {
	return e.err.Error()
}

func (e *CommandError) Unwrap() error {
	return e.err
}

// handleCommandBlock handles a command block and returns a command result and a text.
func (c *Controller) handleCommandBlock(ctx context.Context, renderer Renderer, fm *frontmatter.Frontmatter, file string, block *parser.Block, hidden bool) (*render.CommandResult, string, error) {
	cmd := render.NewCommand(ctx, []string{"sh", "-c"}, filepath.Join(fm.Dir, block.Dir), nil)
	s := strings.Join(block.Lines[1:len(block.Lines)-1], "\n")
	// run a command
	result := cmd.Run(s)

	// check the command exit code
	if block.Result == parser.SuccessResult && result.ExitCode != 0 {
		return result, "", NewCommandError(errors.New("command failed"), s, result.CombinedOutput).WithLocation(block.StartLine, block.EndLine)
	}
	if block.Result == parser.FailureResult && result.ExitCode == 0 {
		return result, "", NewCommandError(errors.New("command must fail but succeeded"), s, result.CombinedOutput)
	}

	// check the command result based on checks
	for _, e := range block.Checks {
		check := &config.Check{
			Expr: e,
		}
		if err := c.check(check, result); err != nil {
			return result, "", NewCommandError(err, result.Command, result.CombinedOutput).WithExpr(check.Expr)
		}
	}

	if hidden {
		// If the block is hidden, the command result isn't rendered.
		return result, "", nil
	}

	// render the block
	txt := strings.Join(block.Lines, "\n")
	txt, err := c.render(renderer, file, fm.GetDelim(), txt, result)
	if err != nil {
		return result, "", NewCommandError(err, s, result.CombinedOutput)
	}
	return result, txt, nil
}

// handleCheckBlock handles a check block.
func (c *Controller) handleCheckBlock(block *parser.Block, result *render.CommandResult) error {
	// Remove code block markers ``` and get the content of the block.
	txt := strings.Join(block.Lines[1:len(block.Lines)-1], "\n")
	checks := struct {
		Checks []*config.Check
	}{}
	// Unmarshal the content of the block as YAML.
	if err := yaml.Unmarshal([]byte(txt), &checks); err != nil {
		return fmt.Errorf("unmarshal a checks block as YAML: %w", NewCommandError(err, result.Command, result.CombinedOutput).WithChecks(txt))
	}
	// Checks the command result based on checks.
	for _, check := range checks.Checks {
		if err := c.check(check, result); err != nil {
			return NewCommandError(err, result.Command, result.CombinedOutput).WithExpr(check.Expr)
		}
	}
	return nil
}

// handleOtherBlock handles a other block.
func (c *Controller) handleOtherBlock(renderer Renderer, delim *frontmatter.Delim, file string, block *parser.Block, result *render.CommandResult) (string, error) {
	return c.render(renderer, file, delim, strings.Join(block.Lines, "\n"), result)
}

// handleOutBlock handles a out block.
func (c *Controller) handleOutBlock(renderer Renderer, delim *frontmatter.Delim, file string, block *parser.Block, result *render.CommandResult) (string, error) {
	return c.render(renderer, file, delim, strings.Join(block.Lines, "\n"), result)
}

// handleBlock handles a block and converts it to a text and gets a command result.
func (c *Controller) handleBlock(ctx context.Context, renderer Renderer, fm *frontmatter.Frontmatter, file string, block *parser.Block, result *render.CommandResult) (*render.CommandResult, string, error) {
	switch block.Kind {
	case parser.HiddenBlock:
		return c.handleCommandBlock(ctx, renderer, fm, file, block, true)
	case parser.MainBlock:
		return c.handleCommandBlock(ctx, renderer, fm, file, block, false)
	case parser.CheckBlock:
		return result, "", c.handleCheckBlock(block, result)
	case parser.OtherBlock:
		txt, err := c.handleOtherBlock(renderer, fm.GetDelim(), file, block, result)
		return result, txt, err
	case parser.OutBlock:
		txt, err := c.handleOutBlock(renderer, fm.GetDelim(), file, block, result)
		return result, txt, err
	default:
		return result, "", fmt.Errorf("unknown block kind: %v", block.Kind)
	}
}

// setFrontmatterDir renders the front matter's dir.
func (c *Controller) setFrontmatterDir(renderer Renderer, fm *frontmatter.Frontmatter, file, dest, cfgPath string) error {
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
	parser := &expr.Parser{}
	b, err := parser.Eval(check.Expr, result)
	if err != nil {
		return fmt.Errorf("evaluate a check: %w", err)
	}
	if !b {
		return errors.New("a check is false")
	}
	return nil
}

// render renders a template.
func (c *Controller) render(renderer Renderer, file string, delim *frontmatter.Delim, txt string, result *render.CommandResult) (string, error) {
	tpl := renderer.NewTemplate().Funcs(render.Funcs(c.fs, file))
	if delim != nil {
		tpl.Delims(delim.GetLeft(), delim.GetRight())
	}
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

// setRenderer configres the renderer based on config.
func (c *Controller) setRenderer(renderer Renderer, cfg *config.Config) error {
	renderer.SetDelims(cfg.Delim.GetLeft(), cfg.Delim.GetRight())
	return nil
}
