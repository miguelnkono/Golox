package interpreter

import (
	"fmt"
	"golox/stmt"
)

type LoxFunction struct {
	declaration stmt.Function[any]
}

func NewLoxFunction(declaration stmt.Function[any]) *LoxFunction {
	return & LoxFunction{
		declaration: declaration,
	}
}

func (l *LoxFunction) Call(interpreter *Interpreter, arguments []any) (result any) {
	environment := NewEnclosedEnvironment(interpreter.globals)

	for i := 0; i < len(l.declaration.Parameters); i++ {
		environment.Define(l.declaration.Parameters[i].Lexeme, arguments[i])
	}

	defer func() {
		if r := recover(); r != nil {
			if ret, ok := r.(returnSignal); ok {
				result = ret.value
			} else {
				panic(r) 
			}
		}
	}()
	interpreter.executeBlock(l.declaration.Body, environment)

	return nil
}

func (l *LoxFunction) Arity() int {
	return len(l.declaration.Parameters)
}

func (l *LoxFunction) String() string { 
	return fmt.Sprintf("<fn %s >\n", l.declaration.Name.Lexeme); 
}
