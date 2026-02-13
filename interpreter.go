package main

import (
	"fmt"
	"golox/errors"
	"golox/expr"
	"golox/stmt"
	"golox/token"
)

type runtimeError struct {
	tok     token.Token
	message string
}

func (r runtimeError) Error() string {
	return fmt.Sprintf("[line %d] RuntimeError: %s", r.tok.Line, r.message)
}

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(),
	}
}

func (i *Interpreter) Interpret(stmts []stmt.Statement[any]) {
	defer func() {
		if r := recover(); r != nil {
			if rerr, ok := r.(runtimeError); ok {
				fmt.Println(rerr.Error())
				errors.HadRuntimeError = true
			} else {
				// This is a programming error, not a Lox runtime error
				panic(r)
			}
		}
	}()

	for _, st := range stmts {

		i.execute(st)
	}
}

func (i *Interpreter) execute(st stmt.Statement[any]) {
	st.Accept(i)
}

// statement visitor
func (i *Interpreter) VisitExpressionStmt(ep *stmt.ExpressionStmt[any]) {
	i.evaluate(ep.Expr)
}

func (i *Interpreter) VisitPrintStmt(ep *stmt.PrintStmt[any]) {
	value := i.evaluate(ep.Expr)
	fmt.Println(i.stringify(value))
}

func (i *Interpreter) VisitVarStmt(v *stmt.VarStmt[any]) {
	var value token.Object
	if v.Initializer != nil {
		value = i.evaluate(v.Initializer)
	}
	i.environment.Define(v.Name.Lexeme, value)
}

// expression visitor
func (i *Interpreter) VisitVariable(v *expr.Variable[any]) any {

	value, err := i.environment.Get(v.Name.Lexeme)
	if err != nil {
		panic(i.error(v.Name, "%s", err.Error()))
	}

	return value
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
		// Handle number + number
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r
			}
		}
		// Handle string + string
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r
			}
		}
		panic(i.error(*b.Operator, "Operands must be two numbers or two strings"))

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

func (i *Interpreter) VisitAssignment(a *expr.Assignment[any]) any {
	value := i.evaluate(a.Exp)
	i.environment.Assign(a.Tok.Lexeme, value)

	return value
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

func (i *Interpreter) stringify(value any) string {
	if value == nil {
		return "nil"
	}

	if num, ok := value.(float64); ok {
		if num == float64(int64(num)) {
			return fmt.Sprintf("%d", int64(num))
		}
		return fmt.Sprintf("%g", num)
	}

	return fmt.Sprintf("%v", value)
}
