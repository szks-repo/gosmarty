package object

import "fmt"

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BoolType
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
