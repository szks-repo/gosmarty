package object

import (
	"bytes"
	"strings"
)

// Array は配列オブジェクトを表します。
// Goのsliceをラップします。
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ArrayType }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := make([]string, len(a.Elements))
	for i, e := range a.Elements {
		elements[i] = e.Inspect()
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
