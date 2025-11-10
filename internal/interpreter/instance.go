package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

type Interpreter struct {
	stack		util.Stack[StackValue]
	program		[]Instruction
	memory		map[string]StackValue
	ip			int
	eop			int
}

func CreateNew(tokens []types.Token) (*Interpreter, error) {
	var stack util.Stack[StackValue] = util.Stack[StackValue]{}

	program, err := ProcessTokens(tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to create interpreter: %w", err)
	}

	return &Interpreter{
		stack: stack,
		program: program,
		memory: make(map[string]StackValue),
		ip: 0,
		eop: int(len(program)),
	}, nil
}
