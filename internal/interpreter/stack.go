package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) pop() (types.DataType, error) {
	var value types.DataType
	value, ok := ip.stack.Pop()
	if !ok {
		return value, fmt.Errorf("failed to ipop. stack is empty.")
	}

	ip.runtimev("new stack size %d\n", ip.stack.Len())

	return value, nil
}

func (ip *Interpreter) peek() (types.DataType, error) {
	return ip.peekOffset(0)
}

func (ip *Interpreter) peekOffset(offset int) (types.DataType, error) {
	var value types.DataType
	idx := ip.stack.Len() - 1 - offset
	if idx < 0 {
		return value, fmt.Errorf("failed to ipeek. idx(%d) < 0", idx)
	}
	if idx >= ip.stack.Len() {
		return value, fmt.Errorf("failed to ipek. idx(%d) >= len(%d)", idx, ip.stack.Len())
	}

	return ip.stack[ip.stack.Len() - 1 - offset], nil
}

func (ip *Interpreter) push(value types.DataType) error {
	if !value.IsAny() {
		return fmt.Errorf("failed to push stack. value given, but unknown type.")
	}

	ip.stack.Push(value)
	ip.runtimev("new stack size %d\n", ip.stack.Len())

	return nil
}

func (ip *Interpreter) ipush(item int) error {
	return ip.push(types.TypeInt(item))
}

func (ip *Interpreter) spush(item string) error {
	return ip.push(types.TypeString(item))
}

func (ip *Interpreter) ppush(item string) error {
	return ip.push(types.TypePath(item))
}
