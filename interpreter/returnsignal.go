package interpreter

type returnSignal struct {
	value any
}

func NewReturnSignal(value any) returnSignal {
	return returnSignal{value: value}
}
