package ast

import (
	"bytes"

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

// ListNode はノードのシーケンス（リスト）を表します。
// テンプレートはTextNodeとActionNodeなどの連続したリストと見なせます。
type ListNode struct {
	Nodes []Node // 子ノードのリスト
}

func (ln *ListNode) TokenLiteral() string {
	if len(ln.Nodes) > 0 {
		return ln.Nodes[0].TokenLiteral()
	}
	return ""
}
func (ln *ListNode) String() string {
	var out bytes.Buffer
	for _, n := range ln.Nodes {
		out.WriteString(n.String())
	}
	return out.String()
}

// TextNode はプレーンなテキストを表します。
type TextNode struct {
	Token token.Token // The token.TEXT token
	Value string
}

func (tn *TextNode) TokenLiteral() string { return tn.Token.Literal }
func (tn *TextNode) String() string       { return tn.Value }

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

// Identifier は変数 ($foo) を表します。
type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string
}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return "$" + i.Value }

// IfNode は {if ...} ... {else} ... {/if} 構文を表します
type IfNode struct {
	Token       token.Token // The 'if' token
	Condition   Node        // ifの条件式
	Consequence *ListNode   // 条件が真の場合に実行されるノードリスト
	Alternative *ListNode   // 条件が偽の場合に実行されるノードリスト (else, elseif)
}

func (in *IfNode) TokenLiteral() string { return in.Token.Literal }
func (in *IfNode) String() string       { /* デバッグ用の実装 */ return "if" }

// PipeNode は {$left | right} のようなパイプライン式を表します
type PipeNode struct {
	Token    token.Token // The '|' token
	Left     Node        // パイプの左辺（値を提供する式）
	Function *Identifier // 適用する関数（修飾子）
}

func (pn *PipeNode) TokenLiteral() string { return pn.Token.Literal }
func (pn *PipeNode) String() string {
	// デバッグ用の実装
	return "(" + pn.Left.String() + " | " + pn.Function.String() + ")"
}
