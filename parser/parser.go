package parser

import (
	"fmt"
	"golox/expr"
	"golox/stmt"
	"golox/token"
	"slices"
)

// program        -> declaration* EOF ;
// declaration    -> varDecl
//                | statement ;
// varDecl        -> "var" IDENTIFY ( "=" expression )? ";" ;
// statement      -> exprStmt
//                | printStmt ;
// exprStmt       -> expression ";" ;
// printStmt      -> "print" expression ";" ;

// expression     → assignment ;

// assignment     → ( call "." )? IDENTIFIER "=" assignment
//                | logic_or ;
//
// logic_or       → logic_and ( "or" logic_and )* ;
// logic_and      → equality ( "and" equality )* ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
//
// unary          → ( "!" | "-" ) unary | call ;
// call           → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
// primary        → "true" | "false" | "nil" | "this"
//                | NUMBER | STRING | IDENTIFIER | "(" expression ")"
//                | "super" "." IDENTIFIER ;

type parseError struct {
	message string
	token   token.Token
}

func (p parseError) Error() string {
	if p.token.TokenType == token.EOF {
		return fmt.Sprintf("[line %d] ParseError at end: %s", p.token.Line, p.message)
	}
	return fmt.Sprintf("[line %d] ParseError at '%s': %s", p.token.Line, p.token.Lexeme, p.message)
}

type Parser struct {
	Tokens  []token.Token
	current int
	errors  []error // Collect multiple errors
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		current: 0,
		errors:  []error{},
	}
}

func (p *Parser) Parse() (stmts []stmt.Statement[any], err []error) {
	var statements []stmt.Statement[any]

	for !p.isAtEnd() {
		decl := p.declaration()
		if decl != nil {
			statements = append(statements, decl)
		}
	}

	// Return any collected errors
	if len(p.errors) > 0 {
		return statements, p.errors
	}

	return statements, nil
}

func (p *Parser) declaration() stmt.Statement[any] {
	defer func() {
		if r := recover(); r != nil {
			if pe, ok := r.(parseError); ok {
				p.errors = append(p.errors, pe)
				p.synchronize()
			} else {
				panic(r)
			}
		}
	}()

	if p.match(token.VAR) {
		return p.varDecleration()
	}
	return p.statement()
}

func (p *Parser) varDecleration() stmt.Statement[any] {
	name := p.consume(token.IDENTIFIER, "Expected a variable name.")

	var initializer expr.Expression[any] = nil
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}

	p.consume(token.SEMICOLON, "Expected ';' after variable declaration!")
	return stmt.NewVarStmt(name, initializer)
}

func (p *Parser) statement() stmt.Statement[any] {
	if p.match(token.PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() stmt.Statement[any] {
	exp := p.expression()
	p.consume(token.SEMICOLON, "Expected ';' after value.")
	return stmt.NewPrintStmt(exp)
}

func (p *Parser) expressionStatement() stmt.Statement[any] {
	exp := p.expression()
	p.consume(token.SEMICOLON, "Expected ';' after value.")
	return stmt.NewExpressionStmt(exp)
}

func (p *Parser) expression() expr.Expression[any] {
	return p.assignment()
}

func (p *Parser) assignment() expr.Expression[any] {
	exp := p.equality()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if variable, ok := exp.(*expr.Variable[any]); ok {
			name := variable.Name
			return expr.NewAssignment(name, value)
		}

		panic(p.error(equals, "Invalid assignment target!"))
	}

	return exp
}

func (p *Parser) equality() expr.Expression[any] {
	expression := p.comparasion()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparasion()
		expression = expr.NewBinary(expression, operator, right)
	}

	return expression
}

func (p *Parser) comparasion() expr.Expression[any] {
	expression := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expression = expr.NewBinary(expression, operator, right)
	}

	return expression
}

func (p *Parser) term() expr.Expression[any] {
	expression := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expression = expr.NewBinary(expression, operator, right)
	}

	return expression
}

func (p *Parser) factor() expr.Expression[any] {
	expression := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expression = expr.NewBinary(expression, operator, right)
	}

	return expression
}

func (p *Parser) unary() expr.Expression[any] {
	for p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return expr.NewUnary(operator, right)
	}

	return p.primary()
}

func (p *Parser) primary() expr.Expression[any] {
	switch {
	case p.match(token.FALSE):
		return expr.NewLiteral[any](false)
	case p.match(token.TRUE):
		return expr.NewLiteral[any](true)
	case p.match(token.NIL):
		return expr.NewLiteral[any](nil)
	case p.match(token.STRING, token.NUMBER):
		tok := p.previous()
		return expr.NewLiteral[any](tok.Literal)
	case p.match(token.LEFT_PAREN):
		expression := p.expression()
		p.consume(token.RIGHT_PAREN, "Expected ')' after expression")
		return expr.NewGrouping(expression)
	case p.match(token.IDENTIFIER):
		return expr.NewVariable[any](p.previous())
	}

	panic(p.error(p.peek(), "Expected expression"))
}

// Utility functions
func (p *Parser) consume(tokenType token.TokenType, message string) token.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	panic(p.error(p.peek(), message))
}

func (p *Parser) error(tok token.Token, message string) parseError {
	return parseError{token: tok, message: message}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == token.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case token.CLASS, token.FUN, token.VAR, token.FOR,
			token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) match(tokens ...token.TokenType) bool {
	found := slices.ContainsFunc(tokens, func(t token.TokenType) bool {
		return p.check(t)
	})

	if found {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) check(tokType token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokType
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == token.EOF
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) peek() token.Token {
	return p.Tokens[p.current]
}

func (p *Parser) previous() token.Token {
	if p.current <= 0 {
		return p.Tokens[0]
	}
	return p.Tokens[p.current-1]
}
