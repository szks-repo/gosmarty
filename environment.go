package gosmarty

import (
	"github.com/szks-repo/gosmarty/modifier"
	"github.com/szks-repo/gosmarty/object"
)

// Environment は変数名と、それが束縛する値(Object)を保持します。
type Environment struct {
	vars      map[string]object.Object
	modifiers map[string]modifier.Modifier
}

func NewEnvironment(opt ...EnvOption) (*Environment, error) {
	env := &Environment{
		vars: make(map[string]object.Object),
	}
	for _, fn := range opt {
		_ = fn(env) //todo: reporting all errors
	}

	return env, nil
}

func (e *Environment) GetVar(name string) (object.Object, bool) {
	obj, ok := e.vars[name]
	return obj, ok
}

func (e *Environment) setVar(name string, val object.Object) {
	e.vars[name] = val
}

type EnvOption = func(env *Environment) error

func WithVariable(name string, value any) EnvOption {
	return func(env *Environment) error {
		obj, err := object.NewObjectFromAny(value)
		if err != nil {
			return err
		}

		env.setVar(name, obj)
		return nil
	}
}

func WithStringVariable(name, value string) EnvOption {
	return func(env *Environment) (err error) {
		env.setVar(name, object.NewString(value))
		return
	}
}

func WithBoolVariable(name string, value bool) EnvOption {
	return func(env *Environment) (err error) {
		env.setVar(name, object.NewBool(value))
		return
	}
}
