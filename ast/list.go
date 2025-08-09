package ast

import "bytes"

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
