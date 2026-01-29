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

type Interpreter struct {
}

func (i *Interpreter) Interpret() string {
	defer func() {
		if r := recover(); r != nil {
			if recovered, ok := r.(runtimeError); ok {
				panic(recovered)
			} else {
			}
		}
	}()
  
  // TODO: Finish this function today.
  return ""
}

func (i *Interpreter) VisitBinary(binary *expr.Binary[any]) any {
	left := i.evaluate(binary.Left)
	right := i.evaluate(binary.Right)

	switch binary.Operator.TokenType {

	// arithmetic operators
	case token.MINUS:
		i.checkNumberOperands(left, right, binary.Operator.String())
		return left.(float64) - right.(float64)
	case token.SLASH:
		i.checkNumberOperands(left, right, binary.Operator.String())
		return left.(float64) / right.(float64)
	case token.STAR:
		i.checkNumberOperands(left, right, binary.Operator.String())
		return left.(float64) * right.(float64)
	case token.PLUS:
		// TODO: need to add robust type checking here.
		switch left.(type) {
		case float64:
			return left.(float64) + right.(float64)
		case string:
			return fmt.Sprintf("%s+%s", left.(string), right.(string))
		default:
			panic(i.error(fmt.Sprintf("%s: Operands must be numbers or strings", binary.Operator.String())))
		}

		// comparison operators
	case token.GREATER:
		i.checkNumberOperands(left, right, binary.Operator)
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		i.checkNumberOperands(left, right, binary.Operator)
		return left.(float64) >= right.(float64)
	case token.LESS:
		i.checkNumberOperands(left, right, binary.Operator)
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		i.checkNumberOperands(left, right, binary.Operator)
		return left.(float64) <= right.(float64)

		// equality operators
	case token.BANG_EQUAL:
		return !i.isEqual(left, right)
	case token.EQUAL_EQUAL:
		return i.isEqual(left, right)
	}

	return nil
}

func (i *Interpreter) VisitUnary(unary *expr.Unary[any]) any {
	object := i.evaluate(unary)

	if unary.Operator.TokenType == token.MINUS {
		i.checkNumberOperand(object, unary.Operator.String())
		return -float64(object.(float64))
	}

	if unary.Operator.TokenType == token.BANG {
		return i.isTruthy(object)
	}

	// unreachable
	return nil
}

func (i *Interpreter) VisitGrouping(grouping *expr.Grouping[any]) any {
	return i.evaluate(grouping.Expression)
}

func (i *Interpreter) VisitLiteral(literal *expr.Literal[any]) any {
	return literal.Value
}

// helper functions
func (i *Interpreter) isEqual(left, right token.Object) bool {
	if left == nil && right == nil {
		return false
	}
	if left == nil {
		return false
	}
	if right == nil {
		return false
	}

	// TODO: will find maybe a better equality check
	return left == right
}

func (i *Interpreter) evaluate(expression expr.Expression[any]) any {
	return expression.Accept(i)
}

func (i *Interpreter) isTruthy(object any) bool {
	// false and null are falsey
	if object == nil {
		return false
	}

	switch object := object.(type) {
	case bool:
		return object
	}

	return true
}

// function to check on errors
func (i *Interpreter) checkNumberOperands(tok1, tok2 any, operand token.Object) {
	switch tok1.(type) {
	case float64:
		switch tok2.(type) {
		case float64:
			return
		}
	}

	// panic(i.error("Operand must be a number", fmt.Sprint("%v", tok.String())))
	panic(i.error(fmt.Sprintf("%s Operand must be a number", operand.(string))))
}

func (i *Interpreter) checkNumberOperand(tok any, operand token.Object) {
	switch tok.(type) {
	case float64:
		return
	}

	// panic(i.error("Operand must be a number", fmt.Sprint("%v", tok.String())))
	panic(i.error(fmt.Sprintf("%s Operand must be a numbers", operand.(string))))
}

func (r runtimeError) Error() string {
	return fmt.Sprintf("error: %s", r.message)
}

func (i *Interpreter) error(format string, any ...any) error {
	return runtimeError{message: fmt.Sprintf(format, any...)}
}
