package object

import "fmt"

type Boolean struct {
	Value bool
}

func NewBool[T ~bool](v T) *Boolean {
	return &Boolean{Value: bool(v)}
}

func (b *Boolean) Type() ObjectType {
	return BoolType
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
