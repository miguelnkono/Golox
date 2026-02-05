package main

import (
	"fmt"
	"golox/expr"
	"golox/token"
)

type runtimeError struct {
	tok     token.Token
	message string
}

func (r runtimeError) Error() string {
	return fmt.Sprintf("[line %d] RuntimeError: %s", r.tok.Line, r.message)
}

type Interpreter struct{}

func (i *Interpreter) Interpret(expression expr.Expression[any]) {
	defer func() {
		if r := recover(); r != nil {
			if rerr, ok := r.(runtimeError); ok {
				fmt.Println(rerr.Error())
			} else {
				// programming error, not a Lox runtime error
				panic(r)
			}
		}
	}()

	value := i.evaluate(expression)

	if value == nil {
		fmt.Println("nil")
		return
	}

	fmt.Println(value)
}

func (i *Interpreter) VisitBinary(b *expr.Binary[any]) any {
	left := i.evaluate(b.Left)
	right := i.evaluate(b.Right)

	switch b.Operator.TokenType {

	// Arithmetic
	case token.MINUS:
		i.checkNumberOperands(left, right, *b.Operator)
		return left.(float64) - right.(float64)

	case token.SLASH:
		i.checkNumberOperands(left, right, *b.Operator)
		if right.(float64) == 0 {
			panic(i.error(*b.Operator, "Division by zero"))
		}
		return left.(float64) / right.(float64)

	case token.STAR:
		i.checkNumberOperands(left, right, *b.Operator)
		return left.(float64) * right.(float64)

	case token.PLUS:
		switch l := left.(type) {
		case float64:
			if r, ok := right.(float64); ok {
				return l + r
			}
		case string:
			if r, ok := right.(string); ok {
				return l + r
			}
		}
		panic(i.error(
			*b.Operator,
			"Operands must be two numbers or two strings",
		))

	// Comparison
	case token.GREATER:
		i.checkNumberOperands(left, right, *b.Operator)
		return left.(float64) > right.(float64)

	case token.GREATER_EQUAL:
		i.checkNumberOperands(left, right, *b.Operator)
		return left.(float64) >= right.(float64)

	case token.LESS:
		i.checkNumberOperands(left, right, *b.Operator)
		return left.(float64) < right.(float64)

	case token.LESS_EQUAL:
		i.checkNumberOperands(left, right, *b.Operator)
		return left.(float64) <= right.(float64)

	// Equality
	case token.BANG_EQUAL:
		return !i.isEqual(left, right)

	case token.EQUAL_EQUAL:
		return i.isEqual(left, right)
	}

	// Unreachable
	return nil
}

func (i *Interpreter) VisitUnary(u *expr.Unary[any]) any {
	right := i.evaluate(u.Right)

	switch u.Operator.TokenType {
	case token.MINUS:
		i.checkNumberOperand(right, *u.Operator)
		return -right.(float64)

	case token.BANG:
		return !i.isTruthy(right)
	}

	return nil
}

func (i *Interpreter) VisitGrouping(g *expr.Grouping[any]) any {
	return i.evaluate(g.Expression)
}

func (i *Interpreter) VisitLiteral(l *expr.Literal[any]) any {
	return l.Value
}

func (i *Interpreter) evaluate(e expr.Expression[any]) any {
	return e.Accept(i)
}

func (i *Interpreter) isTruthy(value any) bool {
	if value == nil {
		return false
	}

	if b, ok := value.(bool); ok {
		return b
	}

	return true
}

func (i *Interpreter) isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a == b
}

func (i *Interpreter) checkNumberOperands(left, right any, op token.Token) {
	_, lok := left.(float64)
	_, rok := right.(float64)

	if lok && rok {
		return
	}

	panic(i.error(op, "Operands must be numbers"))
}

func (i *Interpreter) checkNumberOperand(value any, op token.Token) {
	if _, ok := value.(float64); ok {
		return
	}

	panic(i.error(op, "Operand must be a number"))
}

func (i *Interpreter) error(tok token.Token, format string, args ...any) runtimeError {
	return runtimeError{
		tok:     tok,
		message: fmt.Sprintf(format, args...),
	}
}
