package gosmarty

import (
	"github.com/szks-repo/gosmarty/ast"
	"github.com/szks-repo/gosmarty/object"
)

var (
	NULL = &object.Null{}
)

// Eval はASTノードを評価する中心的な関数
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// ノードのリスト
	case *ast.ListNode:
		return evalNodes(node.Nodes, env)

	// アクション {$...}
	case *ast.ActionNode:
		// ActionNodeの中の式を評価する
		return Eval(node.Pipe, env)

	// テキスト
	case *ast.TextNode:
		return &object.String{Value: node.Value}

	// 識別子 (変数)
	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.IfNode:
		return evalIfNode(node, env)
	}
	return nil
}

// evalNodes はノードのスライスを評価し、結果を連結する
func evalNodes(nodes []ast.Node, env *object.Environment) object.Object {
	var result string
	for _, node := range nodes {
		evaluated := Eval(node, env)
		// 評価結果がNULLでなければ、文字列として連結する
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			result += evaluated.Inspect()
		}
	}
	return &object.String{Value: result}
}

// evalIdentifier は環境から変数の値を探して返す
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	// 簡単のため、見つからなければ空文字を返す
	// Smartyのデフォルトの動作に近い
	return &object.String{Value: ""}
}

func evalIfNode(in *ast.IfNode, env *object.Environment) object.Object {
	condition := Eval(in.Condition, env)

	if isTruthy(condition) {
		return Eval(in.Consequence, env)
	} else if in.Alternative != nil {
		return Eval(in.Alternative, env)
	} else {
		return NULL
	}
}

// isTruthy はオブジェクトが「真」であるかを判定するヘルパー
func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Null:
		return false
	case *object.String:
		return obj.Value != ""
	case *object.Boolean:
		return obj.Value
	// case *object.Number:
	//  return obj.Value != 0
	default:
		return true // NULLと空文字以外は真とみなす
	}
}
