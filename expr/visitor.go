package expr

type Visitor[T any] interface {
	VisitBinary(binary *Binary[T]) T
	VisitUnary(unary *Unary[T]) T
	VisitGrouping(grouping *Grouping[T]) T
	VisitLiteral(literal *Literal[T]) T
}
