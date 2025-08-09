package gosmarty

import (
	"errors"
	"strings"

	"github.com/szks-repo/gosmarty/ast"
	"github.com/szks-repo/gosmarty/lexer"
	"github.com/szks-repo/gosmarty/object"
	"github.com/szks-repo/gosmarty/parser"
)

func New(name string) *GoSmarty {
	return &GoSmarty{
		name: name,
	}
}

func (gsm *GoSmarty) Parse(input string) (*GoSmarty, error) {
	l := lexer.New(input)
	p := parser.New(l)
	tree := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}

	return &GoSmarty{
		name: gsm.name,
		tree: tree,
	}, nil
}

type GoSmarty struct {
	name string
	tree *ast.Tree
}

func (gsm *GoSmarty) Exec(env *object.Environment) object.Object {
	return Eval(gsm.tree.Root, env)
}
