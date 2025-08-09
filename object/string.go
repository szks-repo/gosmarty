package object

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type String struct {
	Value string
}

func StringFromInteger[I constraints.Integer](i I) String {
	return String{Value: fmt.Sprintf("%d", i)}
}

func (s *String) Type() ObjectType {
	return StringType
}

func (s *String) Inspect() string {
	return s.Value
}
