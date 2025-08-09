package eval

import (
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"

	"github.com/szks-repo/gosmarty/object"
)

type Modifier func(object.Object, ...any) object.Object

var msgPrinter = message.NewPrinter(language.Japanese)

// Builtin variable modifiers
// https://www.smarty.net/docs/en/language.modifiers.tpl
// - capitalize
// - cat
// - count_characters
// - count_paragraphs
// - count_sentences
// - count_words
// - date_format
// - default
// - escape
// - from_charset
// - indent
// - lower
// - nl2br
// - regex_replace
// - replace
// - spacify
// - string_format
// - strip
// - strip_tags
// - to_charset
// - truncate
// - unescape
// - upper
// - wordwrap
var builtinModifiers = map[string]Modifier{
	"nl2br": func(input object.Object, args ...any) object.Object {
		if input.Type() != object.StringType {
			return &object.String{} // またはエラーオブジェクト
		}

		str := input.(*object.String).Value
		return &object.String{Value: strings.ReplaceAll(str, "\n", "<br>")}
	},
	"number_format": func(input object.Object, args ...any) object.Object {
		if input.Type() != object.NumberType {
			return &object.String{}
		}

		val := input.(*object.Number).Value
		return &object.String{Value: msgPrinter.Sprint(number.Decimal(val))}
	},
	"upper": func(input object.Object, args ...any) object.Object {
		if input.Type() != object.StringType {
			return &object.String{}
		}

		return &object.String{Value: strings.ToUpper(input.Inspect())}
	},
	"lower": func(input object.Object, args ...any) object.Object {
		if input.Type() != object.StringType {
			return &object.String{}
		}

		return &object.String{Value: strings.ToLower(input.Inspect())}
	},
	"devtest1": func(input object.Object, args ...any) object.Object {
		if input.Type() != object.StringType {
			return &object.String{Value: ""}
		}
		str := input.(*object.String).Value
		return &object.String{Value: str + "_test1"}
	},
}
