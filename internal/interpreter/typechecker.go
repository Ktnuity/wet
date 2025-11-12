package interpreter

import (
	"fmt"
	"strings"

	"github.com/ktnuity/wet/internal/errors"
	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

type StackPreview struct {
	stack			util.Stack[types.ValueType]
	refStack		util.Stack[types.ValueType]
	cause			*Instruction
	refCause		*Instruction
}

type TypeCheckData struct {
	inst			[]Instruction
	procs			map[string]Proc
	typeStack		util.Stack[types.ValueType]
	typeStackStack	util.Stack[StackPreview]
	eop				int
}

func pushStack(data *TypeCheckData, inst *Instruction) {
	data.typeStackStack.Push(StackPreview{
		stack: data.typeStack.Clone(),
		cause: inst,
	})
}

func popStack(data *TypeCheckData) *StackPreview {
	preview, ok := data.typeStackStack.Pop()
	if !ok {
		return nil
	}

	return &preview
}

func peekStack(data *TypeCheckData) *StackPreview {
	preview, ok := data.typeStackStack.PeekRef()
	if !ok {
		return nil
	}

	return preview
}

func compareStacks(current util.Stack[types.ValueType], other util.Stack[types.ValueType]) error {
	currentLen := len(current)
	poppedLen := len(other)

	if currentLen != poppedLen {
		return fmt.Errorf("failed to validate stack size: current(%d) != popped(%d)", currentLen, poppedLen)
	}

	for idx := range currentLen {
		if current[idx] != other[idx] {
			return fmt.Errorf("failed to validate stack: current[%d](%s) != popped[%d](%s) (len=%d). current=(%s) popped=(%s).", idx, current[idx].Format(), idx, other[idx].Format(), currentLen, typeDumpStack(current), typeDumpStack(other))
		}
	}

	return nil
}

type TypeCheckSubData struct {
	base			*TypeCheckData
	inst			*Instruction
	token			*types.Token
	tip				*int
}

type TypeResult struct {
	skip			bool
}

type typeCallback = func(*TypeCheckSubData) (*TypeResult, error)

func typeCheck(data *TypeCheckData) error {
	tip := 0

	handleClause := func(data *TypeCheckSubData, callback typeCallback) (bool, error) {
		result, err := callback(data)
		if err != nil {
			return false, fmt.Errorf("failed to type check. clause failed: %w", err)
		} else if result != nil {
			if !result.skip {
				tip++
			}
			return true, nil
		}

		return false, nil
	}

	typeCallbacks := []typeCallback{
		typeCheckLiteral,
		typeCheckPrint,
		typeCheckPrimary,
		typeCheckMemory,
		typeCheckArithmetic,
		typeCheckStack,
		typeCheckBranch,
		typeCheckLogical,
		typeCheckBitwise,
		typeCheckBoolean,
		typeCheckTools,
		typeCheckExtra,
	}

	for tip < data.eop {
		inst := &data.inst[tip]
		token := inst.Token

		if token == nil {
			return fmt.Errorf("failed to type check. token at index %d is nil.", tip)
		}

		typecheckv("Stack[%d]: %s\n", len(data.typeStack), typeDumpStack(data.typeStack))
		typecheckv("%d: %s\n", tip, token.Value)

		subData := TypeCheckSubData{
			data, inst, token, &tip,
		}

		var next bool
		var err error

		for _, callback := range typeCallbacks {
			next, err = handleClause(&subData, callback)
			if err != nil {
				return err
			} else if next {
				break
			}
		}

		p := errors.PrepareTypeCheck("call")
		proc, ok := data.procs[token.Value]
		if ok {
			ins := proc.Token.Extra.Proc.In
			outs := proc.Token.Extra.Proc.Out

			if data.typeStack.Len() < len(ins) {
				p.Stack(len(ins), data.typeStack.Len())
			}

			if err := compareStacks(data.typeStack[data.typeStack.Len()-len(ins):], ins); err != nil {
				p.CallLeadingError(token.Value, err)
			}

			data.typeStack = data.typeStack[:data.typeStack.Len()-len(ins)]
			data.typeStack = append(data.typeStack, outs...)

			tip++
		} else if !next {
			errors.BadTypeCheck(token.Value, "Unknown operator.")
		}
	}

	return nil
}

func typeDumpStack(stack util.Stack[types.ValueType]) string {
	typeStrings := make([]string, stack.Len())
	for idx := range stack.Len() {
		typeStrings[idx] = stack[idx].Format()
	}

	return strings.Join(typeStrings, " ")
}

func typeOk() (*TypeResult, error) {
	return &TypeResult{false}, nil
}

func typeSkip() (*TypeResult, error) {
	return &TypeResult{true}, nil
}

func typeNext() (*TypeResult, error) {
	return nil, nil
}

func typeBad(err error) (*TypeResult, error) {
	return nil, err
}

func typeCheckLiteral(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("", types.TokenTypeNumber) {
		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("", types.TokenTypeString) {
		d.base.typeStack.Push(types.ValueTypeString)
	} else if d.token.Equals("", types.TokenTypePath) {
		d.base.typeStack.Push(types.ValueTypePath)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckPrint(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals(".", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck(".", "Failed to print int.")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypeInt {
			p.ExpectType(t1.Format(), "int")
		}
	} else if d.token.Equals("puts", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("puts", "Failed to print string.")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypeString {
			p.ExpectType(t1.Format(), "string")
		}
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckPrimary(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("int", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("int")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypeString {
			p.ExpectType(t1.Format(), "string")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("string", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("string")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 == types.ValueTypeString {
			d.base.typeStack.Push(types.ValueTypeString)
		} else if t1 == types.ValueTypeInt {
			d.base.typeStack.Push(types.ValueTypeInt)
		} else {
			p.ExpectType(t1.Format(), "int", "string")
		}
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckMemory(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("store", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("store")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetNameValue(-1, "name")
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypeString {
			p.ExpectNameType("name", t1.Format(), "string")
		}

		if t2 != types.ValueTypeInt && t2 != types.ValueTypeString {
			p.ExpectNameType("value", t2.Format(), "int", "string")
		}
	} else if d.token.Equals("iload", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("iload")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetNameValue(-1, "name")
		}

		if t1 != types.ValueTypeString {
			p.ExpectNameType("name", t1.Format(), "string")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("sload", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("sload")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetNameValue(-1, "name")
		}

		if t1 != types.ValueTypeString {
			p.ExpectNameType("name", t1.Format(), "string")
		}

		d.base.typeStack.Push(types.ValueTypeString)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckArithmetic(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("+", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("+")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t2 == types.ValueTypeString {
			if t1 == types.ValueTypeString {
				d.base.typeStack.Push(types.ValueTypeString)
			} else if t1 == types.ValueTypeInt {
				d.base.typeStack.Push(types.ValueTypeString)
			} else {
				p.With("Second value is string.").ExpectNameType("first value", t1.Format(), "int", "string")
			}
		} else if t2 == types.ValueTypeInt {
			if t1 == types.ValueTypeInt {
				d.base.typeStack.Push(types.ValueTypeInt)
			} else {
				p.With("Second value is int.").ExpectNameType("first value", t1.Format(), "int")
			}
		} else {
			p.GetValue(1)
		}
	} else if d.token.Equals("-", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("-")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypeInt {
			p.ExpectNameType("first value", t1.Format(), "int")
		}

		if t2 != types.ValueTypeInt {
			p.ExpectNameType("second value", t2.Format(), "int")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("*", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("*")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypeInt {
			p.ExpectNameType("first value", t1.Format(), "int")
		}

		if t2 != types.ValueTypeInt {
			p.ExpectNameType("second value", t2.Format(), "int")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("/", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("/")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypeInt {
			p.ExpectNameType("first value", t1.Format(), "int")
		}

		if t2 != types.ValueTypeInt {
			p.ExpectNameType("second value", t2.Format(), "int")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("%", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("%")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypeInt {
			p.ExpectNameType("first value", t1.Format(), "int")
		}

		if t2 != types.ValueTypeInt {
			p.ExpectNameType("second value", t2.Format(), "int")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("++", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("++")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypeInt {
			p.ExpectType(t1.Format(), "int")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("--", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("--")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypeInt {
			p.ExpectType(t1.Format(), "int")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckStack(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("dup", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("dup")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypeInt && t1 != types.ValueTypeString && t1 != types.ValueTypePath {
			p.UnexpectedType(-1, t1.Format())
		}

		d.base.typeStack.Push(t1)
		d.base.typeStack.Push(t1)
	} else if d.token.Equals("drop", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("drop")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}
	} else if d.token.Equals("swap", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("swap")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypeInt && t1 != types.ValueTypeString && t1 != types.ValueTypePath {
			p.UnexpectedType(0, t1.Format())
		}

		if t2 != types.ValueTypeInt && t2 != types.ValueTypeString && t2 != types.ValueTypePath {
			p.UnexpectedType(1, t2.Format())
		}

		d.base.typeStack.Push(t1)
		d.base.typeStack.Push(t2)
	} else if d.token.Equals("over", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("over")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t2 := d.base.typeStack[len(d.base.typeStack)-2]

		if t2 != types.ValueTypeInt && t2 != types.ValueTypeString && t2 != types.ValueTypePath {
			p.UnexpectedType(1, t2.Format())
		}

		d.base.typeStack.Push(t2)
	} else if d.token.Equals("2dup", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("2dup")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1 := d.base.typeStack[len(d.base.typeStack)-2]
		t2 := d.base.typeStack[len(d.base.typeStack)-1]

		if t1 != types.ValueTypeInt && t1 != types.ValueTypeString && t1 != types.ValueTypePath {
			p.UnexpectedType(0, t1.Format())
		}

		if t2 != types.ValueTypeInt && t2 != types.ValueTypeString && t2 != types.ValueTypePath {
			p.UnexpectedType(1, t2.Format())
		}

		d.base.typeStack.Push(t1)
		d.base.typeStack.Push(t2)
	} else if d.token.Equals("2swap", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("2swap")
		if d.base.typeStack.Len() < 4 {
			p.Stack(4, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		t3, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(2)
		}

		t4, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(3)
		}

		if t1 != types.ValueTypeInt && t1 != types.ValueTypeString && t1 != types.ValueTypePath {
			p.UnexpectedType(0, t1.Format())
		}

		if t2 != types.ValueTypeInt && t2 != types.ValueTypeString && t2 != types.ValueTypePath {
			p.UnexpectedType(1, t2.Format())
		}

		if t3 != types.ValueTypeInt && t3 != types.ValueTypeString && t3 != types.ValueTypePath {
			p.UnexpectedType(2, t3.Format())
		}

		if t4 != types.ValueTypeInt && t4 != types.ValueTypeString && t4 != types.ValueTypePath {
			p.UnexpectedType(3, t4.Format())
		}

		d.base.typeStack.Push(t2)
		d.base.typeStack.Push(t1)
		d.base.typeStack.Push(t4)
		d.base.typeStack.Push(t3)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckBranch(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("if", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("if")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		pushStack(d.base, d.inst)
	} else if d.token.Equals("unless", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("unless")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		pushStack(d.base, d.inst)
	} else if d.token.Equals("else", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("else")
		preview := popStack(d.base)
		if preview == nil {
			p.Throw("No neighbor stack found.")
		}

		if preview.cause == nil || preview.cause.Token == nil {
			p.Throw("No connected keyword found.")
		}

		if !preview.cause.Token.Equals("if", types.TokenTypeKeyword) && !preview.cause.Token.Equals("unless", types.TokenTypeKeyword) {
			p.ExpectNameType("cause", preview.cause.Token.Value, "if", "unless")
		}

		d.base.typeStackStack.Push(StackPreview{
			stack: d.base.typeStack.Clone(),
			cause: d.inst,
		})

		d.base.typeStack = preview.stack.Clone()
	} else if d.token.Equals("do", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("do")
		preview := popStack(d.base)
		if preview == nil {
			p.Throw("No neighbor stack found.")
		}

		if preview.cause == nil || preview.cause.Token == nil {
			p.Throw("No connected keyword found.")
		}

		if !preview.cause.Token.Equals("while", types.TokenTypeKeyword) && !preview.cause.Token.Equals("until", types.TokenTypeKeyword) {
			p.ExpectNameType("cause", preview.cause.Token.Value, "while", "until")
		}

		pushStack(d.base, d.inst)
	} else if d.token.Equals("end", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("end")
		preview := popStack(d.base)
		if preview == nil {
			p.Throw("No neighbor stack found.")
		}

		if preview.cause == nil || preview.cause.Token == nil {
			p.Throw("No connected keyword found.")
		}

		if preview.cause.Token.Equals("if", types.TokenTypeKeyword) {
			if err := compareStacks(d.base.typeStack, preview.stack); err != nil {
				p.ConnectedTokenError("if", err)
			}
		} else if preview.cause.Token.Equals("unless", types.TokenTypeKeyword) {
			if err := compareStacks(d.base.typeStack, preview.stack); err != nil {
				p.ConnectedTokenError("unless", err)
			}
		} else if preview.cause.Token.Equals("else", types.TokenTypeKeyword) {
			if err := compareStacks(d.base.typeStack, preview.stack); err != nil {
				p.ConnectedTokenError("else", err)
			}
		} else if preview.cause.Token.Equals("do", types.TokenTypeKeyword) {
			if err := compareStacks(d.base.typeStack, preview.stack); err != nil {
				p.ConnectedTokenError("do", err)
			}
		} else if preview.cause.Token.Equals("proc", types.TokenTypeKeyword) {
			out := preview.cause.Token.Extra.Proc.Out
			if d.base.typeStack.Len() < len(out) {
				p.Stack(len(out), d.base.typeStack.Len())
			}

			if err := compareStacks(d.base.typeStack[d.base.typeStack.Len()-len(out):], out); err != nil {
				p.ConnectedTokenError("proc", err)
			}

			d.base.typeStack = d.base.typeStack[:d.base.typeStack.Len()-len(out)]

			typecheckv("Exit proc %s\n", preview.cause.Token.Extra.Proc.Name)
		} else {
			p.Throw(fmt.Sprintf("Connected keyword '%s' is not implemented.", preview.cause.Token.Value))
		}
	} else if d.token.Equals("while", types.TokenTypeKeyword) {
		pushStack(d.base, d.inst)
	} else if d.token.Equals("until", types.TokenTypeKeyword) {
		pushStack(d.base, d.inst)
	} else if d.token.Equals("proc", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("proc")
		pushStack(d.base, d.inst)

		for _, inType := range d.inst.Token.Extra.Proc.In {
			d.base.typeStack.Push(inType)
		}

		if d.inst.Next == -1 {
			p.Throw("Skip is undefined.")
		}

		typecheckv("Enter proc %s\n", d.token.Extra.Proc.Name)
	} else if d.token.Equals("ret", types.TokenTypeKeyword) {
		errors.BadTypeCheck("ret", "Not implemented.")
	} else if d.token.Equals("dret", types.TokenTypeKeyword) {
		errors.BadTypeCheck("dret", "Not implemented.")
	} else if d.token.Equals("iret", types.TokenTypeKeyword) {
		errors.BadTypeCheck("iret", "Not implemented.")
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckLogical(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("=", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("=")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		_, ok = d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("!=", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("!=")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		_, ok = d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("<", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("<")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 == types.ValueTypeString {
			if t2 != types.ValueTypeString {
				p.With(fmt.Sprintf("Cannot compare string and %s.", t2.Format())).Throw()
			}
		} else if t1 == types.ValueTypeInt {
			if t2 != types.ValueTypeInt {
				p.With(fmt.Sprintf("Cannot compare int and %s.", t2.Format())).Throw()
			}
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals(">", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck(">")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 == types.ValueTypeString {
			if t2 != types.ValueTypeString {
				p.With(fmt.Sprintf("Cannot compare string and %s.", t2.Format())).Throw()
			}
		} else if t1 == types.ValueTypeInt {
			if t2 != types.ValueTypeInt {
				p.With(fmt.Sprintf("Cannot compare int and %s.", t2.Format())).Throw()
			}
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("<=", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("<=")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 == types.ValueTypeString {
			if t2 != types.ValueTypeString {
				p.With(fmt.Sprintf("Cannot compare string and %s.", t2.Format())).Throw()
			}
		} else if t1 == types.ValueTypeInt {
			if t2 != types.ValueTypeInt {
				p.With(fmt.Sprintf("Cannot compare int and %s.", t2.Format())).Throw()
			}
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals(">=", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck(">=")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 == types.ValueTypeString {
			if t2 != types.ValueTypeString {
				p.With(fmt.Sprintf("Cannot compare string and %s.", t2.Format())).Throw()
			}
		} else if t1 == types.ValueTypeInt {
			if t2 != types.ValueTypeInt {
				p.With(fmt.Sprintf("Cannot compare int and %s.", t2.Format())).Throw()
			}
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("!", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("!")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckBitwise(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("~", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("~")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypeInt {
			p.ExpectType(t1.Format(), "int")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("&", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("&")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypeInt || t2 != types.ValueTypeInt {
			p.With(fmt.Sprintf("2 int values expected, got %s and %s.", t1.Format(), t2.Format())).Throw()
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("|", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("|")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypeInt || t2 != types.ValueTypeInt {
			p.With(fmt.Sprintf("2 int values expected, got %s and %s.", t1.Format(), t2.Format())).Throw()
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("^", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("^")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypeInt || t2 != types.ValueTypeInt {
			p.With(fmt.Sprintf("2 int values expected, got %s and %s.", t1.Format(), t2.Format())).Throw()
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckBoolean(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("&&", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("&&")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		_, ok = d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("||", types.TokenTypeSymbol) {
		p := errors.PrepareTypeCheck("||")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		_, ok = d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckTools(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("download", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("download")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypePath {
			p.ExpectNameType("first value", t1.Format(), "path")
		}

		if t2 != types.ValueTypeString {
			p.ExpectNameType("second value", t2.Format(), "string")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("readfile", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("readfile")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypePath {
			p.ExpectType(t1.Format(), "path")
		}

		d.base.typeStack.Push(types.ValueTypeString)
		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("copy", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("copy")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypePath {
			p.ExpectNameType("first value", t1.Format(), "path")
		}

		if t2 != types.ValueTypePath {
			p.ExpectNameType("second value", t2.Format(), "path")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("exist", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("exist")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypePath {
			p.ExpectType(t1.Format(), "path")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("touch", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("touch")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypePath {
			p.ExpectType(t1.Format(), "path")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("mkdir", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("mkdir")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypePath {
			p.ExpectType(t1.Format(), "path")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("rm", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("rm")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypePath {
			p.ExpectType(t1.Format(), "path")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("unzip", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("unzip")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypePath {
			p.ExpectNameType("first value", t1.Format(), "path")
		}

		if t2 != types.ValueTypePath {
			p.ExpectNameType("second value", t2.Format(), "path")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("lsf", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("lsf")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypePath {
			p.ExpectType(t1.Format(), "path")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("getf", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("getf")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypePath {
			p.ExpectNameType("first value", t1.Format(), "path")
		}

		if t2 != types.ValueTypeInt {
			p.ExpectNameType("second value", t2.Format(), "int")
		}

		d.base.typeStack.Push(types.ValueTypeString)
	} else if d.token.Equals("lsd", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("lsd")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypePath {
			p.ExpectType(t1.Format(), "path")
		}

		d.base.typeStack.Push(types.ValueTypeInt)
	} else if d.token.Equals("getd", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("getd")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypePath {
			p.ExpectNameType("first value", t1.Format(), "path")
		}

		if t2 != types.ValueTypeInt {
			p.ExpectNameType("second value", t2.Format(), "int")
		}

		d.base.typeStack.Push(types.ValueTypeString)
	} else if d.token.Equals("concat", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("concat")
		if d.base.typeStack.Len() < 2 {
			p.Stack(2, d.base.typeStack.Len())
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(0)
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(1)
		}

		if t1 != types.ValueTypeString {
			p.ExpectNameType("first value", t1.Format(), "string")
		}

		if t2 != types.ValueTypeString && t2 != types.ValueTypePath {
			p.ExpectNameType("second value", t2.Format(), "string", "path")
		}

		d.base.typeStack.Push(t2)
	} else if d.token.Equals("tostring", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("tostring")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypePath && t1 != types.ValueTypeInt && t1 != types.ValueTypeString {
			p.ExpectType(t1.Format(), "path", "int", "string")
		}

		d.base.typeStack.Push(types.ValueTypeString)
	} else if d.token.Equals("token", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("token")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypeString {
			p.ExpectType(t1.Format(), "string")
		}

		d.base.typeStack.Push(types.ValueTypePath)
	} else if d.token.Equals("absolute", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("absolute")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypeString {
			p.ExpectType(t1.Format(), "string")
		}

		d.base.typeStack.Push(types.ValueTypePath)
	} else if d.token.Equals("relative", types.TokenTypeKeyword) {
		p := errors.PrepareTypeCheck("relative")
		if d.base.typeStack.Len() < 1 {
			p.Empty()
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			p.GetValue(-1)
		}

		if t1 != types.ValueTypeString {
			p.ExpectType(t1.Format(), "string")
		}

		d.base.typeStack.Push(types.ValueTypePath)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckExtra(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("exit", types.TokenTypeKeyword) {
	} else {
		return typeNext()
	}

	return typeOk()
}
