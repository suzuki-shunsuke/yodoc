package parser

import (
	"bufio"
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
	HiddenBlock Kind = iota
	MainBlock
	CheckBlock
	OtherBlock
	OutBlock
)

type Block struct {
	Kind  Kind
	Lines []string
}

func (p *Parser) Parse(r io.Reader) ([]*Block, error) {
	scanner := bufio.NewScanner(r)
	state := &State{
		Current: &Block{
			Kind: OutBlock,
		},
	}

	for scanner.Scan() {
		if err := p.parse(scanner.Text(), state); err != nil {
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return append(state.Blocks, state.Current), nil
}

func (p *Parser) parse(line string, state *State) error {
	// #-yodoc hidden
	// ```
	// #-yodoc run
	// #-yodoc #
	// #-yodoc checks
	if strings.HasPrefix(line, "#-yodoc # ") {
		return nil
	}
	if strings.HasPrefix(line, "```") {
		state.Blocks = append(state.Blocks, state.Current)
		if state.Current.Kind == OutBlock {
			// start block
			// TODO
			state.Current = &Block{
				Lines: []string{line},
			}
			return nil
		}
		// end block
		state.Current.Lines = append(state.Current.Lines, line)
		state.Current = &Block{
			Kind: OutBlock,
		}
		return nil
	}
	switch {
	case strings.HasPrefix(line, "#-yodoc hidden"):
		state.Current.Kind = HiddenBlock
		return nil
	case strings.HasPrefix(line, "#-yodoc run"):
		state.Current.Kind = MainBlock
		return nil
	case strings.HasPrefix(line, "#-yodoc checks"):
		state.Current.Kind = CheckBlock
		return nil
	default:
		state.Current.Lines = append(state.Current.Lines, line)
		return nil
	}
}
