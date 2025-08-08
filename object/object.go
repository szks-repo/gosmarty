package object

import "fmt"

// ObjectType はオブジェクトの種類を表す文字列です。
type ObjectType string

// オブジェクトの種類を定義します。
const (
	STRING_OBJ  = "STRING"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
)

// Object は評価器が扱う全てのデータが実装すべきインターフェースです。
type Object interface {
	Type() ObjectType
	Inspect() string // デバッグや出力のためにオブジェクトの状態を文字列で返す
}

//--- 具体的なオブジェクトの実装 ---

// String は文字列を表すオブジェクトです。
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Boolean は真偽値を表すオブジェクトです。
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// Null は値が存在しないことを表すオブジェクトです。
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
	NULL  = &Null{}
)
