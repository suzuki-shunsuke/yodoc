package expr

import (
	"errors"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type Parser struct{}

func (p *Parser) Eval(input string, result any) (bool, error) {
	prog, err := p.compile(input)
	if err != nil {
		return false, err
	}
	return p.eval(prog, result)
}

func (p *Parser) compile(input string) (*vm.Program, error) {
	prog, err := expr.Compile(input, expr.AsBool())
	if err != nil {
		return nil, fmt.Errorf("compile an expression: %w", err)
	}
	return prog, nil
}

func (p *Parser) eval(prog *vm.Program, result any) (bool, error) {
	output, err := expr.Run(prog, result)
	if err != nil {
		return false, fmt.Errorf("evaluate an expression: %w", err)
	}
	b, ok := output.(bool)
	if !ok {
		return false, errors.New("the result of the expression isn't a boolean value")
	}
	return b, nil
}
