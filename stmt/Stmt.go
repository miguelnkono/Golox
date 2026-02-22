package stmt

import (
	"golox/expr"
	"golox/token"
)

// visitor
type Visitor[T any] interface {
	VisitExpressionStmt(ep *ExpressionStmt[T])
	VisitPrintStmt(ep *PrintStmt[T])
	VisitVarStmt(v *VarStmt[T])
	VisitBlockStmt(b *BlockStmt[T])
	VisitIfStmt(i *IfStmt[T])
	VisitWhileStmt(w *WhileStmt[T])
}

// statement class
type Statement[T any] interface {
	Accept(visitor Visitor[T])
}

// if statement
type WhileStmt[T any] struct {
	Condition expr.Expression[T]
	Body      Statement[T]
}

func NewWhileStmt[T any](condition expr.Expression[T], body Statement[T]) *WhileStmt[T] {
	return &WhileStmt[T]{
		Condition: condition,
		Body:      body,
	}
}

func (w *WhileStmt[T]) Accept(visitor Visitor[T]) {
	visitor.VisitWhileStmt(w)
}

type IfStmt[T any] struct {
	Condition  expr.Expression[T]
	ThenBranch Statement[T]
	ElseBranch Statement[T]
}

func NewIfStmt[T any](condition expr.Expression[T], thenBranch, elseBranch Statement[T]) *IfStmt[T] {
	return &IfStmt[T]{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (i *IfStmt[T]) Accept(visitor Visitor[T]) {
	visitor.VisitIfStmt(i)
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

// Var statement class
type VarStmt[T any] struct {
	Name        token.Token
	Initializer expr.Expression[T]
}

func NewVarStmt[T any](name token.Token, initializer expr.Expression[T]) *VarStmt[T] {
	return &VarStmt[T]{
		Name:        name,
		Initializer: initializer,
	}
}
func (v *VarStmt[T]) Accept(visitor Visitor[T]) {
	visitor.VisitVarStmt(v)
}

// block statement
type BlockStmt[T any] struct {
	Stmts []Statement[T]
}

func NewBlockStmt[T any](stmts []Statement[T]) *BlockStmt[T] {
	return &BlockStmt[T]{
		Stmts: stmts,
	}
}

func (b *BlockStmt[T]) Accept(visitor Visitor[T]) {
	visitor.VisitBlockStmt(b)
}
