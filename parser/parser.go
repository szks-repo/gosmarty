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

const (
	_ int = iota
	LOWEST
	OR
	AND
	COMPARISON
)

var precedences = map[token.TokenType]int{
	token.OR:    OR,
	token.AND:   AND,
	token.EQ:    COMPARISON,
	token.NOTEQ: COMPARISON,
	token.LT:    COMPARISON,
	token.LTE:   COMPARISON,
	token.GT:    COMPARISON,
	token.GTE:   COMPARISON,
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
	case token.FOREACH:
		return p.parseForeachTag()
	case token.FOREACHELSE:
		p.errors = append(p.errors, "unexpected {foreachelse} without matching {foreach}")
		p.consumeUntil(token.RDELIM)
		return nil
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
	p.nextToken() // 'if' を消費 -> curTokenは 条件式の先頭

	cond := p.parseExpression(LOWEST)
	if cond == nil {
		return nil
	}
	node.Condition = cond

	// 2. {if ...} の閉じ '}' を消費
	if !p.curTokenIs(token.RDELIM) {
		p.errors = append(p.errors, fmt.Sprintf("expected RDELIM after if condition, got %s", p.curToken.Type))
		return nil
	}
	// '}' を消費
	p.nextToken()

	// 3. Consequence (then節) のパース
	node.Consequence = p.parseBlockUntil(token.ELSE, token.ELSEIF, token.ENDIF)

	// 4. ElseIf ブランチのパース
	for p.curTokenIs(token.LDELIM) && p.peekTokenIs(token.ELSEIF) {
		// '{' を消費
		p.nextToken()
		elseifNode := &ast.ElseIfNode{Token: p.curToken}
		// 'elseif' を消費
		p.nextToken()

		cond := p.parseExpression(LOWEST)
		if cond == nil {
			return nil
		}

		elseifNode.Condition = cond

		if !p.curTokenIs(token.RDELIM) {
			p.errors = append(p.errors, "expected RDELIM for elseif tag")
			return nil
		}
		// '}' を消費
		p.nextToken()

		elseifNode.Consequence = p.parseBlockUntil(token.ELSEIF, token.ELSE, token.ENDIF)
		node.ElseIfs = append(node.ElseIfs, elseifNode)
	}

	// 5. Alternative (else節) のパース (存在すれば)
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

	// 6. 終了タグ {/if} を消費
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

func (p *Parser) parseForeachTag() *ast.ForeachNode {
	// curTokenは '{'
	p.nextToken() // '{' を消費 -> curTokenは 'foreach'
	node := &ast.ForeachNode{Token: p.curToken}
	p.nextToken() // 'foreach' を消費 -> curTokenは最初の属性名

	for !p.curTokenIs(token.RDELIM) && !p.curTokenIs(token.EOF) {
		if !p.curTokenIs(token.IDENT) {
			p.errors = append(p.errors, fmt.Sprintf("expected attribute name for foreach, got %s", p.curToken.Type))
			return nil
		}

		attrName := p.curToken.Literal
		p.nextToken()

		if !p.curTokenIs(token.ASSIGN) {
			p.errors = append(p.errors, fmt.Sprintf("expected '=' after foreach attribute %q", attrName))
			return nil
		}
		p.nextToken()

		switch attrName {
		case "from":
			expr := p.parseExpression(LOWEST)
			if expr == nil {
				return nil
			}
			node.Source = expr
		case "item":
			name, ok := p.parseForeachVariableName()
			if !ok {
				return nil
			}
			if node.Item != "" {
				p.errors = append(p.errors, "duplicate item attribute in foreach")
				return nil
			}
			node.Item = name
		case "key":
			name, ok := p.parseForeachVariableName()
			if !ok {
				return nil
			}
			if node.Key != "" {
				p.errors = append(p.errors, "duplicate key attribute in foreach")
				return nil
			}
			node.Key = name
		case "name":
			name, ok := p.parseForeachName()
			if !ok {
				return nil
			}
			node.Name = name
		default:
			p.errors = append(p.errors, fmt.Sprintf("unsupported foreach attribute: %s", attrName))
			return nil
		}
	}

	if !p.curTokenIs(token.RDELIM) {
		p.errors = append(p.errors, "expected RDELIM to close foreach tag")
		return nil
	}
	// '}' を消費
	p.nextToken()

	if node.Source == nil {
		p.errors = append(p.errors, "foreach requires from attribute")
		return nil
	}
	if node.Item == "" {
		p.errors = append(p.errors, "foreach requires item attribute")
		return nil
	}

	node.Body = p.parseBlockUntil(token.FOREACHELSE, token.ENDFOREACH)

	if p.curTokenIs(token.LDELIM) && p.peekTokenIs(token.FOREACHELSE) {
		// '{' を消費
		p.nextToken()
		// 'foreachelse' を消費
		p.nextToken()

		if !p.curTokenIs(token.RDELIM) {
			p.errors = append(p.errors, "expected RDELIM for foreachelse tag")
			return nil
		}
		// '}' を消費
		p.nextToken()

		node.Alternative = p.parseBlockUntil(token.ENDFOREACH)
	}

	if !(p.curTokenIs(token.LDELIM) && p.peekTokenIs(token.ENDFOREACH)) {
		p.errors = append(p.errors, "expected {/foreach} tag")
		return nil
	}
	// '{' を消費
	p.nextToken()
	// '/foreach' を消費
	p.nextToken()

	if !p.curTokenIs(token.RDELIM) {
		p.errors = append(p.errors, "expected RDELIM for /foreach tag")
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

func (p *Parser) parseForeachVariableName() (string, bool) {
	switch p.curToken.Type {
	case token.DOLLAR:
		p.nextToken()
		if !p.curTokenIs(token.IDENT) {
			p.errors = append(p.errors, fmt.Sprintf("expected IDENT after '$' in foreach attribute, got %s", p.curToken.Type))
			return "", false
		}
		name := p.curToken.Literal
		p.nextToken()
		return name, true
	case token.IDENT:
		name := p.curToken.Literal
		p.nextToken()
		return name, true
	default:
		p.errors = append(p.errors, fmt.Sprintf("expected variable name in foreach attribute, got %s", p.curToken.Type))
		return "", false
	}
}

func (p *Parser) parseForeachName() (string, bool) {
	switch p.curToken.Type {
	case token.IDENT, token.STRING:
		name := p.curToken.Literal
		p.nextToken()
		return name, true
	default:
		p.errors = append(p.errors, fmt.Sprintf("expected identifier or string for foreach name attribute, got %s", p.curToken.Type))
		return "", false
	}
}

func (p *Parser) consumeUntil(t token.TokenType) {
	// 現在のトークンから指定のトークンまで読み飛ばす
	p.nextToken()
	for !p.curTokenIs(token.EOF) && !p.curTokenIs(t) {
		p.nextToken()
	}
	if p.curTokenIs(t) {
		p.nextToken()
	}
}

func (p *Parser) parseVariableTagWithPipeline() ast.Node {
	// curTokenは '{'
	p.nextToken() // '{' を消費 -> curTokenは '$'

	// 最初に左辺の式をパース
	left := p.parseExpression(LOWEST)
	if left == nil {
		return nil
	}

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

func (p *Parser) parsePrimaryExpr_backup() ast.Node {
	var left ast.Node

	switch p.curToken.Type {
	case token.DOLLAR:
		p.nextToken() // '$' を消費
		if !isIdentLike(p.curToken.Type) {
			p.errors = append(p.errors, fmt.Sprintf("expected IDENT-like token, got %s", p.curToken.Type))
			return nil
		}
		left = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken() // 識別子を消費
	case token.NUMBER:
		left = p.parseNumberLiteral()
	default:
		// 他のプライマリー式（文字列リテラルなど）もここに追加できる
		p.errors = append(p.errors, fmt.Sprintf("unexpected token for primary expression: %s", p.curToken.Type))
		return nil
	}

	for p.curTokenIs(token.DOT) {
		dotToken := p.curToken
		// '.' を消費
		p.nextToken()

		if !isIdentLike(p.curToken.Type) {
			p.errors = append(p.errors, fmt.Sprintf("expected IDENT-like token after '.', got %s", p.curToken.Type))
			return nil
		}

		right := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}

		left = &ast.FieldAccess{
			Token: dotToken,
			Left:  left,
			Right: right,
		}
		// プロパティ識別子を消費
		p.nextToken()
	}

	return left
}

func (p *Parser) parsePrimaryExpr() ast.Node {
	var left ast.Node

	switch p.curToken.Type {
	case token.DOLLAR:
		p.nextToken() // '$' を消費
		if !p.curTokenIs(token.IDENT) {
			p.errors = append(p.errors, fmt.Sprintf("expected IDENT, got %s", p.curToken.Type))
			return nil
		}
		left = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken() // 識別子を消費
	case token.NUMBER:
		left = p.parseNumberLiteral()
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token for primary expression: %s", p.curToken.Type))
		return nil
	}

	// .prop や [index] のような後置演算子をパースするループ
	for {
		switch p.curToken.Type {
		case token.DOT:
			dotToken := p.curToken
			p.nextToken() // '.' を消費

			if !isIdentLike(p.curToken.Type) {
				p.errors = append(p.errors, fmt.Sprintf("expected IDENT-like token after '.', got %s", p.curToken.Type))
				return nil
			}

			right := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			left = &ast.FieldAccess{
				Token: dotToken,
				Left:  left,
				Right: right,
			}
			p.nextToken() // プロパティ識別子を消費

		case token.LBRACKET: // ここを修正します
			bracketToken := p.curToken
			p.nextToken() // '[' を消費

			index := p.parseExpression(LOWEST) // インデックス内の式 ('0'など) をパース
			if index == nil {
				return nil
			}

			// 現在のトークンが ']' であることを確認
			if !p.curTokenIs(token.RBRACKET) {
				msg := fmt.Sprintf("expected token to be ], got %s instead", p.curToken.Type)
				p.errors = append(p.errors, msg)
				return nil
			}

			left = &ast.IndexExpression{
				Token: bracketToken,
				Left:  left,
				Index: index,
			}

			p.nextToken() // ']' を消費して次に進む

		default:
			// 後置演算子がなければループを抜ける
			return left
		}
	}

	return left
}

func (p *Parser) parseExpression(precedence int) ast.Node {
	left := p.parsePrimaryExpr()
	if left == nil {
		return nil
	}

	for !p.curTokenIs(token.RDELIM) && !p.curTokenIs(token.PIPE) && !p.curTokenIs(token.RBRACKET) && !p.curTokenIs(token.EOF) {
		curPrec := p.curPrecedence()
		if precedence >= curPrec {
			break
		}
		tokType := p.curToken.Type
		if _, ok := precedences[tokType]; !ok {
			break
		}

		tok := p.curToken
		p.nextToken()

		right := p.parseExpression(curPrec)
		if right == nil {
			return nil
		}

		left = &ast.InfixExpression{
			Token:    tok,
			Left:     left,
			Operator: tok.Literal,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
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

func isIdentLike(t token.TokenType) bool {
	switch t {
	case token.IDENT,
		token.FOREACH,
		token.FOREACHELSE,
		token.IF,
		token.ELSE,
		token.ELSEIF,
		token.ENDIF,
		token.ENDFOREACH,
		token.AND,
		token.OR,
		token.LITERAL,
		token.ENDLITERAL:
		return true
	default:
		return false
	}
}
