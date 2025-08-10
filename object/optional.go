package object

type Optional struct {
	value Object
	some  bool
}

func NewOptional(obj Object) *Optional {
	if obj == nil {
		return &Optional{value: NULL}
	}
	if obj.Type() == NullType {
		return &Optional{value: NULL}
	}

	return &Optional{value: obj, some: true}
}

func (o *Optional) Type() ObjectType {
	return OptionalType
}

func (o *Optional) Inspect() string {
	if o.some {
		return o.value.Inspect()
	}
	return "null"
}

func (o *Optional) Unwrap() Object {
	if !o.some {
		return NewNull()
	}
	return o.value
}

func (o Optional) Some() bool {
	return o.some
}
