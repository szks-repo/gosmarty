package ast

import (
	"strings"

	"github.com/szks-repo/gosmarty/token"
)

// MemberAccess は {$obj.prop} のようなプロパティアクセスを表します
type FieldAccess struct {
	Token token.Token // The '.' token
	Left  Node        // ドットの左側にあるオブジェクト (Identifier or another FieldAccess)
	Right *Identifier // アクセスされるプロパティ
}

func (ma *FieldAccess) TokenLiteral() string {
	return ma.Token.Literal
}

func (ma *FieldAccess) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(ma.Left.String())
	out.WriteString(".")
	out.WriteString(ma.Right.Value) // プロパティ名は '$' なしで表示
	out.WriteString(")")

	return out.String()
}
