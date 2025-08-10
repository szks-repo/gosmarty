package object

import (
	"strings"
)

// Array は配列オブジェクトを表します。
// Goのsliceをラップします。
type Array struct {
	Value []Object
}

func (a *Array) Type() ObjectType {
	return ArrayType
}

func (a *Array) Inspect() string {
	var out strings.Builder

	elements := make([]string, len(a.Value))
	for i, e := range a.Value {
		elements[i] = e.Inspect()
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
