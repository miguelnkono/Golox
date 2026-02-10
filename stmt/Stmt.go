package stmt

import "golox/expr"

// visitor
type Visitor[T any] interface {
	VisitExpressionStmt(ep *ExpressionStmt[T])
	VisitPrintStmt(ep *PrintStmt[T])
}

// statement class
type Statement[T any] interface {
	Accept(visitor Visitor[T])
}

// expression statement class
type ExpressionStmt[T any] struct {
	Expr expr.Expression[T]
}

func NewExpressionStmt[T any](e expr.Expression[T]) *ExpressionStmt[T] {
	return &ExpressionStmt[T]{
		Expr: e,
	}
}
func (es *ExpressionStmt[T]) Accept(visitor Visitor[T]) {
	visitor.VisitExpressionStmt(es)
}

// print statement class
type PrintStmt[T any] struct {
	Expr expr.Expression[T]
}

func NewPrintStmt[T any](e expr.Expression[T]) *PrintStmt[T] {
	return &PrintStmt[T]{
		Expr: e,
	}
}
func (es *PrintStmt[T]) Accept(visitor Visitor[T]) {
	visitor.VisitPrintStmt(es)
}
