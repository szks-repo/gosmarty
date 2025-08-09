package gosmarty

import (
	"strings"

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
	case *ast.PipeNode:
		return evalPipeNode(node, env)
	}

	return nil
}

// evalNodes はノードのスライスを評価し、結果を連結する
func evalNodes(nodes []ast.Node, env *object.Environment) object.Object {
	var result string
	for _, node := range nodes {
		evaluated := Eval(node, env)
		// 評価結果がNULLでなければ、文字列として連結する
		if evaluated != nil && evaluated.Type() != object.NullType {
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

// Builtin はテンプレート内で使用可能なGoの関数の型
type Builtin func(input object.Object) object.Object

var builtins = map[string]Builtin{
	"nl2br": func(input object.Object) object.Object {
		if input.Type() != object.StringType {
			return &object.String{Value: ""} // またはエラーオブジェクト
		}

		str := input.(*object.String).Value
		return &object.String{Value: strings.ReplaceAll(str, "\n", "<br>")}
	},
	"devtest1": func(input object.Object) object.Object {
		if input.Type() != object.StringType {
			return &object.String{Value: ""}
		}
		str := input.(*object.String).Value
		return &object.String{Value: str + "_test1"}
	},
}

func evalPipeNode(node *ast.PipeNode, env *object.Environment) object.Object {
	// 1. 左辺を評価する
	left := Eval(node.Left, env)

	// 2. 関数名を取得
	funcName := node.Function.Value

	// 3. 関数レジストリから関数を探す
	fn, ok := builtins[funcName]
	if !ok {
		// エラー処理: 未定義の関数
		// ここでは空文字を返す
		return &object.String{Value: ""}
	}

	// 4. 関数を実行して結果を返す
	return fn(left)
}
