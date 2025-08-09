package ast

import (
	"github.com/szks-repo/gosmarty/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

// Tree はパースされたテンプレート全体を表すASTのルートです。
// text/template の parse.Tree に相当します。
type Tree struct {
	Name string    // テンプレート名（将来的な利用のため）
	Root *ListNode // ノードツリーのルート
}

func (t *Tree) String() string {
	return t.Root.String()
}

// ActionNode は評価されるべきアクション（例: {$name}）を表します。
// `{{...}}` に相当します。
type ActionNode struct {
	Token token.Token // The '{' (LDELIM) token
	Pipe  Node        // 評価されるべき式のパイプライン（将来の拡張用）
}

func (an *ActionNode) TokenLiteral() string { return an.Token.Literal }
func (an *ActionNode) String() string {
	if an.Pipe != nil {
		return "{" + an.Pipe.String() + "}"
	}
	return ""
}
