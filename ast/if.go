package ast

import "github.com/szks-repo/gosmarty/token"

// IfNode は {if ...} ... {else} ... {/if} 構文を表します
type IfNode struct {
	Token       token.Token // The 'if' token
	Condition   Node        // ifの条件式
	Consequence *ListNode   // 条件が真の場合に実行されるノードリスト
	Alternative *ListNode   // 条件が偽の場合に実行されるノードリスト (else, elseif)
}

func (in *IfNode) TokenLiteral() string {
	return in.Token.Literal
}

func (in *IfNode) String() string {
	/* デバッグ用の実装 */
	return "if"
}

//workflow test