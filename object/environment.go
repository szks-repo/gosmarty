package object

// Environment は変数名と、それが束縛する値(Object)を保持します。
type Environment struct {
	store map[string]Object
}

// NewEnvironment は新しいEnvironmentのインスタンスを生成して返します。
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Get は指定された名前の変数を環境から取得します。
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set は環境に新しい変数の束縛を追加します。
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}