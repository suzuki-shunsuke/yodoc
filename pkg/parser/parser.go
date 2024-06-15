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

type Block struct {
	Kind     Kind
	Lines    []string
	ReadFile string
	Dir      string
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

func (p *Parser) parse(line string, state *State) {
	// #-yodoc hidden
	// ```
	// #-yodoc run
	// #-yodoc #
	// #-yodoc checks
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
		return
	case strings.HasPrefix(line, "#-yodoc checks"):
		state.Current.Kind = CheckBlock
		return
	default:
		state.Current.Lines = append(state.Current.Lines, line)
		return
	}
}
