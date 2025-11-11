package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

type Interpreter struct {
	stack		util.Stack[types.DataType]
	program		[]Instruction
	procs		map[string]Proc
	memory		map[string]types.DataType
	ip			int
	eop			int
	calls		util.Stack[int]
}

func CreateNew(tokens []types.Token) (*Interpreter, error) {
	var stack util.Stack[types.DataType] = util.Stack[types.DataType]{}

	result, err := ProcessTokens(tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to create interpreter: %w", err)
	}

	inst := &Interpreter{
		stack: stack,
		program: result.inst,
		procs: result.proc,
		memory: make(map[string]types.DataType),
		ip: 0,
		eop: int(len(result.inst)),
		calls: util.Stack[int]{},
	}

	tcd := &TypeCheckData{
		inst: inst.program,
		procs: result.proc,
		typeStack: util.Stack[types.ValueType]{},
		typeStackStack: util.Stack[StackPreview]{},
		eop: len(inst.program),
	}

	err = typeCheck(tcd)
	if err != nil {
		return nil, fmt.Errorf("failed to create interpreter: %w", err)
	}

	return inst, nil
}
