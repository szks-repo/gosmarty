package ast

import "github.com/szks-repo/gosmarty/token"

// PipeNode は {$left | right} のようなパイプライン式を表します
type PipeNode struct {
	Token    token.Token // The '|' token
	Left     Node        // パイプの左辺（値を提供する式）
	Function *Identifier // 適用する関数（修飾子）
}

func (pn *PipeNode) TokenLiteral() string {
	return pn.Token.Literal
}

func (pn *PipeNode) String() string {
	// デバッグ用の実装
	return "(" + pn.Left.String() + " | " + pn.Function.String() + ")"
}
