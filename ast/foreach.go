package ast

import "github.com/szks-repo/gosmarty/token"

// ForeachNode は {foreach ...} ブロックを表します。
type ForeachNode struct {
	Token       token.Token // 'foreach' トークン
	Source      Node        // from属性で指定された反復対象
	Key         string      // key属性で指定された変数名（任意）
	Item        string      // item属性で指定された変数名
	Name        string      // name属性（現状未使用だが将来利用を考慮）
	Body        *ListNode   // foreach本体のノード
	Alternative *ListNode   // {foreachelse} ブロック
}

func (fn *ForeachNode) TokenLiteral() string {
	return fn.Token.Literal
}

func (fn *ForeachNode) String() string {
	// デバッグ用に簡易表現を返す
	return "foreach"
}
