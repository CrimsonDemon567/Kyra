package parser

import (
	"fmt"
	"kyra/internal/lexer"
)

// ---------------------------
// def name(args): block
// ---------------------------

func parseDefFunc(p *Parser) Stmt {
	p.expect(lexer.K_DEF, "def")

	name := p.expect(lexer.IDENT, "function name").Lexeme

	args := parseFuncArgs(p)

	var returnType *TypeNode
	if p.match(lexer.ARROW) {
		returnType = parseType(p)
	}

	p.expect(lexer.COLON, "expected ':' after def signature")

	body := parseIndentedBlock(p)

	return &FuncDef{
		Name:       name,
		Args:       args,
		ReturnType: returnType,
		Body:       body,
	}
}

// ---------------------------
// func name(args) { block }
// func name(args) = expr
// ---------------------------

func parseFuncForms(p *Parser) Stmt {
	p.expect(lexer.K_FUNC, "func")

	name := p.expect(lexer.IDENT, "function name").Lexeme

	args := parseFuncArgs(p)

	var returnType *TypeNode
	if p.match(lexer.ARROW) {
		returnType = parseType(p)
	}

	// One-liner: func name(args) = expr
	if p.match(lexer.ASSIGN) {
		expr := p.parseExpression()
		return &FuncOneLiner{
			Name:       name,
			Args:       args,
			ReturnType: returnType,
			Expr:       expr,
		}
	}

	// Block: func name(args) { ... }
	if p.match(lexer.LBRACE) {
		body := parseBraceBlock(p)
		return &FuncExprDef{
			Name:       name,
			Args:       args,
			ReturnType: returnType,
			Body:       body,
		}
	}

	panic("Expected '=' or '{' after func signature")
}

// ---------------------------
// Function arguments
// ---------------------------

func parseFuncArgs(p *Parser) []FuncArg {
	p.expect(lexer.LPAREN, "expected '(' after function name")

	args := []FuncArg{}

	if p.peek().Type != lexer.RPAREN {
		for {
			name := p.expect(lexer.IDENT, "argument name").Lexeme

			var t *TypeNode
			if p.match(lexer.COLON) {
				t = parseType(p)
			}

			args = append(args, FuncArg{
				Name: name,
				Type: t,
			})

			if !p.match(lexer.COMMA) {
				break
			}
		}
	}

	p.expect(lexer.RPAREN, "expected ')' after arguments")

	return args
}

// ---------------------------
// Type parsing
// ---------------------------

func parseType(p *Parser) *TypeNode {
	tok := p.peek()

	switch tok.Type {
	case lexer.K_I32, lexer.K_I64,
		lexer.K_F32, lexer.K_F64,
		lexer.K_BOOL, lexer.K_STRING,
		lexer.K_VOID:

		p.next()
		return &TypeNode{Name: tok.Lexeme}

	default:
		panic(fmt.Sprintf("Unexpected type: %s (%s)", tok.Type, tok.Lexeme))
	}
}
