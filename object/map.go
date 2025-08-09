package object

import (
	"bytes"
	"strings"
)

// Map はマップ（ハッシュ）オブジェクトを表します。
// Goのmapをラップします。
type Map struct {
	Pairs map[string]Object
}

func (m *Map) Type() ObjectType { return MapType }
func (m *Map) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range m.Pairs {
		pairs = append(pairs, key+":"+value.Inspect())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
