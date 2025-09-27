package ast

import "github.com/szks-repo/gosmarty/token"

// InfixExpression represents a binary operation like `$a > $b` or `$a and $b`.
type InfixExpression struct {
	Token    token.Token
	Left     Node
	Operator string
	Right    Node
}

func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *InfixExpression) String() string {
	left := ""
	if ie.Left != nil {
		left = ie.Left.String()
	}
	right := ""
	if ie.Right != nil {
		right = ie.Right.String()
	}
	return "(" + left + " " + ie.Operator + " " + right + ")"
}
