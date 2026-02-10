package expr

import "golox/token"

type Visitor[T any] interface {
	VisitBinary(binary *Binary[T]) T
	VisitUnary(unary *Unary[T]) T
	VisitGrouping(grouping *Grouping[T]) T
	VisitLiteral(literal *Literal[T]) T
	VisitVariable(variable *Variable[T]) T
}

type Expression[T any] interface {
	Accept(visitor Visitor[T]) T
}

type Binary[T any] struct {
	Left     Expression[T]
	Operator *token.Token
	Right    Expression[T]
}

func NewBinary[T any](left Expression[T], operator token.Token, right Expression[T]) Expression[T] {
	return &Binary[T]{
		Left:     left,
		Right:    right,
		Operator: &operator,
	}
}

func (b *Binary[T]) Accept(visitor Visitor[T]) T {
	return visitor.VisitBinary(b)
}

type Grouping[T any] struct {
	Expression Expression[T]
}

func NewGrouping[T any](expression Expression[T]) Expression[T] {
	return &Grouping[T]{Expression: expression}
}

func (g *Grouping[T]) Accept(visitor Visitor[T]) T {
	return visitor.VisitGrouping(g)
}

type Literal[T any] struct {
	Value token.Object
}

func NewLiteral[T any](value any) Expression[T] {
	return &Literal[T]{Value: value}
}

func (l *Literal[T]) Accept(visitor Visitor[T]) T {
	return visitor.VisitLiteral(l)
}

type Unary[T any] struct {
	Operator *token.Token
	Right    Expression[T]
}

func NewUnary[T any](operator token.Token, right Expression[T]) Expression[T] {
	return &Unary[T]{
		Operator: &operator,
		Right:    right,
	}
}

func (u *Unary[T]) Accept(visitor Visitor[T]) T {
	return visitor.VisitUnary(u)
}

type Variable[T any] struct {
	Name token.Token
}

func NewVariable[T any](name token.Token) Expression[T] {
	return &Variable[T]{
		Name: name,
	}
}

func (v *Variable[T]) Accept(visitor Visitor[T]) T {
	return visitor.VisitVariable(v)
}
