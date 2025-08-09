package eval

import (
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"

	"github.com/szks-repo/gosmarty/object"
)

type Modifier func(obj object.Object) object.Object

var msgPrinter = message.NewPrinter(language.Japanese)

var builtinModifiers = map[string]Modifier{
	"nl2br": func(input object.Object) object.Object {
		if input.Type() != object.StringType {
			return &object.String{} // またはエラーオブジェクト
		}

		str := input.(*object.String).Value
		return &object.String{Value: strings.ReplaceAll(str, "\n", "<br>")}
	},
	"number_format": func(input object.Object) object.Object {
		if input.Type() != object.NumberType {
			return &object.String{}
		}

		val := input.(*object.Number).Value
		return &object.String{Value: msgPrinter.Sprint(number.Decimal(val))}
	},
	"upper": func(input object.Object) object.Object {
		if input.Type() != object.StringType {
			return &object.String{}
		}

		return &object.String{Value: strings.ToUpper(input.Inspect())}
	},
	"devtest1": func(input object.Object) object.Object {
		if input.Type() != object.StringType {
			return &object.String{Value: ""}
		}
		str := input.(*object.String).Value
		return &object.String{Value: str + "_test1"}
	},
}
