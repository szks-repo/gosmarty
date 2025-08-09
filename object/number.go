package object

import "fmt"

type Number struct {
	Value float64
}

func (n *Number) Type() ObjectType {
	return NumberType
}

func (n *Number) Inspect() string {
	return fmt.Sprintf("%g", n.Value)
}
