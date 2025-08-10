package ast

import (
	"strings"

	"github.com/szks-repo/gosmarty/token"
)

// IndexExpression は array[index] のようなインデックスアクセスを表します
type IndexExpression struct {
	Token token.Token // The '[' token
	Left  Node        // インデックスでアクセスされる対象 (Identifier, MemberAccess など)
	Index Node        // インデックス式 (NumberLiteral, Identifier など)
}

func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IndexExpression) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}
