package modifier

import (
	"strings"
	"sync"

	phpstring "github.com/szks-repo/go-php-functions/string"

	"github.com/szks-repo/gosmarty/object"
)

type Modifier func(input object.Object, args ...any) object.Object

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
var registry = map[string]Modifier{
	"nl2br": func(input object.Object, args ...any) object.Object {
		if input.Type() != object.StringType {
			return object.NULL // またはエラーオブジェクト
		}

		return object.NewString(strings.ReplaceAll(input.Inspect(), "\n", "<br>"))
	},
	"number_format": func(input object.Object, args ...any) object.Object {
		if input.Type() != object.NumberType {
			return object.NULL
		}

		val := input.(*object.Number).Value
		return object.NewString(phpstring.NumberFormat[float64](val))
	},
	"upper": func(input object.Object, args ...any) object.Object {
		if input.Type() != object.StringType {
			return object.NULL
		}

		return object.NewString(strings.ToUpper(input.Inspect()))
	},
	"lower": func(input object.Object, args ...any) object.Object {
		if input.Type() != object.StringType {
			return object.NULL
		}

		return object.NewString(strings.ToLower(input.Inspect()))
	},
	"wordwrap": func(input object.Object, args ...any) object.Object {
		if input.Type() != object.StringType {
			return object.NULL
		}

		// TODO:
		// parse args
		// 0: width
		// 1: breakWord
		// 2: cutLong
		var opt phpstring.WordwrapOpt

		return object.NewString(phpstring.Wordwrap(input.Inspect(), opt))
	},
}

var registryMu sync.RWMutex

func Get(name string) (Modifier, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	mod, ok := registry[name]
	return mod, ok
}

func Register(name string, mod Modifier) bool {
	registryMu.Lock()
	defer registryMu.Unlock()
	_, overrided := registry[name]
	registry[name] = mod

	return overrided
}
