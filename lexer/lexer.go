package lexer

import (
	"unicode"

	"github.com/szks-repo/gosmarty/token"
)

// Lexerの状態を定義
type lexerState int

const (
	stateText lexerState = iota // デリミタ外のテキストを解析中
	stateTag                    // デリミタ内のタグを解析中
)

type Lexer struct {
	input   []rune
	pos     int
	readPos int
	ch      rune
	state   lexerState
}

func New(input string) *Lexer {
	l := &Lexer{
		input: []rune(input),
		state: stateText,
	}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	if l.state == stateText {
		return l.nextTokenInText()
	}
	return l.nextTokenInTag()
}

// stateText時のトークン生成
func (l *Lexer) nextTokenInText() token.Token {
	var tok token.Token
	// `{` が見つかるか、入力が終わるまでを読む
	pos := l.pos
	for l.ch != '{' && l.ch != 0 {
		l.readChar()
	}

	// `{` の前の文字列をTEXTトークンとして返す
	if l.pos > pos {
		tok.Type = token.TEXT
		tok.Literal = string(l.input[pos:l.pos])
		return tok
	}

	// `{` が見つかった場合
	if l.ch == '{' {
		// タグモードに移行
		l.state = stateTag
		return l.nextTokenInTag()
	}

	return token.Token{Type: token.EOF, Literal: ""}
}

// stateTag時のトークン生成（元のNextTokenのロジックに近い）
func (l *Lexer) nextTokenInTag() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '{':
		if l.peekChar() == '*' { // コメント {* ... *}
			l.readChar()
			tok.Type = token.COMMENT
			tok.Literal = l.readComment()
		} else {
			tok = newToken(token.LDELIM, l.ch)
		}
	case '}':
		tok = newToken(token.RDELIM, l.ch)
		l.state = stateText // テキストモードに復帰
	case '$':
		tok = newToken(token.DOLLAR, l.ch)
	case '"', '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString(l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if unicode.IsLetter(l.ch) || l.ch == '/' { // /if, /foreach のため
			literal := l.readIdentifier()
			tok.Type = token.LookupIdent(literal)
			tok.Literal = literal
			return tok
		} else if unicode.IsDigit(l.ch) {
			tok.Type = token.NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos
	// `/` で始まる場合（終了タグ）
	if l.ch == '/' {
		l.readChar()
	}
	for unicode.IsLetter(l.ch) || unicode.IsDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}

	return string(l.input[pos:l.pos])
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) peekChar() rune {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for unicode.IsDigit(l.ch) {
		l.readChar()
	}
	return string(l.input[pos:l.pos])
}

func (l *Lexer) readString(quote rune) string {
	pos := l.pos + 1
	for {
		l.readChar()
		if l.ch == quote || l.ch == 0 {
			break
		}
	}
	return string(l.input[pos:l.pos])
}

func (l *Lexer) readComment() string {
	pos := l.pos + 1
	for {
		l.readChar()
		if l.ch == '*' && l.peekChar() == '}' {
			break
		}
		if l.ch == 0 {
			break
		}
	}
	commentBody := l.input[pos:l.pos]
	l.readChar()
	l.readChar()
	return string(commentBody)
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
