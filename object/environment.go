package object

import "fmt"

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
		obj, err := toObject(value)
		if err != nil {
			return err
		}

		env.setVar(name, obj)
		return nil
	}
}

// Goのネイティブな型をgosmartyのObjectに変換
func toObject(i any) (Object, error) {
	switch i := i.(type) {
	case string:
		return &String{Value: i}, nil
	case int:
		return &Number{Value: float64(i)}, nil
	case int64:
		return &Number{Value: float64(i)}, nil
	case float64:
		return &Number{Value: i}, nil
	case bool:
		if i {
			return TRUE, nil
		}
		return FALSE, nil
	case nil:
		return NULL, nil
	case []any:
		elements := make([]Object, len(i))
		for idx, elem := range i {
			obj, err := toObject(elem)
			if err != nil {
				return nil, err
			}
			elements[idx] = obj
		}
		return &Array{Elements: elements}, nil
	case map[string]any:
		pairs := make(map[string]Object)
		for key, val := range i {
			obj, err := toObject(val)
			if err != nil {
				return nil, err
			}
			pairs[key] = obj
		}
		return &Map{Pairs: pairs}, nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", i)
	}
}
