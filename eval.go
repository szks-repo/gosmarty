package gosmarty

import (
	"github.com/szks-repo/gosmarty/ast"
	"github.com/szks-repo/gosmarty/modifier"
	"github.com/szks-repo/gosmarty/object"
)

var (
	NULL = object.NewNull()
)

// Eval はASTノードを評価する中心的な関数
func Eval(node ast.Node, env *Environment) object.Object {
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
		return object.NewString(node.Value)
	// 識別子 (変数)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FieldAccess:
		return evalFieldAccess(node, env)
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.IfNode:
		return evalIfNode(node, env)
	case *ast.PipeNode:
		return evalPipeNode(node, env)
	}

	return nil
}

// evalNodes はノードのスライスを評価し、結果を連結する
func evalNodes(nodes []ast.Node, env *Environment) object.Object {
	var result string
	for _, node := range nodes {
		evaluated := Eval(node, env)
		// 評価結果がNULLでなければ、文字列として連結する
		if evaluated != nil && evaluated.Type() != object.NullType {
			result += evaluated.Inspect()
		}
	}
	return object.NewString(result)
}

// evalIdentifier は環境から変数の値を探して返す
func evalIdentifier(node *ast.Identifier, env *Environment) object.Object {
	if val, ok := env.GetVar(node.Value); ok {
		return val
	}

	return NULL
}

func evalFieldAccess(node *ast.FieldAccess, env *Environment) object.Object {
	// 1. 左辺を評価する (e.g., $user -> MapObject)
	left := Eval(node.Left, env)

	// 2. 左辺がMapでなければエラー (NULLを返す)
	if left.Type() != object.MapType {
		return NULL
	}

	// 3. Mapからプロパティを取得する
	objMap := left.(*object.Map)
	propName := node.Right.Value // (e.g., "id")

	// 4. プロパティが存在すればその値を、なければNULLを返す
	if val, ok := objMap.Pairs[propName]; ok {
		return val
	}

	return NULL
}

func evalIfNode(in *ast.IfNode, env *Environment) object.Object {
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
	case *object.Number:
		return obj.Value != 0
	default:
		return true // NULLと空文字以外は真とみなす
	}
}

func evalPipeNode(node *ast.PipeNode, env *Environment) object.Object {
	// 1. 左辺を評価する
	left := Eval(node.Left, env)

	funcName := node.Function.Value
	fn, ok := modifier.Get(funcName)
	if !ok {
		// エラー処理: 未定義の関数
		// ここでは空文字を返す
		return NULL
	}

	return fn(left)
}
