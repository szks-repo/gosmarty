package ast

import "github.com/szks-repo/gosmarty/token"

// Identifier は変数 ($foo) を表します。
type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return "$" + i.Value
}
