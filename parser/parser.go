package parser

import (
	"fmt"
	"strconv"

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

func (p *Parser) ParseProgram() *ast.Tree {
	tree := &ast.Tree{
		Root: &ast.ListNode{Nodes: []ast.Node{}},
	}

	for p.curToken.Type != token.EOF {
		var node ast.Node
		switch p.curToken.Type {
		case token.TEXT:
			node = p.parseTextNode()
		case token.LDELIM:
			// '{' を見つけたら、次のトークンを覗き見てどのタグか判断する
			// If it's a comment, consume it and continue
			if p.peekToken.Type == token.COMMENT {
				p.nextToken() // Consume LDELIM
				p.nextToken() // Consume COMMENT
				// コメントをスキップした後、現在のトークンがTEXTであれば、それを処理する
			}
			node = p.parseTag()
		default:
			p.nextToken()
			continue
		}

		if node != nil {
			tree.Root.Nodes = append(tree.Root.Nodes, node)
		}
	}
	return tree
}

func (p *Parser) parseTextNode() *ast.TextNode {
	node := &ast.TextNode{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken() // TEXTトークンを消費
	return node
}

// parseTag は `{` の次のトークンを見て、どの構文か判断し、パースを振り分ける
func (p *Parser) parseTag() ast.Node {
	switch p.peekToken.Type {
	case token.DOLLAR:
		return p.parseVariableTagWithPipeline()
	// todo: consider this case
	// case token.NUMBER:
	// return p.parseVariableTagWithPipeline()
	case token.IF:
		return p.parseIfTag()
	default:
		// エラー処理：不明なタグ
		p.errors = append(p.errors, fmt.Sprintf("unknown tag type: %s", p.peekToken.Type))
		p.nextToken() // エラーリカバリーのため進める
		return nil
	}
}

// parseVariableTag は {$name} のような変数をパースする
func (p *Parser) parseVariableTag() *ast.ActionNode {
	// curTokenは '{'
	node := &ast.ActionNode{Token: p.curToken}
	// '{' を消費 -> curTokenは '$'
	p.nextToken()

	// '$' を消費 -> curTokenは 識別子
	p.nextToken()
	if !p.curTokenIs(token.IDENT) {
		p.errors = append(p.errors, fmt.Sprintf("expected IDENT, got %s", p.curToken.Type))
		return nil
	}
	node.Pipe = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	// 識別子を消費
	p.nextToken()

	if !p.curTokenIs(token.RDELIM) {
		p.errors = append(p.errors, fmt.Sprintf("expected RDELIM, got %s", p.curToken.Type))
		return nil
	}

	// '}' を消費
	p.nextToken()
	return node
}

// parseIfTag は {if ...} ブロック全体をパースする
func (p *Parser) parseIfTag() *ast.IfNode {
	// curTokenは '{'
	p.nextToken() // '{' を消費 -> curTokenは 'if'
	node := &ast.IfNode{Token: p.curToken}
	p.nextToken() // 'if' を消費 -> curTokenは '$'

	// 1. 条件式のパース
	if p.curToken.Type != token.DOLLAR {
		p.errors = append(p.errors, "expected variable expression for if condition")
		return nil
	}

	// '$' を消費 -> curTokenは 識別子
	p.nextToken()
	if !p.curTokenIs(token.IDENT) {
		p.errors = append(p.errors, fmt.Sprintf("expected IDENT for condition, got %s", p.curToken.Type))
		return nil
	}
	node.Condition = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	// 識別子を消費
	p.nextToken()

	// 2. {if ...} の閉じ '}' を消費
	if !p.curTokenIs(token.RDELIM) {
		p.errors = append(p.errors, fmt.Sprintf("expected RDELIM after if condition, got %s", p.curToken.Type))
		return nil
	}
	// '}' を消費
	p.nextToken()

	// 3. Consequence (then節) のパース
	node.Consequence = p.parseBlockUntil(token.ELSE, token.ENDIF)

	// 4. Alternative (else節) のパース (存在すれば)
	if p.curTokenIs(token.LDELIM) && p.peekTokenIs(token.ELSE) {
		// '{' を消費
		p.nextToken()
		// 'else' を消費
		p.nextToken()

		if !p.curTokenIs(token.RDELIM) {
			p.errors = append(p.errors, "expected RDELIM for else tag")
			return nil
		}
		// '}' を消費
		p.nextToken()

		node.Alternative = p.parseBlockUntil(token.ENDIF)
	}

	// 5. 終了タグ {/if} を消費
	if !(p.curTokenIs(token.LDELIM) && p.peekTokenIs(token.ENDIF)) {
		p.errors = append(p.errors, "expected {/if} tag")
		return nil
	}
	// '{' を消費
	p.nextToken()
	// '/if' を消費
	p.nextToken()

	if !p.curTokenIs(token.RDELIM) {
		p.errors = append(p.errors, "expected RDELIM for /if tag")
		return nil
	}
	// '}' を消費
	p.nextToken()

	return node
}

// parseBlockUntil は指定された終了トークンが見つかるまでノードをパースし続ける
func (p *Parser) parseBlockUntil(endTokens ...token.TokenType) *ast.ListNode {
	block := &ast.ListNode{Nodes: []ast.Node{}}

	for {
		if p.curTokenIs(token.EOF) {
			p.errors = append(p.errors, "unexpected EOF, unclosed block")
			return block
		}
		if p.curTokenIs(token.LDELIM) {
			for _, endToken := range endTokens {
				if p.peekTokenIs(endToken) {
					return block
				}
			}
		}

		var stmt ast.Node
		switch p.curToken.Type {
		case token.TEXT:
			stmt = p.parseTextNode()
		case token.LDELIM:
			stmt = p.parseTag()
		default:
			p.nextToken()
			continue
		}
		if stmt != nil {
			block.Nodes = append(block.Nodes, stmt)
		}
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) parseVariableTagWithPipeline() ast.Node {
	// curTokenは '{'
	p.nextToken() // '{' を消費 -> curTokenは '$'

	// 最初に左辺の式をパース
	left := p.parsePrimaryExpr()

	// '|' が続く限りパイプラインを構築
	for p.curTokenIs(token.PIPE) {
		pipeToken := p.curToken

		// '|' を消費
		p.nextToken()

		if !p.curTokenIs(token.IDENT) {
			p.errors = append(p.errors, "expected modifier function name after '|'")
			return nil
		}

		// 新しいPipeNodeを作成し、それまでの式を左辺に設定
		left = &ast.PipeNode{
			Token: pipeToken,
			Left:  left,
			Function: &ast.Identifier{
				Token: p.curToken,
				Value: p.curToken.Literal,
			},
		}
		// 関数名を消費
		p.nextToken()
	}

	if !p.curTokenIs(token.RDELIM) {
		p.errors = append(p.errors, fmt.Sprintf("expected RDELIM, got %s", p.curToken.Type))
		return nil
	}

	// '}' を消費
	p.nextToken()

	return &ast.ActionNode{
		Token: token.Token{
			Type:    token.LDELIM,
			Literal: "{",
		},
		Pipe: left,
	}
}

// パイプラインの元となる最初の式をパースするヘルパー
func (p *Parser) parsePrimaryExpr() ast.Node {
	switch p.curToken.Type {
	case token.DOLLAR:
		// '
		p.nextToken()
		if !p.curTokenIs(token.IDENT) {
			p.errors = append(p.errors, fmt.Sprintf("expected IDENT, got %s", p.curToken.Type))
			return nil
		}
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		// 識別子を消費
		p.nextToken()
		return ident
	case token.NUMBER:
		return p.parseNumberLiteral()
	default:
		// 他のプライマリー式（文字列リテラルなど）もここに追加できる
		return nil
	}
}

func (p *Parser) parseNumberLiteral() ast.Node {
	lit := &ast.NumberLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float64", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	// 数値トークンを消費
	p.nextToken()
	return lit
}
