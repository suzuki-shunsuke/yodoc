package parser

import (
	"bufio"
	"fmt"
	"strings"
)

type Parser struct{}

type State struct {
	Blocks  []*Block
	Current *Block
}

type Kind int

const (
	UnknownBlock Kind = iota
	HiddenBlock
	MainBlock
	CheckBlock
	OtherBlock
	OutBlock
)

type Result int

const (
	UnknownResult Result = iota
	SuccessResult
	FailureResult
)

type Block struct {
	Kind      Kind
	Lines     []string
	ReadFile  string
	Dir       string
	Result    Result
	Checks    []string
	StartLine int
	EndLine   int
}

func (p *Parser) Parse(ln int, lastLine string, scanner *bufio.Scanner) ([]*Block, error) {
	state := &State{
		Current: &Block{
			Kind:      OutBlock,
			StartLine: 1,
			Lines:     []string{lastLine},
		},
	}

	for scanner.Scan() {
		p.parse(ln, scanner.Text(), state)
		ln++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan a template file: %w", err)
	}
	state.Current.EndLine = ln - 1
	return append(state.Blocks, state.Current), nil
}

func (p *Parser) parse(ln int, line string, state *State) { //nolint:cyclop,funlen
	// #-yodoc hidden
	// ```
	// #-yodoc run
	// #-yodoc #
	// #-yodoc check
	// #-yodoc check [<expr>]
	// #-yodoc dir
	if strings.HasPrefix(line, "#-yodoc # ") {
		return
	}
	if strings.HasPrefix(line, "#-yodoc dir ") {
		state.Current.Dir = strings.TrimPrefix(line, "#-yodoc dir ")
		return
	}
	if strings.HasPrefix(line, "```") {
		state.Blocks = append(state.Blocks, state.Current)
		if state.Current.Kind == OutBlock {
			// start block
			state.Current.EndLine = ln - 1
			state.Current = &Block{
				Lines:     []string{line},
				StartLine: ln,
			}
			return
		}
		if state.Current.Kind == UnknownBlock {
			state.Current.Kind = OtherBlock
		}
		// end block
		state.Current.Lines = append(state.Current.Lines, line)
		state.Current.EndLine = ln
		state.Current = &Block{
			Kind:      OutBlock,
			StartLine: ln + 1,
		}
		return
	}
	switch {
	case strings.HasPrefix(line, "#-yodoc hidden"):
		state.Current.Kind = HiddenBlock
		return
	case strings.HasPrefix(line, "#-yodoc run"):
		state.Current.Kind = MainBlock
		state.Current.Result = SuccessResult
		return
	case strings.HasPrefix(line, "#!yodoc run"):
		state.Current.Kind = MainBlock
		state.Current.Result = FailureResult
		return
	case strings.HasPrefix(line, "#-yodoc check"):
		if state.Current.Kind == MainBlock {
			e := strings.TrimPrefix(line, "#-yodoc check")
			if e != "" {
				state.Current.Checks = append(state.Current.Checks, e)
			}
			return
		}
		state.Current.Kind = CheckBlock
		return
	default:
		state.Current.Lines = append(state.Current.Lines, line)
		return
	}
}
