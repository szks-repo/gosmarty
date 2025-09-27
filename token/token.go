package token

// Token は字句解析器(Lexer)が生成するトークンを表す構造体です。
type Token struct {
	Type    TokenType
	Literal string
}

// TokenType はトークンの種類を表す文字列です。
type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Delimiters
	LBRACE     = "{"        // {
	RBRACE     = "}"        // }
	LDELIM     = "LDELIM"   // Smartyの左デリミタ (e.g., {)
	RDELIM     = "RDELIM"   // Smartyの右デリミタ (e.g., })
	LITERAL    = "LITERAL"  // {literal}
	ENDLITERAL = "/LITERAL" // {/literal}
	COMMENT    = "COMMENT"  // {* ... *}

	// 識別子 + リテラル
	IDENT    = "IDENT" // 変数名など (例: foo, bar)
	DOLLAR   = "$"
	PIPE     = "|"
	DOT      = "."
	LBRACKET = "["
	RBRACKET = "]"
	STRING   = "STRING" // "foo" or 'bar'
	NUMBER   = "NUMBER" // 12345
	TEXT     = "TEXT"   // デリミタの外にあるプレーンなテキスト

	// 演算子
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	EQ    = "=="
	NOTEQ = "!="
	LT    = "<"
	LTE   = "<="
	GT    = ">"
	GTE   = ">="
	AND   = "and"

	IF          = "if"
	ELSE        = "else"
	ELSEIF      = "elseif"
	ENDIF       = "/if"
	FOREACH     = "foreach"
	FOREACHELSE = "foreachelse"
	ENDFOREACH  = "/foreach"
)

var keywords = map[string]TokenType{
	"if":          IF,
	"else":        ELSE,
	"elseif":      ELSEIF,
	"/if":         ENDIF,
	"foreach":     FOREACH,
	"foreachelse": FOREACHELSE,
	"/foreach":    ENDFOREACH,
	"literal":     LITERAL,
	"/literal":    ENDLITERAL,
	"and":         AND,
}

// LookupIdent は識別子がキーワードかどうかを判定します。
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
