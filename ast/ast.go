// gosmarty/ast/ast.go
package ast

import (
	"bytes"

	"github.com/szks-repo/gosmarty/token"
)

// Node は全てのASTノードが実装すべき基本インターフェースです。
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement は文を表すノードです。
type Statement interface {
	Node
	statementNode()
}

// Expression は式を表すノードです。
type Expression interface {
	Node
	expressionNode()
}

// Program はASTのルートノードです。
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// TextStatement はプレーンなテキストを表します。
type TextStatement struct {
	Token token.Token // The token.TEXT token
	Value string
}

func (ts *TextStatement) statementNode()       {}
func (ts *TextStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *TextStatement) String() string       { return ts.Value }

// ExpressionStatement は {$foo} のような式文を表します。
type ExpressionStatement struct {
	Token      token.Token // The '{' token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return "{" + es.Expression.String() + "}"
	}
	return ""
}

// Identifier は変数 ($foo) を表します。
type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return "$" + i.Value }
