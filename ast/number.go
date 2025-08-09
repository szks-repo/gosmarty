package ast

import "github.com/szks-repo/gosmarty/token"

// NumberLiteral は数値を表します
type NumberLiteral struct {
	Token token.Token // The token.NUMBER token
	Value float64
}

func (nl *NumberLiteral) TokenLiteral() string {
	return nl.Token.Literal
}

func (nl *NumberLiteral) String() string {
	return nl.Token.Literal
}
