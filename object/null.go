package object

type Null struct{}

func NewNull() *Null {
	return &Null{}
}

func (n *Null) Type() ObjectType {
	return NullType
}

func (n *Null) Inspect() string {
	return "null"
}
