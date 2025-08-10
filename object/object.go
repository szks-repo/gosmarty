package object

import (
	"fmt"
	"reflect"
	"time"
)

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
	NULL  = &Null{}
)

type ObjectType int

// オブジェクトの種類を定義します。
const (
	StringType ObjectType = iota + 1
	BoolType
	NullType
	NumberType
	ArrayType
	MapType
	TimeType
)

type Object interface {
	Type() ObjectType
	// デバッグや出力のためにオブジェクトの状態を文字列で返す
	Inspect() string
}

// Goのネイティブな型をgosmartyのObjectに変換
func NewObjectFromAny(i any) (Object, error) {
	switch i := i.(type) {
	case string:
		return &String{Value: i}, nil
	case int:
		return &Number{Value: float64(i)}, nil
	case int64:
		return &Number{Value: float64(i)}, nil
	case uint:
		return &Number{Value: float64(i)}, nil
	case uint64:
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
	case []string:
		values := make([]Object, len(i))
		for idx, elem := range i {
			values[idx] = NewString(elem)
		}
		return &Array{Value: values}, nil
	case []any:
		values := make([]Object, len(i))
		for idx, elem := range i {
			obj, err := NewObjectFromAny(elem)
			if err != nil {
				return nil, err
			}
			values[idx] = obj
		}
		return &Array{Value: values}, nil
	case map[string]any:
		pairs := make(map[string]Object)
		for key, val := range i {
			obj, err := NewObjectFromAny(val)
			if err != nil {
				return nil, err
			}
			pairs[key] = obj
		}
		return &Map{Value: pairs}, nil
	case time.Time:
		return NewTime(i), nil
	// todo: support go stdlib package types
	// case *big.Rat:
	// case *big.Int:
	default:
		rv := reflect.ValueOf(i)
		// underlying types or structs
		switch rv.Kind() {
		case reflect.String:
			return &String{Value: rv.String()}, nil
		case reflect.Int:
			return &Number{Value: float64(rv.Int())}, nil
		case reflect.Int64:
			return &Number{Value: float64(rv.Int())}, nil
		case reflect.Uint:
			return &Number{Value: float64(rv.Uint())}, nil
		case reflect.Uint64:
			return &Number{Value: float64(rv.Uint())}, nil
		case reflect.Bool:
			return &Boolean{Value: rv.Bool()}, nil
		case reflect.Struct:
			//todo
		}

		return nil, fmt.Errorf("unsupported type: %T", i)
	}
}
