package object

// Environment は変数名と、それが束縛する値(Object)を保持します。
type Environment struct {
	store map[string]Object
}

func NewEnvironment(opt ...EnvOption) (*Environment, error) {
	env := &Environment{
		store: make(map[string]Object),
	}
	for _, fn := range opt {
		_ = fn(env) //todo: reporting all errors
	}

	return env, nil
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) setVar(name string, val Object) {
	e.store[name] = val
}

type EnvOption = func(env *Environment) error

func WithVariable(name string, value any) EnvOption {
	return func(env *Environment) error {
		obj, err := newObjectFromAny(value)
		if err != nil {
			return err
		}

		env.setVar(name, obj)
		return nil
	}
}
