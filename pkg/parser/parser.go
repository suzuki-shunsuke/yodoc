package parser

import (
	"bufio"
	"fmt"
	"io"
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
	Kind     Kind
	Lines    []string
	ReadFile string
	Dir      string
	Result   Result
	Checks   []string
}

func (p *Parser) Parse(r io.Reader) ([]*Block, error) {
	scanner := bufio.NewScanner(r)
	state := &State{
		Current: &Block{
			Kind: OutBlock,
		},
	}

	for scanner.Scan() {
		p.parse(scanner.Text(), state)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan a template file: %w", err)
	}
	return append(state.Blocks, state.Current), nil
}

func (p *Parser) parse(line string, state *State) { //nolint:cyclop
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
			state.Current = &Block{
				Lines: []string{line},
			}
			return
		}
		if state.Current.Kind == UnknownBlock {
			state.Current.Kind = OtherBlock
		}
		// end block
		state.Current.Lines = append(state.Current.Lines, line)
		state.Current = &Block{
			Kind: OutBlock,
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
