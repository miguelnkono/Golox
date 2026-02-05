package parser

import (
	"fmt"
	"golox/expr"
	"golox/token"
	"slices"
)

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

func (p *Parser) Parse() (expression expr.Expression[any], err error) {
	defer func() {
		if r := recover(); r != nil {
			if pe, ok := r.(parseError); ok {
				expression = nil
				err = pe
			} else {
				panic(r) // Programming error, not parse error
			}
		}
	}()

	expression = p.expression()

	if p.match(token.SEMICOLON) {
		// expression statement (temporary)
	}

	if !p.isAtEnd() {
		return nil, p.error(p.peek(), "Expected end of expression")
	}

	return expression, nil
}

// Expression parsing methods remain the same...
func (p *Parser) expression() expr.Expression[any] {
	return p.equality()
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
