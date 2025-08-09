package object

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
	NULL  = &Null{}
)

type ObjectType int

// オブジェクトの種類を定義します。
const (
	StringType ObjectType = iota + 1
	BoolType
	NullType
)

type Object interface {
	Type() ObjectType
	// デバッグや出力のためにオブジェクトの状態を文字列で返す
	Inspect() string
}
