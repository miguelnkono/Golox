package interpreter

import (
	"time"
)

type Clock struct {}

func NewClock() *Clock { return &Clock{}; }

func (c *Clock) Arity() int { return 0; }
func (c *Clock) Call(interpreter *Interpreter , arguments []any) any {
	return float64(time.Now().UnixMilli()) / 1000.0
}
func (c *Clock) String() string { return "<native fn>;"}
