package interpreter

type LoxCallable interface {
	Arity() int
	Call(*Interpreter, []any) any
}
