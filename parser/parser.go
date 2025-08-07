// gosmarty/parser/parser.go
package parser

import (
	"fmt"

	"github.com/szks-repo/gosmarty/ast"
	"github.com/szks-repo/gosmarty/lexer"
	"github.com/szks-repo/gosmarty/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	// 2つトークンを読み込み、curTokenとpeekTokenをセットする
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram はパース処理を開始し、ASTのルートノードを返します。
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// NOTE: この実装は非常に簡略化されています。
	// 本来はLexerがTEXTモードとTAGモードを区別し、
	// Parserはそれに応じてparseTextStatementやparseTagStatementを呼び出すべきです。
	for p.curToken.Type != token.EOF {
		var stmt ast.Statement
		if p.curToken.Type == token.LDELIM {
			stmt = p.parseTagStatement()
		} else {
			// ここではLDELIM以外のすべてをTEXTとみなす仮実装
			stmt = p.parseTextStatement()
		}

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		// nextTokenは各parseXXX関数内で呼び出されるべき
		// ここでの呼び出しは仮
		if p.curToken.Type != token.LDELIM {
			p.nextToken()
		}
	}
	return program
}

// parseTextStatement はプレーンテキストをパースします。
func (p *Parser) parseTextStatement() *ast.TextStatement {
	// この関数は、LexerがTEXTトークンを返すことを前提としています。
	// 今回の仮実装では、現在のトークンをそのまま使います。
	stmt := &ast.TextStatement{Token: p.curToken, Value: p.curToken.Literal}
	return stmt
}

// parseTagStatement は { ... } の中身をパースします。
func (p *Parser) parseTagStatement() ast.Statement {
	// { の次のトークンに基づいて処理を分岐
	switch p.peekToken.Type {
	case token.DOLLAR:
		p.nextToken() // consume {
		p.nextToken() // consume $
		return p.parseExpressionStatement()
	// case token.IF:
	// 	return p.parseIfStatement()
	default:
		// 未知のタグはエラーとしてスキップ
		msg := fmt.Sprintf("no parsing function for %s found", p.peekToken.Type)
		p.errors = append(p.errors, msg)
		p.nextToken() // { を消費
		return nil
	}
}

// parseExpressionStatement は {$foo} のような式をパースします
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken} // curToken is now '$'

	// ここで式をパースするロジック
	// 今回は識別子のみを仮実装
	stmt.Expression = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// } が来るまで進める
	for !p.curTokenIs(token.RDELIM) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	if !p.expectPeek(token.RDELIM) {
		// } がない場合はエラーだが、ここではnilを返す
		return nil
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p.Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
