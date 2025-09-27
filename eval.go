package gosmarty

import (
	"sort"
	"strings"

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
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.IfNode:
		return evalIfNode(node, env)
	case *ast.PipeNode:
		return evalPipeNode(node, env)
	case *ast.ForeachNode:
		return evalForeachNode(node, env)
	}

	return nil
}

// evalNodes はノードのスライスを評価し、結果を連結する
func evalNodes(nodes []ast.Node, env *Environment) object.Object {
	var result string
	for _, node := range nodes {
		evaluated := Eval(node, env)
		// 評価結果がNULLでなければ、文字列として連結する
	L:
		if evaluated != nil {
			if evaluated.Type() == object.OptionalType {
				evaluated = evaluated.(*object.Optional).Unwrap()
				goto L
			}
			if evaluated.Type() != object.NullType {
				result += evaluated.Inspect()
			}
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
	if val, ok := objMap.Value[propName]; ok {
		return val
	}

	return NULL
}

func evalInfixExpression(node *ast.InfixExpression, env *Environment) object.Object {
	switch node.Operator {
	case "and":
		left := unwrapOptional(Eval(node.Left, env))
		if left == nil {
			left = NULL
		}
		if !isTruthy(left) {
			return object.FALSE
		}
		right := unwrapOptional(Eval(node.Right, env))
		if right == nil {
			right = NULL
		}
		if isTruthy(right) {
			return object.TRUE
		}
		return object.FALSE
	case ">", ">=", "<", "<=", "==", "!=":
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalComparisonExpression(node.Operator, left, right)
	default:
		return NULL
	}
}

func evalComparisonExpression(op string, leftObj, rightObj object.Object) object.Object {
	left := unwrapOptional(leftObj)
	if left == nil {
		left = NULL
	}
	right := unwrapOptional(rightObj)
	if right == nil {
		right = NULL
	}

	switch op {
	case ">", ">=", "<", "<=":
		lNum, lOk := left.(*object.Number)
		rNum, rOk := right.(*object.Number)
		if !lOk || !rOk {
			return object.FALSE
		}
		var result bool
		switch op {
		case ">":
			result = lNum.Value > rNum.Value
		case ">=":
			result = lNum.Value >= rNum.Value
		case "<":
			result = lNum.Value < rNum.Value
		case "<=":
			result = lNum.Value <= rNum.Value
		}
		return boolObject(result)
	case "==", "!=":
		result := objectsEqual(left, right)
		if op == "!=" {
			result = !result
		}
		return boolObject(result)
	default:
		return NULL
	}
}

func objectsEqual(left, right object.Object) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}

	if left.Type() != right.Type() {
		return false
	}

	switch l := left.(type) {
	case *object.Null:
		return true
	case *object.Number:
		r := right.(*object.Number)
		return l.Value == r.Value
	case *object.String:
		r := right.(*object.String)
		return l.Value == r.Value
	case *object.Boolean:
		r := right.(*object.Boolean)
		return l.Value == r.Value
	default:
		return left == right
	}
}

func boolObject(val bool) object.Object {
	if val {
		return object.TRUE
	}
	return object.FALSE
}

