package config

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type Check struct {
	Expr string
	expr *vm.Program
}

func (c *Check) Build() error {
	if c == nil {
		return nil
	}
	if c.Expr == "" {
		return nil
	}
	p, err := expr.Compile(c.Expr)
	if err != nil {
		return fmt.Errorf("compile an expression: %w", err)
	}
	c.expr = p
	return nil
}

func (c *Check) GetExpr() *vm.Program {
	if c == nil {
		return nil
	}
	return c.expr
}
