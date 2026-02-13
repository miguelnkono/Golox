package main

import (
	"fmt"
	"golox/errors"
	"golox/token"
	"strconv"
)

var keywords = map[string]token.TokenType{
	"and":    token.AND,
	"or":     token.OR,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"true":   token.TRUE,
	"for":    token.FOR,
	"fun":    token.FUN,
	"nil":    token.NIL,
	"print":  token.PRINT,
	"super":  token.SUPER,
	"return": token.RETURN,
	"this":   token.THIS,
	"var":    token.VAR,
	"while":  token.WHILE,
	"if":     token.IF,
}

type Scanner struct {
	source string
	tokens []token.Token

	start, current, line int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []token.Token{},
		current: 0,
		start:   0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []token.Token {
	for !s.isAtEnd() {
		s.start = s.current

		s.scanToken()
	}

	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {

	// single character token
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)

		// one or two character(s) token
	case '=':
		s.addToken(If(s.match('='), token.EQUAL_EQUAL, token.EQUAL))
	case '<':
		s.addToken(If(s.match('='), token.LESS_EQUAL, token.LESS))
	case '>':
		s.addToken(If(s.match('='), token.GREATER_EQUAL, token.GREATER))
	case '!':
		s.addToken(If(s.match('='), token.BANG_EQUAL, token.BANG))
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			for !(s.peek() == '*' && s.peekNext() == '/') && !s.isAtEnd() {
				if s.peek() == '\n' {
					s.line++
				}

				s.advance()
			}

			if s.isAtEnd() {
				errors.Perror(s.line, "Unterminated comment.")
				return
			}
			s.advance()
			s.advance()
		} else {
			s.addToken(token.SLASH)
		}

		// character that don't have any sens
	case ' ':
	case '\t':
	case '\r':
		break
	case '\n':
		s.line++

		// long character (string)
	case '"':
		s.stringLiteral()

	default:
		if isDigit(char) {
			s.number()
		} else if isAlpha(char) {
			s.identifier()
		} else {
			errors.Perror(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) identifier() {

	for isAlphaNumberic(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, present := keywords[text]
	if present != true {
		tokenType = token.IDENTIFIER
	}
	s.addToken(tokenType)
}

func isAlphaNumberic(char byte) bool {
	return isDigit(char) || isAlpha(char)
}

func isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		char == '_'
}
func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	digit, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		errors.Perror(s.line, fmt.Sprintf("Can not convert %s to a floating point number", s.source[s.start:s.current]))
	}
	s.addTokenLiteral(token.NUMBER, digit)
}

func (s *Scanner) peekNext() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *Scanner) stringLiteral() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		errors.Perror(s.line, "Unterminated string")
		return
	}

	s.advance()

	text := s.source[s.start+1 : s.current-1]
	s.addTokenLiteral(token.STRING, text)
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}

	return s.source[s.current]
}

func If[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

func (s *Scanner) match(char byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != char {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) advance() byte {
	char := s.source[s.current]
	s.current++
	return char
}

func (s *Scanner) addToken(tokenType token.TokenType) {
	s.addTokenLiteral(tokenType, nil)
}

func (s *Scanner) addTokenLiteral(tokenType token.TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(tokenType, text, literal, s.line))
}
