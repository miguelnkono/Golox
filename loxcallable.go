package main

import (
	"golox/token"
)

type LoxCallable interface {
	Arity() int
	Call(*Interpreter, []token.Object) token.Object
}
