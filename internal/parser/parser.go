package parser

import (
	"fmt"
	"kyra/internal/lexer"
)

// Parser is the main structure that consumes tokens and produces an AST.
type Parser struct {
	tokens []lexer.Token
	pos    int
}

// New creates a new parser instance.
func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// ---------------------------
// Token helpers
// ---------------------------

func (p *Parser) peek() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peekNext() lexer.Token {
	if p.pos+1 >= len(p.tokens) {
		return lexer.Token{Type: lexer.EOF}
	}
	return p.tokens[p.pos+1]
}

func (p *Parser) next() lexer.Token {
	tok := p.peek()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return tok
}

func (p *Parser) match(types ...lexer.TokenType) bool {
	tok := p.peek()
	for _, t := range types {
		if tok.Type == t {
			p.next()
			return true
		}
	}
	return false
}

func (p *Parser) expect(t lexer.TokenType, msg string) lexer.Token {
	tok := p.peek()
	if tok.Type != t {
		panic(fmt.Sprintf("Parse error at line %d: expected %s but got %s (%s)",
			tok.Line, t, tok.Type, msg))
	}
	return p.next()
}

// ---------------------------
// Entry point
// ---------------------------

// Parse parses a full Kyra module.
func (p *Parser) Parse() *AST {
	ast := &AST{
		ModuleName: "main",
		Imports:    []*UseStmt{},
		TopLevel:   []Stmt{},
	}

	// Parse imports first
	for p.match(lexer.NEWLINE) {
	}

	for p.peek().Type == lexer.K_USE {
		use := p.parseUse()
		ast.Imports = append(ast.Imports, use)
		for p.match(lexer.NEWLINE) {
		}
	}

	// Parse top-level statements
	for p.peek().Type != lexer.EOF {
		stmt := p.parseTopLevel()
		if stmt != nil {
			ast.TopLevel = append(ast.TopLevel, stmt)
		}
		for p.match(lexer.NEWLINE) {
		}
	}

	return ast
}

// ---------------------------
// Top-level parsing
// ---------------------------

func (p *Parser) parseTopLevel() Stmt {
	tok := p.peek()

	switch tok.Type {

	case lexer.K_DEF:
		return p.parseDefFunction()

	case lexer.K_FUNC:
		return p.parseFuncVariants()

	case lexer.K_LET:
		return p.parseLet()

	case lexer.K_IF:
		return p.parseIf()

	case lexer.K_WHILE:
		return p.parseWhile()

	case lexer.K_FOR:
		return p.parseFor()

	case lexer.K_RETURN:
		return p.parseReturn()

	case lexer.K_EXIT:
		p.next()
		return &ExitStmt{}

	case lexer.K_PASS:
		p.next()
		return &PassStmt{}

	case lexer.IDENT, lexer.NUMBER, lexer.STRING, lexer.LPAREN:
		return p.parseExprStmt()

	case lexer.INDENT:
		p.next()
		return p.parseTopLevel()

	case lexer.DEDENT:
		p.next()
		return nil

	case lexer.NEWLINE:
		p.next()
		return nil

	default:
		panic(fmt.Sprintf("Unexpected token at top-level: %s (%s)", tok.Type, tok.Lexeme))
	}
}

// ---------------------------
// Use statement
// ---------------------------

func (p *Parser) parseUse() *UseStmt {
	p.expect(lexer.K_USE, "use statement")

	parts := []string{}
	isStdlib := false

	// Detect stdlib prefix: sdt/
	if p.peek().Type == lexer.IDENT && p.peek().Lexeme == "sdt" {
		isStdlib = true
		p.next()
		p.expect(lexer.SLASH, "expected '/' after sdt")
	}

	// Parse module path: a/b/c
	for {
		tok := p.expect(lexer.IDENT, "module name")
		parts = append(parts, tok.Lexeme)

		if !p.match(lexer.SLASH) {
			break
		}
	}

	return &UseStmt{
		Path:     parts,
		IsStdlib: isStdlib,
	}
}

// ---------------------------
// Let statement
// ---------------------------

func (p *Parser) parseLet() Stmt {
	p.expect(lexer.K_LET, "let statement")

	name := p.expect(lexer.IDENT, "variable name").Lexeme
	p.expect(lexer.ASSIGN, "assignment '='")

	expr := p.parseExpression()

	return &LetStmt{
		Name: name,
		Expr: expr,
	}
}

// ---------------------------
// Expression statement
// ---------------------------

func (p *Parser) parseExprStmt() Stmt {
	expr := p.parseExpression()
	return &ExprStmt{Expr: expr}
}

// ---------------------------
// Return
// ---------------------------

func (p *Parser) parseReturn() Stmt {
	p.expect(lexer.K_RETURN, "return")
	value := p.parseExpression()
	return &ReturnStmt{Value: value}
}

// ---------------------------
// Delegation to other files
// ---------------------------

func (p *Parser) parseIf() Stmt {
	return parseIfStmt(p)
}

func (p *Parser) parseWhile() Stmt {
	return parseWhileStmt(p)
}

func (p *Parser) parseFor() Stmt {
	return parseForStmt(p)
}

func (p *Parser) parseDefFunction() Stmt {
	return parseDefFunc(p)
}

func (p *Parser) parseFuncVariants() Stmt {
	return parseFuncForms(p)
}

func (p *Parser) parseExpression() Expr {
	return parseExpr(p, 0)
}
