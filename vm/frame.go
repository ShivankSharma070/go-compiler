package vm

import (
	"github.com/ShivankSharma070/go-compiler/code"
	"github.com/ShivankSharma070/go-compiler/object"
)

type Frame struct {
	c *object.Closure
	ip          int
	basePointer int
}

func NewFrame(c *object.Closure, basePointer int) *Frame {
	return &Frame{
		c:          c,
		ip:          -1,
		basePointer: basePointer,
	}
}
func (f *Frame) Instructions() code.Instructions {
	return f.c.Fn.Instructions
}
