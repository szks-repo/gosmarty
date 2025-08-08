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

// ParseProgram はパース処理を開始し、*ast.Tree を返します。
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
			node = p.parseActionNode()
		default:
			// LDELIMでもTEXTでもないトークンは無視して進む
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

// parseActionNode は {$...} をパースします
func (p *Parser) parseActionNode() ast.Node {
	// 現在のトークンは LDELIM ({)
	node := &ast.ActionNode{Token: p.curToken}

	p.nextToken() // { を消費

	// タグの中身をパース
	switch p.curToken.Type {
	case token.DOLLAR:
		node.Pipe = p.parsePipe()
	case token.IF,token.ENDIF,token.ELSEIF,token.ELSE:
		return p.parseIfNode()
	default:
		msg := fmt.Sprintf("unexpected token in tag: got %s", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	// } を期待
	if !p.curTokenIs(token.RDELIM) {
		msg := fmt.Sprintf("expected RDELIM, got %s instead", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	p.nextToken() // } を消費

	return node
}

// parsePipe は $name のような式をパースします
func (p *Parser) parsePipe() ast.Node {
	// 現在のトークンは $
	p.nextToken() // $ を消費

	if !p.curTokenIs(token.IDENT) {
		return nil // エラー処理
	}

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // 識別子を消費
	return ident
}

// ... (その他のヘルパー関数 curTokenIs, peekTokenIs, expectPeek, peekError は変更なし) ...
// NOTE: curTokenIs, peekTokenIs, expectPeek, peekError といった
// ヘルパー関数は変更がないため、元のコードをそのまま使用してください。
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
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
	msg := fmt.Sprintf("expected next token to be %s, got %s(%s) instead",
		t, p.peekToken.Type, p.peekToken.Literal)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseIfNode() *ast.IfNode {
	// 現在のトークンは 'if'
	node := &ast.IfNode{Token: p.curToken}

	p.nextToken() // 'if' を消費

	// --- 条件式のパース ---
	// ここでは簡単のため、{$foo} のみ対応
	if p.curToken.Type == token.DOLLAR {
		node.Condition = p.parsePipe()
	} else {
		// エラー処理
		return nil
	}

	// --- Consequence (then節) のパース ---
	node.Consequence = &ast.ListNode{}
	// {/if} または {else} が来るまでパースを続ける
	for !p.curTokenIs(token.ENDIF) && !p.curTokenIs(token.ELSE) && !p.curTokenIs(token.EOF) {
		var stmt ast.Node
		switch p.curToken.Type {
		case token.TEXT:
			stmt = p.parseTextNode()
		case token.LDELIM:
			// {if} の中にも {$var} のようなタグは書ける
			stmt = p.parseActionNode()
		default:
			p.nextToken()
			continue
		}
		node.Consequence.Nodes = append(node.Consequence.Nodes, stmt)
	}

	// --- Alternative (else節) のパース ---
	if p.curTokenIs(token.ELSE) {
		p.nextToken() // 'else' を消費

		node.Alternative = &ast.ListNode{}
		// {/if} が来るまでパースを続ける
		for !p.curTokenIs(token.ENDIF) && !p.curTokenIs(token.EOF) {
			var stmt ast.Node
			switch p.curToken.Type {
			case token.TEXT:
				stmt = p.parseTextNode()
			case token.LDELIM:
				stmt = p.parseActionNode()
			default:
				p.nextToken()
				continue
			}
			node.Alternative.Nodes = append(node.Alternative.Nodes, stmt)
		}
	}

	if !p.curTokenIs(token.ENDIF) {
		// {/if} がない場合はエラー
		p.peekError(token.ENDIF)
		return nil
	}
	p.nextToken() // {/if} を消費

	return node
}
