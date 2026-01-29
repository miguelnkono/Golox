package tools

import (
	"fmt"
	"golox/expr"
	"strings"
)

type AstPrinter struct{}

func (a *AstPrinter) Print(expression expr.Expression[any]) string {
	return expression.Accept(a).(string)
}

func (a *AstPrinter) VisitBinary(binary *expr.Binary[any]) any {
	return a.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
}

func (a *AstPrinter) VisitUnary(unary *expr.Unary[any]) any {
	return a.parenthesize(unary.Operator.Lexeme, unary.Right)
}

func (a *AstPrinter) VisitGrouping(grouping *expr.Grouping[any]) any {
	return a.parenthesize("group", grouping.Expression)
}

func (a *AstPrinter) VisitLiteral(literal *expr.Literal[any]) any {
	if literal.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", literal.Value)
}

func (a *AstPrinter) parenthesize(name string, expressions ...expr.Expression[any]) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, exp := range expressions {
		builder.WriteString(" ")
		// We assert to string here because we know AstPrinter always returns strings
		result := exp.Accept(a)
    fmt.Fprintf(&builder, "%v", result)
	}

	builder.WriteString(")")

	return builder.String()
}
