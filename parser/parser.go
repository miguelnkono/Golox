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
// statement      → exprStmt
//                | forStmt
//                | ifStmt
//                | printStmt
//                | returnStmt
//                | whileStmt
//                | block ;

// exprStmt       → expression ";" ;
// forStmt        → "for" "(" ( varDecl | exprStmt | ";" )
//                            expression? ";"
//                            expression? ")" statement ;
// ifStmt         → "if" "(" expression ")" statement
//                  ( "else" statement )? ;
// printStmt      → "print" expression ";" ;
// returnStmt     → "return" expression? ";" ;
// whileStmt      → "while" "(" expression ")" statement ;
// block          → "{" declaration* "}" ;
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
	if p.match(token.LEFT_BRACE) {
		return stmt.NewBlockStmt(p.block())
	}
	if p.match(token.IF) {
		return p.ifStatement()
	}
	if p.match(token.WHILE) {
		return p.whileStmt()
	}
	if p.match(token.FOR) {
		return p.forStmt()
	}

	return p.expressionStatement()
}

func (p *Parser) forStmt() stmt.Statement[any] {
	p.consume(token.LEFT_PAREN, "Expect '(' after the for keyword.")

	var initializer stmt.Statement[any]
	if p.match(token.VAR) {
		initializer = p.varDecleration()
	} else if p.match(token.SEMICOLON) {
		initializer = nil
	} else {
		initializer = p.expressionStatement()
	}

	var condition expr.Expression[any]
	if !p.check(token.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(token.SEMICOLON, "Expect ';' after the loop condition.")

	var increment expr.Expression[any]
	if !p.check(token.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(token.RIGHT_PAREN, "Expect ';' after the incrementation.")

	body := p.statement()

	if increment != nil {
		body = stmt.NewBlockStmt([]stmt.Statement[any]{body, stmt.NewExpressionStmt(increment)})
	}

	if condition == nil {
		condition = expr.NewLiteral[any](true)
	}
	body = stmt.NewWhileStmt(condition, body)

	if initializer != nil {
		body = stmt.NewBlockStmt([]stmt.Statement[any]{initializer, body})
	}

	return body
}

func (p *Parser) whileStmt() stmt.Statement[any] {
	p.consume(token.LEFT_PAREN, "Expect '(' after the while keyword.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after the condition.")
	body := p.statement()

	return stmt.NewWhileStmt(condition, body)
}

// parsing the if statement.
func (p *Parser) ifStatement() stmt.Statement[any] {
	p.consume(token.LEFT_PAREN, "Expect '(' after the if statement.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' at the end of the if statement.")
	thenBranch := p.statement()
	var elseBranch stmt.Statement[any]

	if p.match(token.ELSE) {
		elseBranch = p.statement()
	}

	return stmt.NewIfStmt(condition, thenBranch, elseBranch)
}

func (p *Parser) block() []stmt.Statement[any] {
	var stmts []stmt.Statement[any]

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}

	p.consume(token.RIGHT_BRACE, "Expected '}' after block.")
	return stmts
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
	// exp := p.equality()
	exp := p.or()

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

func (p *Parser) or() expr.Expression[any] {
	exp := p.and()
	for p.match(token.OR) {
		operator := p.previous()
		right := p.and()
		exp = expr.NewLogical(exp, operator, right)
	}

	return exp
}

func (p *Parser) and() expr.Expression[any] {
	exp := p.equality()

	for p.match(token.AND) {
		operator := p.previous()
		right := p.equality()
		exp = expr.NewLogical(exp, operator, right)
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
