package parser

import (
	"fmt"
	"golox/expr"
	"golox/token"
	"strconv"
)

/*

expression    -> equality
equality      -> comparasion (("==" | "!=") comparasion)* ;
comparasion   -> term ((">" | ">=" | "<" | "<=") comparasion)* ;
term          -> factor (("+" | "-") factor)* ;
factor        -> unary (("/" | "*") unary)* ;
unary         -> ("!" | "-") unary
                  | primary ;
primary       -> NUMBER | STRING | "true" | "false" | "nil"
                  | "(" expression ")" ;

*/

type parseError struct {
	message string
	token   token.Token
}

func (p parseError) Error() string {
	if p.token.Lexeme != "" {
		return fmt.Sprintf("Error at '%s' on line %d: %s", p.token.Lexeme, p.token.Line+1, p.message)
	}
	return p.message
}

type Parser struct {
	Tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() (expr.Expression[any], error) {
	var lastErr error

	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(parseError); ok {
				lastErr = e
				p.synchronize()
			} else {
				panic(r)
			}
		}
	}()

	result := p.expression()

	// TODO: This block is need till we'll add statements
	if p.check(token.SEMICOLON) {
		p.advance()
	}

	if !p.isAtEnd() {
		return nil, p.error(p.peek(), "Expect end of expression.")
	}

	return result, lastErr
}

// expression rule
func (p *Parser) expression() expr.Expression[any] {
	return p.equality()
}

// equality rule
func (p *Parser) equality() expr.Expression[any] {
	expression := p.comparasion()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparasion()
		expression = expr.NewBinary(expression, operator, right)
	}

	return expression
}

// comparasion rule
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
	}

	if p.match(token.STRING, token.NUMBER) {
		tok := p.previous()
		return expr.NewLiteral[any](tok.Literal)
	}

	if p.match(token.LEFT_PAREN) {
		expression := p.expression()
		p.consume(token.RIGHT_PAREN, "Expected ')' at the end of a grouping expression.")
		return expr.NewGrouping(expression)
	}

	panic(p.error(p.peek(), fmt.Sprintf("expecting an expression, but found '%s'", p.peek().Lexeme)))
}

// utilities functions
func (p *Parser) consume(tokenType token.TokenType, message string) token.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	panic(parseError{message: fmt.Sprintf("Error at %s on line %d: %s", strconv.Quote(p.peek().Lexeme), p.peek().Line+1, message)})
}

func (p *Parser) error(tok token.Token, message string) error {
	return parseError{token: tok, message: message}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == token.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) match(tokens ...token.TokenType) bool {

	for _, tokenType := range tokens {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	// found := slices.ContainsFunc(tokens, func(t token.TokenType) bool {
	// 	return p.check(t)
	// })

	// if found {
	// 	p.advance()
	// 	return true
	// }

	return false
}

func (p *Parser) check(token token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	tok := p.peek()
	return tok.TokenType == token
}

func (p *Parser) isAtEnd() bool {
	tok := p.peek()
	if tok.TokenType == token.EOF {
		return true
	}
	return false
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
