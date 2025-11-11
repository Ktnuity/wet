package interpreter

import (
	"fmt"
	"strings"

	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

type ValueType uint8
const (
	ValueTypeNone ValueType = iota
	ValueTypeInt
	ValueTypeString
	ValueTypePath
)

func ValueTypeFormat(vt ValueType) string {
	switch vt {
	case ValueTypeNone:		return "None"
	case ValueTypeInt:		return "Int"
	case ValueTypeString:	return "String"
	case ValueTypePath:		return "Path"
	default:				return "Unknown"
	}
}

type StackPreview struct {
	stack			util.Stack[ValueType]
	refStack		util.Stack[ValueType]
	cause			*Instruction
	refCause		*Instruction
}

type TypeCheckData struct {
	inst			[]Instruction
	typeStack		util.Stack[ValueType]
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

func compareStacks(current util.Stack[ValueType], other util.Stack[ValueType]) error {
	currentLen := len(current)
	poppedLen := len(other)

	if currentLen != poppedLen {
		return fmt.Errorf("failed to validate stack size: current(%d) != popped(%d)", currentLen, poppedLen)
	}

	for idx := range currentLen {
		if current[idx] != other[idx] {
			return fmt.Errorf("failed to validate stack: current[%d](%s) != popped[%d](%s) (len=%d). current=(%s) popped=(%s).", idx, ValueTypeFormat(current[idx]), idx, ValueTypeFormat(other[idx]), currentLen, typeDumpStack(current), typeDumpStack(other))
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
	}

	for tip < data.eop {
		inst := &data.inst[tip]
		token := inst.Token

		if token == nil {
			return fmt.Errorf("failed to type check. token at index %d is nil.", tip)
		}

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
	}

	return nil
}

func typeDumpStack(stack util.Stack[ValueType]) string {
	typeStrings := make([]string, stack.Len())
	for idx := range stack.Len() {
		typeStrings[idx] = ValueTypeFormat(stack[idx])
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
		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("", types.TokenTypeString) {
		d.base.typeStack.Push(ValueTypeString)
	} else if d.token.Equals("", types.TokenTypePath) {
		d.base.typeStack.Push(ValueTypePath)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckPrint(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals(".", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf(". operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf(". operator failed. failed to get value."))
		}

		if t1 != ValueTypeInt {
			return typeBad(fmt.Errorf(". operator failed. expected int, got: %s", ValueTypeFormat(t1)))
		}
	} else if d.token.Equals("puts", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("puts operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("puts operator failed. failed to get value."))
		}

		if t1 != ValueTypeString {
			return typeBad(fmt.Errorf("puts operator failed. expected string, got: %s", ValueTypeFormat(t1)))
		}
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckPrimary(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("int", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("int operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("int operator failed. failed to get value."))
		}

		if t1 != ValueTypeString {
			return typeBad(fmt.Errorf("int operator failed. expected string, got: %s", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("string", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("string operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("string operator failed. failed to get value."))
		}

		if t1 == ValueTypeString {
			d.base.typeStack.Push(ValueTypeString)
		} else if t1 == ValueTypeInt {
			d.base.typeStack.Push(ValueTypeInt)
		} else {
			return typeBad(fmt.Errorf("string operator failed. expected int or string, got: %s", ValueTypeFormat(t1)))
		}
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckMemory(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("store", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("store operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("store operator failed. failed to get name."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("store operator failed. failed to get value."))
		}

		if t1 != ValueTypeString {
			return typeBad(fmt.Errorf("store operator failed. expected string name, got: %s", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypeInt {
			return typeBad(fmt.Errorf("store operator failed. expected int value, got: %s", ValueTypeFormat(t2)))
		}
	} else if d.token.Equals("load", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("load operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("load operator failed. failed to get name."))
		}

		if t1 != ValueTypeString {
			return typeBad(fmt.Errorf("load operator failed. expected name string, got: %s", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckArithmetic(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("+", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("+ operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("+ operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("+ operator failed. failed to get second value."))
		}

		if t2 == ValueTypeString {
			if t1 == ValueTypeString {
				d.base.typeStack.Push(ValueTypeString)
			} else if t1 == ValueTypeInt {
				d.base.typeStack.Push(ValueTypeString)
			} else {
				return typeBad(fmt.Errorf("+ operator failed. second value is string, expected first value int or string, got: %s.", ValueTypeFormat(t1)))
			}
		} else if t2 == ValueTypeInt {
			if t1 == ValueTypeInt {
				d.base.typeStack.Push(ValueTypeInt)
			} else {
				return typeBad(fmt.Errorf("+ operator failed. second value is int, expected first value int, got: %s.", ValueTypeFormat(t1)))
			}
		} else {
			return typeBad(fmt.Errorf("+ operator failed. failed to get second value."))
		}
	} else if d.token.Equals("-", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("- operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("- operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("- operator failed. failed to get second value."))
		}

		if t1 != ValueTypeInt {
			return typeBad(fmt.Errorf("- operator failed. expected int first value, got: %s", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypeInt {
			return typeBad(fmt.Errorf("- operator failed. expected int second value, got: %s", ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("*", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("* operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("* operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("* operator failed. failed to get second value."))
		}

		if t1 != ValueTypeInt {
			return typeBad(fmt.Errorf("* operator failed. expected int first value, got: %s", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypeInt {
			return typeBad(fmt.Errorf("* operator failed. expected int second value, got: %s", ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("/", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("/ operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("/ operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("/ operator failed. failed to get second value."))
		}

		if t1 != ValueTypeInt {
			return typeBad(fmt.Errorf("/ operator failed. expected int first value, got: %s", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypeInt {
			return typeBad(fmt.Errorf("/ operator failed. expected int second value, got: %s", ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("%", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("modulo operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("modulo operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("modulo operator failed. failed to get second value."))
		}

		if t1 != ValueTypeInt {
			return typeBad(fmt.Errorf("modulo operator failed. expected int first value, got: %s", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypeInt {
			return typeBad(fmt.Errorf("modulo operator failed. expected int second value, got: %s", ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("++", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("++ operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("++ operator failed. failed to get value."))
		}

		if t1 != ValueTypeInt {
			return typeBad(fmt.Errorf("++ operator failed. expected int value, got: %s", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("--", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("-- operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("-- operator failed. failed to get value."))
		}

		if t1 != ValueTypeInt {
			return typeBad(fmt.Errorf("-- operator failed. expected int value, got: %s", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckStack(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("dup", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("dup operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("dup operator failed. failed to get value."))
		}

		if t1 != ValueTypeInt && t1 != ValueTypeString && t1 != ValueTypePath {
			return typeBad(fmt.Errorf("dup operator failed. unexpected type: %s", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(t1)
		d.base.typeStack.Push(t1)
	} else if d.token.Equals("drop", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("dup operator failed. stack is empty."))
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("dup operator failed. failed to get value."))
		}
	} else if d.token.Equals("swap", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("swap operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("swap operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("swap operator failed. failed to get second value."))
		}

		if t1 != ValueTypeInt && t1 != ValueTypeString && t1 != ValueTypePath {
			return typeBad(fmt.Errorf("swap operator failed. unexpected first type: %s", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypeInt && t2 != ValueTypeString && t2 != ValueTypePath {
			return typeBad(fmt.Errorf("swap operator failed. unexpected second type: %s", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(t1)
		d.base.typeStack.Push(t2)
	} else if d.token.Equals("over", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("over operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t2 := d.base.typeStack[len(d.base.typeStack)-2]

		if t2 != ValueTypeInt && t2 != ValueTypeString && t2 != ValueTypePath {
			return typeBad(fmt.Errorf("over operator failed. unexpected second type: %s", ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(t2)
	} else if d.token.Equals("2dup", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("2dup operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1 := d.base.typeStack[len(d.base.typeStack)-2]
		t2 := d.base.typeStack[len(d.base.typeStack)-1]

		if t1 != ValueTypeInt && t1 != ValueTypeString && t1 != ValueTypePath {
			return typeBad(fmt.Errorf("2dup operator failed. unexpected first type: %s", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypeInt && t2 != ValueTypeString && t2 != ValueTypePath {
			return typeBad(fmt.Errorf("2dup operator failed. unexpected second type: %s", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(t1)
		d.base.typeStack.Push(t2)
	} else if d.token.Equals("2swap", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 4 {
			return typeBad(fmt.Errorf("swap operator failed. stack size is %d. 4 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("swap operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("swap operator failed. failed to get second value."))
		}

		t3, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("swap operator failed. failed to get third value."))
		}

		t4, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("swap operator failed. failed to get fourth value."))
		}

		if t1 != ValueTypeInt && t1 != ValueTypeString && t1 != ValueTypePath {
			return typeBad(fmt.Errorf("swap operator failed. unexpected first type: %s", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypeInt && t2 != ValueTypeString && t2 != ValueTypePath {
			return typeBad(fmt.Errorf("swap operator failed. unexpected second type: %s", ValueTypeFormat(t2)))
		}

		if t3 != ValueTypeInt && t3 != ValueTypeString && t3 != ValueTypePath {
			return typeBad(fmt.Errorf("swap operator failed. unexpected third type: %s", ValueTypeFormat(t3)))
		}

		if t4 != ValueTypeInt && t4 != ValueTypeString && t4 != ValueTypePath {
			return typeBad(fmt.Errorf("swap operator failed. unexpected fourth type: %s", ValueTypeFormat(t4)))
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
		pushStack(d.base, d.inst)
		//return typeBad(fmt.Errorf("if operator failed. not implemented."))
	} else if d.token.Equals("unless", types.TokenTypeKeyword) {
		pushStack(d.base, d.inst)
	} else if d.token.Equals("else", types.TokenTypeKeyword) {
		preview := popStack(d.base)
		if preview == nil {
			return typeBad(fmt.Errorf("else operator failed. no pushed stack."))
		}

		if preview.cause == nil || preview.cause.Token == nil {
			return typeBad(fmt.Errorf("else operator failed. no cause."))
		}

		if !preview.cause.Token.Equals("if", types.TokenTypeKeyword) && !preview.cause.Token.Equals("unless", types.TokenTypeKeyword) {
			return typeBad(fmt.Errorf("else operator failed. if or unless cause expected, got %s.", preview.cause.Token.Value))
		}

		d.base.typeStackStack.Push(StackPreview{
			stack: d.base.typeStack.Clone(),
			cause: d.inst,
		})

		d.base.typeStack = preview.stack.Clone()
	} else if d.token.Equals("do", types.TokenTypeKeyword) {
		preview := popStack(d.base)
		if preview == nil {
			return typeBad(fmt.Errorf("do operator failed. no pushed stack."))
		}

		if preview.cause == nil || preview.cause.Token == nil {
			return typeBad(fmt.Errorf("do operator failed. no cause."))
		}

		if !preview.cause.Token.Equals("while", types.TokenTypeKeyword) && !preview.cause.Token.Equals("until", types.TokenTypeKeyword) {
			return typeBad(fmt.Errorf("do operator failed. while or until cause expected, got %s.", preview.cause.Token.Value))
		}

		pushStack(d.base, d.inst)
		//return typeBad(fmt.Errorf("do operator failed. not implemented."))
	} else if d.token.Equals("end", types.TokenTypeKeyword) {
		preview := popStack(d.base)
		if preview == nil {
			return typeBad(fmt.Errorf("end operator failed. no pushed stack."))
		}

		if preview.cause == nil || preview.cause.Token == nil {
			return typeBad(fmt.Errorf("end operator failed. no cause."))
		}

		if preview.cause.Token.Equals("if", types.TokenTypeKeyword) {
			if err := compareStacks(d.base.typeStack, preview.stack); err != nil {
				return typeBad(fmt.Errorf("end operator failed for if cause: %w", err))
			}
		} else if preview.cause.Token.Equals("unless", types.TokenTypeKeyword) {
			if err := compareStacks(d.base.typeStack, preview.stack); err != nil {
				return typeBad(fmt.Errorf("end operator failed for unless cause: %w", err))
			}
		} else if preview.cause.Token.Equals("else", types.TokenTypeKeyword) {
			if err := compareStacks(d.base.typeStack, preview.stack); err != nil {
				return typeBad(fmt.Errorf("end operator failed for else cause: %w", err))
			}
		} else if preview.cause.Token.Equals("do", types.TokenTypeKeyword) {
			if err := compareStacks(d.base.typeStack, preview.stack); err != nil {
				return typeBad(fmt.Errorf("end operator failed for do cause: %w", err))
			}
		} else {
			return typeBad(fmt.Errorf("end operator failed. %s cause not implemented.", preview.cause.Token.Value))
		}
	} else if d.token.Equals("while", types.TokenTypeKeyword) {
		pushStack(d.base, d.inst)
		//return typeBad(fmt.Errorf("while operator failed. not implemented."))
	} else if d.token.Equals("until", types.TokenTypeKeyword) {
		pushStack(d.base, d.inst)
		//return typeBad(fmt.Errorf("until operator failed. not implemented."))
	} else if d.token.Equals("proc", types.TokenTypeKeyword) {
		if d.inst.Next == -1 {
			return typeBad(fmt.Errorf("proc operator failed. skip is undefined.\n"))
		}

		(*d.tip) = int(d.inst.Next)
		return typeSkip()
	} else if d.token.Equals("ret", types.TokenTypeKeyword) {
		return typeBad(fmt.Errorf("ret operator failed. not implemented."))
	} else if d.token.Equals("dret", types.TokenTypeKeyword) {
		return typeBad(fmt.Errorf("dret operator failed. not implemented."))
	} else if d.token.Equals("iret", types.TokenTypeKeyword) {
		return typeBad(fmt.Errorf("iret operator failed. not implemented."))
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckLogical(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("=", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("= operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("= operator failed. failed to get first value."))
		}

		_, ok = d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("= operator failed. failed to get second value."))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("!=", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("!= operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("!= operator failed. failed to get first value."))
		}

		_, ok = d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("!= operator failed. failed to get second value."))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("<", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("< operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("< operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("< operator failed. failed to get second value."))
		}

		if t1 == ValueTypeString {
			if t2 != ValueTypeString {
				return typeBad(fmt.Errorf("< operator failed. cannot compare string and %s.", ValueTypeFormat(t2)))
			}
		} else if t1 == ValueTypeInt {
			if t2 != ValueTypeInt {
				return typeBad(fmt.Errorf("< operator failed. cannot compare int and %s.", ValueTypeFormat(t2)))
			}
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals(">", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("> operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("> operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("> operator failed. failed to get second value."))
		}

		if t1 == ValueTypeString {
			if t2 != ValueTypeString {
				return typeBad(fmt.Errorf("> operator failed. cannot compare string and %s.", ValueTypeFormat(t2)))
			}
		} else if t1 == ValueTypeInt {
			if t2 != ValueTypeInt {
				return typeBad(fmt.Errorf("> operator failed. cannot compare int and %s.", ValueTypeFormat(t2)))
			}
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("<=", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("<= operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("<= operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("<= operator failed. failed to get second value."))
		}

		if t1 == ValueTypeString {
			if t2 != ValueTypeString {
				return typeBad(fmt.Errorf("<= operator failed. cannot compare string and %s.", ValueTypeFormat(t2)))
			}
		} else if t1 == ValueTypeInt {
			if t2 != ValueTypeInt {
				return typeBad(fmt.Errorf("<= operator failed. cannot compare int and %s.", ValueTypeFormat(t2)))
			}
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals(">=", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf(">= operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf(">= operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf(">= operator failed. failed to get second value."))
		}

		if t1 == ValueTypeString {
			if t2 != ValueTypeString {
				return typeBad(fmt.Errorf(">= operator failed. cannot compare string and %s.", ValueTypeFormat(t2)))
			}
		} else if t1 == ValueTypeInt {
			if t2 != ValueTypeInt {
				return typeBad(fmt.Errorf(">= operator failed. cannot compare int and %s.", ValueTypeFormat(t2)))
			}
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("!", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("! operator failed. stack is empty."))
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("! operator failed. failed to get first value."))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckBitwise(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("~", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("~ operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("~ operator failed. failed to get first value."))
		}

		if t1 != ValueTypeInt {
			return typeBad(fmt.Errorf("~ operator failed. expected int, got %s.", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("&", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("& operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("& operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("& operator failed. failed to get second value."))
		}

		if t1 != ValueTypeInt || t2 != ValueTypeInt {
			return typeBad(fmt.Errorf("& operator failed. 2 int values expected, got %s and %s.", ValueTypeFormat(t1), ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("|", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("| operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("| operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("| operator failed. failed to get second value."))
		}

		if t1 != ValueTypeInt || t2 != ValueTypeInt {
			return typeBad(fmt.Errorf("| operator failed. 2 int values expected, got %s and %s.", ValueTypeFormat(t1), ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("^", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("^ operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("^ operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("^ operator failed. failed to get second value."))
		}

		if t1 != ValueTypeInt || t2 != ValueTypeInt {
			return typeBad(fmt.Errorf("^ operator failed. 2 int values expected, got %s and %s.", ValueTypeFormat(t1), ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckBoolean(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("&&", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("& operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("& operator failed. failed to get first value."))
		}

		_, ok = d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("& operator failed. failed to get second value."))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("||", types.TokenTypeSymbol) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("& operator failed. stack size is %d. 2 is requied.", d.base.typeStack.Len()))
		}

		_, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("& operator failed. failed to get first value."))
		}

		_, ok = d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("& operator failed. failed to get second value."))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else {
		return typeNext()
	}

	return typeOk()
}

func typeCheckTools(d *TypeCheckSubData) (*TypeResult, error) {
	if d.token.Equals("download", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("download operator failed. stack size is %d. 2 is required.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("download operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("download operator failed. failed to get second value."))
		}

		if t1 != ValueTypePath {
			return typeBad(fmt.Errorf("download operator failed. expected path as first value, got %s.", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypeString {
			return typeBad(fmt.Errorf("download operator failed. expected string as second value, got %s.", ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("readfile", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("readfile operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("readfile operator failed. failed to get first value."))
		}

		if t1 != ValueTypePath {
			return typeBad(fmt.Errorf("readfile operator failed. expected path, got %s.", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypeString)
		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("copy", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("copy operator failed. stack size is %d. 2 is required.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("copy operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("copy operator failed. failed to get second value."))
		}

		if t1 != ValueTypePath {
			return typeBad(fmt.Errorf("copy operator failed. expected path as first value, got %s.", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypePath {
			return typeBad(fmt.Errorf("copy operator failed. expected path as second value, got %s.", ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(ValueTypeInt)
		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("exist", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("exist operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("exist operator failed. failed to get first value."))
		}

		if t1 != ValueTypePath {
			return typeBad(fmt.Errorf("exist operator failed. expected path, got %s.", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("touch", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("touch operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("touch operator failed. failed to get first value."))
		}

		if t1 != ValueTypePath {
			return typeBad(fmt.Errorf("touch operator failed. expected path, got %s.", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("mkdir", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("mkdir operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("mkdir operator failed. failed to get first value."))
		}

		if t1 != ValueTypePath {
			return typeBad(fmt.Errorf("mkdir operator failed. expected path, got %s.", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("rm", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("rm operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("rm operator failed. failed to get first value."))
		}

		if t1 != ValueTypePath {
			return typeBad(fmt.Errorf("rm operator failed. expected path, got %s.", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("unzip", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("unzip operator failed. stack size is %d. 2 is required.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("unzip operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("unzip operator failed. failed to get second value."))
		}

		if t1 != ValueTypePath {
			return typeBad(fmt.Errorf("unzip operator failed. expected path as first value, got %s.", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypePath {
			return typeBad(fmt.Errorf("unzip operator failed. expected path as second value, got %s.", ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(ValueTypeInt)
		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("lsf", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("lsf operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("lsf operator failed. failed to get first value."))
		}

		if t1 != ValueTypePath {
			return typeBad(fmt.Errorf("lsf operator failed. expected path, got %s.", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("getf", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("getf operator failed. stack size is %d. 2 is required.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("getf operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("getf operator failed. failed to get second value."))
		}

		if t1 != ValueTypePath {
			return typeBad(fmt.Errorf("getf operator failed. expected path as first value, got %s.", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypeInt {
			return typeBad(fmt.Errorf("getf operator failed. expected int as second value, got %s.", ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(ValueTypeString)
	} else if d.token.Equals("lsd", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("lsd operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("lsd operator failed. failed to get first value."))
		}

		if t1 != ValueTypePath {
			return typeBad(fmt.Errorf("lsd operator failed. expected path, got %s.", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypeInt)
	} else if d.token.Equals("getd", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("getd operator failed. stack size is %d. 2 is required.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("getd operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("getd operator failed. failed to get second value."))
		}

		if t1 != ValueTypePath {
			return typeBad(fmt.Errorf("getd operator failed. expected path as first value, got %s.", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypeInt {
			return typeBad(fmt.Errorf("getd operator failed. expected int as second value, got %s.", ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(ValueTypeString)
	} else if d.token.Equals("concat", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 2 {
			return typeBad(fmt.Errorf("concat operator failed. stack size is %d. 2 is required.", d.base.typeStack.Len()))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("concat operator failed. failed to get first value."))
		}

		t2, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("concat operator failed. failed to get second value."))
		}

		if t1 != ValueTypeString {
			return typeBad(fmt.Errorf("concat operator failed. expected string as first value, got %s.", ValueTypeFormat(t1)))
		}

		if t2 != ValueTypeString && t2 != ValueTypePath {
			return typeBad(fmt.Errorf("concat operator failed. expected string or path as second value, got %s.", ValueTypeFormat(t2)))
		}

		d.base.typeStack.Push(t2)
	} else if d.token.Equals("tostring", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("tostring operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("tostring operator failed. failed to get first value."))
		}

		if t1 != ValueTypePath && t1 != ValueTypeInt && t1 != ValueTypeString {
			return typeBad(fmt.Errorf("tostring operator failed. expected path, int, or string, got %s.", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypeString)
	} else if d.token.Equals("token", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("token operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("token operator failed. failed to get first value."))
		}

		if t1 != ValueTypeString {
			return typeBad(fmt.Errorf("token operator failed. expected string, got %s.", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypePath)
	} else if d.token.Equals("absolute", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("absolute operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("absolute operator failed. failed to get first value."))
		}

		if t1 != ValueTypeString {
			return typeBad(fmt.Errorf("absolute operator failed. expected string, got %s.", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypePath)
	} else if d.token.Equals("relative", types.TokenTypeKeyword) {
		if d.base.typeStack.Len() < 1 {
			return typeBad(fmt.Errorf("relative operator failed. stack is empty."))
		}

		t1, ok := d.base.typeStack.Pop()
		if !ok {
			return typeBad(fmt.Errorf("relative operator failed. failed to get first value."))
		}

		if t1 != ValueTypeString {
			return typeBad(fmt.Errorf("relative operator failed. expected string, got %s.", ValueTypeFormat(t1)))
		}

		d.base.typeStack.Push(ValueTypePath)
	} else {
		return typeNext()
	}

	return typeOk()
}

/*
func typeCheck(d *TypeCheckSubData) (*TypeResult, error) {
	} else {
		return typeNext()
	}

	return typeOk()
}
*/
