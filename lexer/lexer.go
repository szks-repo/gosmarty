package lexer

import (
	"github.com/szks-repo/gosmarty/token"
)

type Lexer struct {
	input        string
	position     int  // 現在の文字の位置
	readPosition int  // 次の文字の位置
	ch           byte // 現在検査中の文字
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // 0はASCIIのNUL文字で、EOFを表す
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken は入力ソースから次のトークンを読み取って返します。
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '{':
		// コメント {* ... *} のハンドリング
		if l.peekChar() == '*' {
			l.readChar() // '*' を消費
			tok.Type = token.COMMENT
			tok.Literal = l.readComment()
		} else {
			tok = newToken(token.LDELIM, l.ch)
		}
	case '}':
		tok = newToken(token.RDELIM, l.ch)
	case '$':
		tok = newToken(token.DOLLAR, l.ch)
	case '"', '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString(l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			// ここでTEXTかIDENTかを判断する必要がある
			// 簡単のため、ここでは識別子として読み取る
			literal := l.readIdentifier()
			tok.Type = token.LookupIdent(literal)
			tok.Literal = literal
			return tok // readIdentifier内でreadCharを呼んでいるため、早期リターン
		} else if isDigit(l.ch) {
			tok.Type = token.NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			// 簡単のため、今はTEXTトークンを単純化
			// 本来は { が現れるまでをTEXTとして読み取るロジックが必要
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString(quote byte) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == quote || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readComment() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '*' && l.peekChar() == '}' {
			break
		}
		if l.ch == 0 { // EOF
			break
		}
	}
	commentBody := l.input[position:l.position]
	l.readChar() // '*' を消費
	l.readChar() // '}' を消費
	return commentBody
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '/'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
