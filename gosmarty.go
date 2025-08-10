package gosmarty

import (
	"errors"
	"strings"

	"github.com/szks-repo/gosmarty/ast"
	"github.com/szks-repo/gosmarty/lexer"
	"github.com/szks-repo/gosmarty/modifier"
	"github.com/szks-repo/gosmarty/object"
	"github.com/szks-repo/gosmarty/parser"
)

type GoSmarty struct {
	templates map[string]*Template
}

func New() *GoSmarty {
	return &GoSmarty{
		templates: make(map[string]*Template, 0),
	}
}

func (gsm *GoSmarty) Parse(input string) (*Template, error) {
	p := parser.New(lexer.New(input))
	tree := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}

	return &Template{
		tree: tree,
	}, nil
}

func RegisterModifier(name string, mod modifier.Modifier) {
	modifier.Register(name, mod)
}

func (gsm *GoSmarty) ExecuteTemplate(name string, env *Environment) object.Object {
	t, ok := gsm.templates[name]
	if !ok {
		panic("ERR: TODO")
	}

	return Eval(t.tree.Root, env)
}

type Template struct {
	tree *ast.Tree
}

func (t *Template) Execute(env *Environment) object.Object {
	return Eval(t.tree.Root, env)
}
