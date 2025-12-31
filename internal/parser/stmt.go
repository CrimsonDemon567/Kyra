package parser

import (
	"fmt"
	"kyra/internal/lexer"
)

// ---------------------------
// IF / ELSE
// ---------------------------

func parseIfStmt(p *Parser) Stmt {
	p.expect(lexer.K_IF, "if")

	cond := p.parseExpression()

	// Two forms:
	// 1) if cond:
	// 2) if cond { ... }

	if p.match(lexer.COLON) {
		// Python-style block
		thenBlock := parseIndentedBlock(p)

		var elseBlock []Stmt
		if p.match(lexer.K_ELSE) {
			p.expect(lexer.COLON, "expected ':' after else")
			elseBlock = parseIndentedBlock(p)
		}

		return &IfStmt{
			Cond: cond,
			Then: thenBlock,
			Else: elseBlock,
		}
	}

	if p.match(lexer.LBRACE) {
		// One-liner block: if cond { stmt }
		thenBlock := parseBraceBlock(p)

		var elseBlock []Stmt
		if p.match(lexer.K_ELSE) {
			if p.match(lexer.LBRACE) {
				elseBlock = parseBraceBlock(p)
			} else {
				p.expect(lexer.COLON, "expected ':' or '{' after else")
				elseBlock = parseIndentedBlock(p)
			}
		}

		return &IfStmt{
			Cond: cond,
			Then: thenBlock,
			Else: elseBlock,
		}
	}

	panic("Expected ':' or '{' after if condition")
}

// ---------------------------
// WHILE
// ---------------------------

func parseWhileStmt(p *Parser) Stmt {
	p.expect(lexer.K_WHILE, "while")

	cond := p.parseExpression()

	if p.match(lexer.COLON) {
		body := parseIndentedBlock(p)
		return &WhileStmt{
			Cond: cond,
			Body: body,
		}
	}

	if p.match(lexer.LBRACE) {
		body := parseBraceBlock(p)
		return &WhileStmt{
			Cond: cond,
			Body: body,
		}
	}

	panic("Expected ':' or '{' after while condition")
}

// ---------------------------
// FOR
// ---------------------------

func parseForStmt(p *Parser) Stmt {
	p.expect(lexer.K_FOR, "for")

	// Syntax:
	// for i 10
	// for i range
	// for i expr

	varName := p.expect(lexer.IDENT, "loop variable").Lexeme

	limit := p.parseExpression()

	if p.match(lexer.COLON) {
		body := parseIndentedBlock(p)
		return &ForStmt{
			VarName: varName,
			Limit:   limit,
			Body:    body,
		}
	}

	if p.match(lexer.LBRACE) {
		body := parseBraceBlock(p)
		return &ForStmt{
			VarName: varName,
			Limit:   limit,
			Body:    body,
		}
	}

	panic("Expected ':' or '{' after for loop")
}

// ---------------------------
// BLOCK PARSING
// ---------------------------

func parseIndentedBlock(p *Parser) []Stmt {
	stmts := []Stmt{}

	// Expect INDENT
	if !p.match(lexer.INDENT) {
		panic("Expected INDENT after ':'")
	}

	for {
		tok := p.peek()

		if tok.Type == lexer.DEDENT {
			p.next()
			break
		}

		if tok.Type == lexer.NEWLINE {
			p.next()
			continue
		}

		stmt := p.parseTopLevel()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}

	return stmts
}

func parseBraceBlock(p *Parser) []Stmt {
	stmts := []Stmt{}

	for {
		tok := p.peek()

		if tok.Type == lexer.RBRACE {
			p.next()
			break
		}

		if tok.Type == lexer.NEWLINE {
			p.next()
			continue
		}

		stmt := p.parseTopLevel()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}

	return stmts
}

// ---------------------------
// RETURN / EXIT / PASS
// ---------------------------

func parseReturnStmt(p *Parser) Stmt {
	p.expect(lexer.K_RETURN, "return")
	value := p.parseExpression()
	return &ReturnStmt{Value: value}
}

func parseExitStmt(p *Parser) Stmt {
	p.expect(lexer.K_EXIT, "exit")
	return &ExitStmt{}
}

func parsePassStmt(p *Parser) Stmt {
	p.expect(lexer.K_PASS, "pass")
	return &PassStmt{}
}

// ---------------------------
// Expression statement
// ---------------------------

func parseExprStmt(p *Parser) Stmt {
	expr := p.parseExpression()
	return &ExprStmt{Expr: expr}
}
