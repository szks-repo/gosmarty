package object

import (
	"strings"
)

// Map はマップ（ハッシュ）オブジェクトを表します。
// Goのmapをラップします。
type Map struct {
	Value map[string]Object
}

func (m *Map) Type() ObjectType {
	return MapType
}

func (m *Map) Inspect() string {
	var out strings.Builder

	var pairs []string
	for key, value := range m.Value {
		pairs = append(pairs, key+":"+value.Inspect())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
