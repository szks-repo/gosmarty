package ast

import "github.com/szks-repo/gosmarty/token"

// TextNode はプレーンなテキストを表します。
type TextNode struct {
	Token token.Token // The token.TEXT token
	Value string
}

func (tn *TextNode) TokenLiteral() string {
	return tn.Token.Literal
}

func (tn *TextNode) String() string {
	return tn.Value
}