func evalIfNode(in *ast.IfNode, env *Environment) object.Object {
	condition := Eval(in.Condition, env)

	if isTruthy(condition) {
		return Eval(in.Consequence, env)
	}

	for _, elseifNode := range in.ElseIfs {
		elseifCondition := Eval(elseifNode.Condition, env)
		if isTruthy(elseifCondition) {
			return Eval(elseifNode.Consequence, env)
		}
	}

	if in.Alternative != nil {
		return Eval(in.Alternative, env)
	}

	return NULL
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
	case *object.Map:
		return len(obj.Value) > 0
	case *object.Array:
		return len(obj.Value) > 0
	case *object.Time:
		return !obj.Value.IsZero()
	case *object.Optional:
		if obj.Some() {
			return isTruthy(obj.Unwrap())
		}
		return false
	default:
		return true
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

func evalIndexExpression(node *ast.IndexExpression, env *Environment) object.Object {
	left := Eval(node.Left, env)
	index := Eval(node.Index, env)

	// ArrayをNumberでインデックスアクセスする場合のみを考慮
	if left.Type() == object.ArrayType && index.Type() == object.NumberType {
		arrObject := left.(*object.Array)
		idx := int(index.(*object.Number).Value)

		// 範囲外アクセスチェック
		if idx < 0 || idx >= len(arrObject.Value) {
			return NULL // 範囲外ならNULLを返す
		}
		return arrObject.Value[idx]
	}

	return NULL
}

func evalForeachNode(node *ast.ForeachNode, env *Environment) object.Object {
	iterable := unwrapOptional(Eval(node.Source, env))
	if iterable == nil {
		iterable = NULL
	}

	prevItem, hadPrevItem := env.GetVar(node.Item)
	var prevKey object.Object
	var hadPrevKey bool
	if node.Key != "" {
		prevKey, hadPrevKey = env.GetVar(node.Key)
	}

	var foreachMap *object.Map
	var loopState *object.Map
	var prevLoopState object.Object
	var hadPrevLoop bool
	if node.Name != "" {
		foreachMap = ensureSmartyForeachMap(env)
		prevLoopState, hadPrevLoop = foreachMap.Value[node.Name]
		loopState = &object.Map{Value: map[string]object.Object{}}
		foreachMap.Value[node.Name] = loopState
	}

	defer func() {
		if node.Item != "" {
			if hadPrevItem {
				env.setVar(node.Item, prevItem)
			} else {
				env.unsetVar(node.Item)
			}
		}
		if node.Key != "" {
			if hadPrevKey {
				env.setVar(node.Key, prevKey)
			} else {
				env.unsetVar(node.Key)
			}
		}
		if node.Name != "" && foreachMap != nil {
			if hadPrevLoop {
				foreachMap.Value[node.Name] = prevLoopState
			} else {
				delete(foreachMap.Value, node.Name)
			}
		}
	}()

	var rendered strings.Builder
	iterated := false

	switch obj := iterable.(type) {
	case *object.Array:
		total := len(obj.Value)
		for idx, elem := range obj.Value {
			iterated = true
			env.setVar(node.Item, elem)
			if node.Key != "" {
				env.setVar(node.Key, &object.Number{Value: float64(idx)})
			}
			updateForeachLoopState(loopState, idx, total)
			appendRendered(&rendered, Eval(node.Body, env))
		}
	case *object.Map:
		if len(obj.Value) > 0 {
			keys := make([]string, 0, len(obj.Value))
			for key := range obj.Value {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			total := len(keys)
			for idx, key := range keys {
				iterated = true
				env.setVar(node.Item, obj.Value[key])
				if node.Key != "" {
					env.setVar(node.Key, object.NewString(key))
				}
				updateForeachLoopState(loopState, idx, total)
				appendRendered(&rendered, Eval(node.Body, env))
			}
		}
	}

	if iterated {
		if rendered.Len() == 0 {
			return NULL
		}
		return object.NewString(rendered.String())
	}

	if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}

	return NULL
}

func appendRendered(b *strings.Builder, obj object.Object) {
	obj = unwrapOptional(obj)
	if obj == nil {
		return
	}
	if obj.Type() == object.NullType {
		return
	}
	b.WriteString(obj.Inspect())
}

func unwrapOptional(obj object.Object) object.Object {
	for obj != nil && obj.Type() == object.OptionalType {
		opt := obj.(*object.Optional)
		if !opt.Some() {
			return NULL
		}
		obj = opt.Unwrap()
	}
	return obj
}

func ensureSmartyForeachMap(env *Environment) *object.Map {
	var smartyMap *object.Map
	if current, ok := env.GetVar("smarty"); ok {
		if existing, ok := current.(*object.Map); ok {
			smartyMap = existing
		}
	}
	if smartyMap == nil {
		smartyMap = &object.Map{Value: map[string]object.Object{}}
		env.setVar("smarty", smartyMap)
	}

	var foreachMap *object.Map
	if existing, ok := smartyMap.Value["foreach"]; ok {
		if casted, ok := existing.(*object.Map); ok {
			foreachMap = casted
		}
	}
	if foreachMap == nil {
		foreachMap = &object.Map{Value: map[string]object.Object{}}
		smartyMap.Value["foreach"] = foreachMap
	}

	return foreachMap
}

func updateForeachLoopState(loopState *object.Map, idx, total int) {
	if loopState == nil || total <= 0 {
		return
	}
	loopState.Value["first"] = object.NewBool(idx == 0)
	loopState.Value["last"] = object.NewBool(idx == total-1)
}
