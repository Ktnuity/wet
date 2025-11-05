package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

type StackValue = util.Either[string, int]

type Interpreter struct {
	stack		util.Stack[StackValue]
	program		[]Instruction
	ip			int
	eop			int
}

func CreateNew(tokens []types.Token) (*Interpreter, error) {
	var stack util.Stack[StackValue] = util.Stack[StackValue]{}

	program, err := ProcessTokens(tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to create interpreter: %v", err)
	}

	return &Interpreter{
		stack: stack,
		program: program,
		ip: 0,
		eop: int(len(program)),
	}, nil
}

func (ip *Interpreter) Run() (bool, error) {
	if ip == nil {
		return false, fmt.Errorf("failed to run interpreter. instance is nil.")
	}

	for ip.ip < ip.eop {
		status, err := ip.Step()
		if err != nil {
			return false, fmt.Errorf("failed to run interpreter. interpreter step failed: %v", err)
		}

		if !status {
			return false, nil
		}
	}

	return true, nil
}

func (ip *Interpreter) runtimev(format string, args...any) bool {
	if !IsVerboseRuntime() {
		return false
	}

	colorYellow := "\033[33m"
	colorReset := "\033[0m"

	var sb strings.Builder
	sb.WriteString(colorYellow)
	sb.WriteString("Runtime: ")
	sb.WriteString(colorReset)
	sb.WriteString(format)

	format = sb.String()

	fmt.Printf(format, args)

	return false
}

func (ip *Interpreter) runtimeverr(format string, args...any) (bool, error) {
	return ip.runtimev(format, args...), nil
}

func (ip *Interpreter) Step() (bool, error) {
	if ip == nil {
		return false, fmt.Errorf("failed to step interpreter. instance is nil.")
	}

	if ip.ip >= ip.eop {
		return ip.runtimeverr("failed to run step. ip >= eop.")
	}

	inst := &ip.program[ip.ip]
	token := inst.Token

	if token == nil {
		return ip.runtimeverr("failed to run step. token at index %d is nil\n", ip.ip)
	}

	ip.runtimev("current stack size %d\n", ip.stack.Len())

	if IsVerboseRuntime() {
		tokenLog := token.Format()
		ip.runtimev("current token [%s] at ip %d\n", tokenLog, ip.ip)
	}

	if token.Equals("", types.TokenTypeNumber) {
		num, ok := token.GetNumberValue()
		if !ok {
			return ip.runtimeverr("failed to run step. token at index %d is has invalid number\n", ip.ip)
		}

		ip.runtimev("pushing %d\n", num)
		ip.ipush(num)
	} else if token.Equals("", types.TokenTypeString) {
		str, ok := token.GetStringValue()
		if !ok {
			return ip.runtimeverr("failed to run step. token at index %d has invalid string\n", ip.ip)
		}

		ip.runtimev("pushing \"%s\"\n", strings.ReplaceAll(str, "\n", "\\n"))
		ip.spush(str)
	} else if token.Equals(".", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. log operator failed. stack is empty.\n")
		}

		ip.runtimev("logging top value.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. log operator failed. failed to get value: %v\n", err)
		}

		n1, ok := v1.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. log operator failed. failed to get int value.")
		}
		ip.runtimev("popped %d\n", n1)

		fmt.Printf("%d\n", n1)
	} else if token.Equals("puts", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. puts operator failed. stack is empty.\n")
		}

		ip.runtimev("printing top string value.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. puts operator failed. failed to get value: %v\n", err)
		}

		s1, ok := v1.Left()
		if !ok {
			return ip.runtimeverr("failed to run step. puts operator failed. value is not a string.\n")
		}
		ip.runtimev("popped \"%s\"\n", s1)

		fmt.Printf("%s", s1)
	} else if token.Equals("int", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. int operator failed. stack is empty.\n")
		}

		ip.runtimev("parsing string to int.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. int operator failed. failed to get value: %v\n", err)
		}

		s1, ok := v1.Left()
		if !ok {
			return ip.runtimeverr("failed to run step. int operator failed. value is not a string.\n")
		}
		ip.runtimev("popped \"%s\"\n", s1)

		parsed, err := strconv.Atoi(s1)
		if err != nil {
			ip.runtimev("parse failed, pushing 0\n")
			ip.ipush(0)
			ip.runtimev("pushed %d\n", 0)
			ip.ipush(0)
			ip.runtimev("pushed %d (false)\n", 0)
		} else {
			ip.ipush(parsed)
			ip.runtimev("pushed %d\n", parsed)
			ip.ipush(1)
			ip.runtimev("pushed %d (true)\n", 1)
		}
	} else if token.Equals("string", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. string operator failed. stack is empty.\n")
		}

		ip.runtimev("converting value to string.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. string operator failed. failed to get value: %v\n", err)
		}

		if s1, ok := v1.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			ip.spush(s1)
			ip.runtimev("pushed \"%s\"\n", s1)
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("popped %d\n", n1)
			result := strconv.Itoa(n1)
			ip.spush(result)
			ip.runtimev("pushed \"%s\"\n", result)
		}
	} else if token.Equals("+", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. + operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("adding two numbers, concating two strings, or concating string + number.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. + operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. + operator failed. failed to get second value: %v\n", err)
		}

		if s2, ok := v2.Left(); ok {
			if s1, ok := v1.Left(); ok {
				ip.runtimev("popped \"%s\"\n", s1)
				ip.runtimev("popped \"%s\"\n", s2)

				result := s2 + s1

				ip.spush(result)
				ip.runtimev("pushed \"%s\"\n", result)
			} else if n1, ok := v1.Right(); ok {
				ip.runtimev("popped %d\n", n1)
				ip.runtimev("popped \"%s\"\n", s2)

				result := fmt.Sprintf("%s%d", s2, n1)

				ip.spush(result)
				ip.runtimev("pushed \"%s\"\n", result)
			} else {
				return ip.runtimeverr("failed to run step. + operator failed. failed to get first value type.\n")
			}
		} else if n2, ok := v2.Right(); ok {
			if n1, ok := v1.Right(); ok {
				ip.runtimev("popped %d\n", n1)
				ip.runtimev("popped %d\n", n2)

				result := n2 + n1

				ip.ipush(result)
				ip.runtimev("pushed %d\n", result)
			} else if v1.IsLeft() {
				return ip.runtimeverr("failed to run step. + operator failed. second value is int, thus first value must be int.\n")
			} else {
				return ip.runtimeverr("failed to run step. + operator failed. failed to get first value type.\n")
			}
		} else {
			return ip.runtimeverr("failed to run step. + operator failed. failed to get second value type.\n")
		}
	} else if token.Equals("-", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. - operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("subtracting two numbers.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. - operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. - operator failed. failed to get second value: %v\n", err)
		}

		if v2.IsLeft() {
			return ip.runtimeverr("failed to run step. - operator failed. cannot subtract from string.\n")
		}

		if v1.IsLeft() {
			return ip.runtimeverr("failed to run step. - operator failed. cannot subtract string.\n")
		}

		n1, ok := v1.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. - operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. - operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		result := n2 - n1

		ip.ipush(result)
		ip.runtimev("pushed %d\n", result)
	} else if token.Equals("*", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. * operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("multiplying two numbers.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. * operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. * operator failed. failed to get second value: %v\n", err)
		}

		if v2.IsLeft() {
			return ip.runtimeverr("failed to run step. * operator failed. cannot multiply string.\n")
		}

		if v1.IsLeft() {
			return ip.runtimeverr("failed to run step. * operator failed. cannot multiply string.\n")
		}

		n1, ok := v1.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. * operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. * operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		result := n1 * n2

		ip.ipush(result)
		ip.runtimev("pushed %d\n", result)
	} else if token.Equals("/", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. / operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("dividing two numbers.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. / operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. / operator failed. failed to get second value: %v\n", err)
		}

		if v2.IsLeft() {
			return ip.runtimeverr("failed to run step. / operator failed. cannot divide string.\n")
		}

		if v1.IsLeft() {
			return ip.runtimeverr("failed to run step. / operator failed. cannot divide string.\n")
		}

		n1, ok := v1.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. / operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. / operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		var result int
		if n1 != 0 {
			result = n2 / n1
		}

		ip.ipush(result)
		ip.runtimev("pushed %d\n", result)
	} else if token.Equals("%", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. % operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("modulo two numbers.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. % operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. % operator failed. failed to get second value: %v\n", err)
		}

		if v2.IsLeft() {
			return ip.runtimeverr("failed to run step. % operator failed. cannot modulo string.\n")
		}

		if v1.IsLeft() {
			return ip.runtimeverr("failed to run step. % operator failed. cannot modulo string.\n")
		}

		n1, ok := v1.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. % operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. % operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		var result int
		if n1 != 0 {
			result = n2 % n1
		}

		ip.ipush(result)
		ip.runtimev("pushed %d\n", result)
	} else if token.Equals("dup", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. dup operator failed. stack is empty.\n")
		}

		ip.runtimev("duplicate top stack value.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. dup operator failed. failed to get value: %v\n", err)
		}

		if s1, ok := v1.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			ip.spush(s1)
			ip.runtimev("pushed \"%s\"\n", s1)
			ip.spush(s1)
			ip.runtimev("pushed \"%s\"\n", s1)
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("popped %d\n", n1)
			ip.ipush(n1)
			ip.runtimev("pushed %d\n", n1)
			ip.ipush(n1)
			ip.runtimev("pushed %d\n", n1)
		}
	} else if token.Equals("drop", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. drop operator failed. stack is empty.\n")
		}

		ip.runtimev("dropping top stack value.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. drop operator failed. failed to pop value: %v\n", err)
		}

		if s1, ok := v1.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("popped %d\n", n1)
		}
	} else if token.Equals("swap", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. swap operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("swapping top two stack values.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. swap operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. swap operator failed. failed to get second value: %v\n", err)
		}

		if s1, ok := v1.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("popped %d\n", n1)
		}

		if s2, ok := v2.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s2)
		} else if n2, ok := v2.Right(); ok {
			ip.runtimev("popped %d\n", n2)
		}

		ip.push(v1)
		if s1, ok := v1.Left(); ok {
			ip.runtimev("pushed \"%s\"\n", s1)
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("pushed %d\n", n1)
		}

		ip.push(v2)
		if s2, ok := v2.Left(); ok {
			ip.runtimev("pushed \"%s\"\n", s2)
		} else if n2, ok := v2.Right(); ok {
			ip.runtimev("pushed %d\n", n2)
		}
	} else if token.Equals("over", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. over operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("pushing 2nd top-most value onto stack.\n")
		v2, err := ip.peekOffset(1)
		if err != nil {
			return ip.runtimeverr("failed to run step. over operator failed. failed to get second value: %v\n", err)
		}

		if s2, ok := v2.Left(); ok {
			ip.runtimev("peeked(1) \"%s\"\n", s2)
			ip.spush(s2)
			ip.runtimev("pushed \"%s\"\n", s2)
		} else if n2, ok := v2.Right(); ok {
			ip.runtimev("peeked(1) %d\n", n2)
			ip.ipush(n2)
			ip.runtimev("pushed %d\n", n2)
		}
	} else if token.Equals("2dup", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. 2dup operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("two-duping top stack values.\n")
		v1, err := ip.peekOffset(1)
		if err != nil {
			return ip.runtimeverr("failed to run step. 2dup operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.peekOffset(0)
		if err != nil {
			return ip.runtimeverr("failed to run step. 2dup operator failed. failed to get second value: %v", err)
		}

		if s1, ok := v1.Left(); ok {
			ip.runtimev("peeked(1) \"%s\"\n", s1)
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("peeked(1) %d\n", n1)
		}

		if s2, ok := v2.Left(); ok {
			ip.runtimev("peeked(0) \"%s\"\n", s2)
		} else if n2, ok := v2.Right(); ok {
			ip.runtimev("peeked(0) %d\n", n2)
		}

		ip.push(v1)
		if s1, ok := v1.Left(); ok {
			ip.runtimev("pushed \"%s\"\n", s1)
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("pushed %d\n", n1)
		}

		ip.push(v2)
		if s2, ok := v2.Left(); ok {
			ip.runtimev("pushed \"%s\"\n", s2)
		} else if n2, ok := v2.Right(); ok {
			ip.runtimev("pushed %d\n", n2)
		}
	} else if token.Equals("2swap", types.TokenTypeKeyword) {
		if ip.stack.Len() < 4 {
			return ip.runtimeverr("failed to run step. 2swap operator failed. stack size is %d. 4 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("two-swapping top stack values.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. 2swap operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. 2swap operator failed. failed to get second value: %v\n", err)
		}

		v3, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. 2swap operator failed. failed to get third value: %v\n", err)
		}

		v4, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. 2swap operator failed. failed to get fourth value: %v\n", err)
		}

		if s1, ok := v1.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("popped %d\n", n1)
		}

		if s2, ok := v2.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s2)
		} else if n2, ok := v2.Right(); ok {
			ip.runtimev("popped %d\n", n2)
		}

		if s3, ok := v3.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s3)
		} else if n3, ok := v3.Right(); ok {
			ip.runtimev("popped %d\n", n3)
		}

		if s4, ok := v4.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s4)
		} else if n4, ok := v4.Right(); ok {
			ip.runtimev("popped %d\n", n4)
		}

		ip.push(v2)
		if s2, ok := v2.Left(); ok {
			ip.runtimev("pushed \"%s\"\n", s2)
		} else if n2, ok := v2.Right(); ok {
			ip.runtimev("pushed %d\n", n2)
		}

		ip.push(v1)
		if s1, ok := v1.Left(); ok {
			ip.runtimev("pushed \"%s\"\n", s1)
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("pushed %d\n", n1)
		}

		ip.push(v4)
		if s4, ok := v4.Left(); ok {
			ip.runtimev("pushed \"%s\"\n", s4)
		} else if n4, ok := v4.Right(); ok {
			ip.runtimev("pushed %d\n", n4)
		}

		ip.push(v3)
		if s3, ok := v3.Left(); ok {
			ip.runtimev("pushed \"%s\"\n", s3)
		} else if n3, ok := v3.Right(); ok {
			ip.runtimev("pushed %d\n", n3)
		}
	} else if token.Equals("if", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. if operator failed. stack is empty.\n")
		}

		ip.runtimev("validating if-condition.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. if operator failed. failed to get condition value: %v\n", err)
		}

		var truthy bool = false
		if s1, ok := v1.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			truthy = len(s1) > 0
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("popped %d\n", n1)
			truthy = n1 != 0
		}

		if !truthy {
			ip.ip = int(inst.Next)
			return true, nil
		}
	} else if token.Equals("else", types.TokenTypeKeyword) {
		if inst.Next != -1 {
			ip.ip = int(inst.Next)
			return true, nil
		}
	} else if token.Equals("do", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. do operator failed. stack is empty.\n")
		}

		ip.runtimev("validating do-condition.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. do operator failed. failed to get condition value: %v\n", err)
		}

		var truthy bool = false
		if s1, ok := v1.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			truthy = len(s1) > 0
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("popped %d\n", n1)
			truthy = n1 != 0
		}

		if !truthy {
			ip.ip = int(inst.Next)
			return true, nil
		}
	} else if token.Equals("end", types.TokenTypeKeyword) {
		if inst.Next != -1 {
			ip.ip = int(inst.Next)
			return true, nil
		}
	} else if token.Equals("while", types.TokenTypeKeyword) {
		ip.runtimev("while (do nothing).\n")
	} else if token.Equals("=", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. = operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("equality-check top stack values.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. = operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. = operator failed. failed to get second value: %v\n", err)
		}

		var result int = 0
		var name string = "false"

		if s1, ok := v1.Left(); ok {
			if s2, ok := v2.Left(); ok {
				ip.runtimev("popped \"%s\"\n", s1)
				ip.runtimev("popped \"%s\"\n", s2)
				if s1 == s2 {
					result = 1
					name = "true"
				}
			} else {
				ip.runtimev("popped int\n")
				ip.runtimev("popped \"%s\"\n", s1)
			}
		} else if n1, ok := v1.Right(); ok {
			if n2, ok := v2.Right(); ok {
				ip.runtimev("popped %d\n", n1)
				ip.runtimev("popped %d\n", n2)
				if n1 == n2 {
					result = 1
					name = "true"
				}
			} else {
				ip.runtimev("popped string\n")
				ip.runtimev("popped %d\n", n1)
			}
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals("!=", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. != operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("inequality-check top stack values.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. != operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. != operator failed. failed to get second value: %v\n", err)
		}

		var result int = 1
		var name string = "true"

		if s1, ok := v1.Left(); ok {
			if s2, ok := v2.Left(); ok {
				ip.runtimev("popped \"%s\"\n", s1)
				ip.runtimev("popped \"%s\"\n", s2)
				if s1 == s2 {
					result = 0
					name = "false"
				}
			} else {
				ip.runtimev("popped int\n")
				ip.runtimev("popped \"%s\"\n", s1)
			}
		} else if n1, ok := v1.Right(); ok {
			if n2, ok := v2.Right(); ok {
				ip.runtimev("popped %d\n", n1)
				ip.runtimev("popped %d\n", n2)
				if n1 == n2 {
					result = 0
					name = "false"
				}
			} else {
				ip.runtimev("popped string\n")
				ip.runtimev("popped %d\n", n1)
			}
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals("<", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. < operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("less-check top stack values.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. < operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. < operator failed. failed to get second value: %v\n", err)
		}

		var result int = 0
		var name string = "false"

		if s1, ok := v1.Left(); ok {
			if s2, ok := v2.Left(); ok {
				ip.runtimev("popped \"%s\"\n", s1)
				ip.runtimev("popped \"%s\"\n", s2)
				if s2 < s1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr("failed to run step. < operator failed. cannot compare string and int.\n")
			}
		} else if n1, ok := v1.Right(); ok {
			if n2, ok := v2.Right(); ok {
				ip.runtimev("popped %d\n", n1)
				ip.runtimev("popped %d\n", n2)
				if n2 < n1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr("failed to run step. < operator failed. cannot compare int and string.\n")
			}
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals(">", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. > operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("greater-check top stack values.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. > operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. > operator failed. failed to get second value: %v\n", err)
		}

		var result int = 0
		var name string = "false"

		if s1, ok := v1.Left(); ok {
			if s2, ok := v2.Left(); ok {
				ip.runtimev("popped \"%s\"\n", s1)
				ip.runtimev("popped \"%s\"\n", s2)
				if s2 > s1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr("failed to run step. > operator failed. cannot compare string and int.\n")
			}
		} else if n1, ok := v1.Right(); ok {
			if n2, ok := v2.Right(); ok {
				ip.runtimev("popped %d\n", n1)
				ip.runtimev("popped %d\n", n2)
				if n2 > n1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr("failed to run step. > operator failed. cannot compare int and string.\n")
			}
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals("<=", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. <= operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("lessorequal-check top stack values.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. <= operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. <= operator failed. failed to get second value: %v\n", err)
		}

		var result int = 0
		var name string = "false"

		if s1, ok := v1.Left(); ok {
			if s2, ok := v2.Left(); ok {
				ip.runtimev("popped \"%s\"\n", s1)
				ip.runtimev("popped \"%s\"\n", s2)
				if s2 <= s1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr("failed to run step. <= operator failed. cannot compare string and int.\n")
			}
		} else if n1, ok := v1.Right(); ok {
			if n2, ok := v2.Right(); ok {
				ip.runtimev("popped %d\n", n1)
				ip.runtimev("popped %d\n", n2)
				if n2 <= n1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr("failed to run step. <= operator failed. cannot compare int and string.\n")
			}
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals(">=", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. >= operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("greaterorequal-check top stack values.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. >= operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. >= operator failed. failed to get second value: %v\n", err)
		}

		var result int = 0
		var name string = "false"

		if s1, ok := v1.Left(); ok {
			if s2, ok := v2.Left(); ok {
				ip.runtimev("popped \"%s\"\n", s1)
				ip.runtimev("popped \"%s\"\n", s2)
				if s2 >= s1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr("failed to run step. >= operator failed. cannot compare string and int.\n")
			}
		} else if n1, ok := v1.Right(); ok {
			if n2, ok := v2.Right(); ok {
				ip.runtimev("popped %d\n", n1)
				ip.runtimev("popped %d\n", n2)
				if n2 >= n1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr("failed to run step. >= operator failed. cannot compare int and string.\n")
			}
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals("!", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. ! operator failed. stack is empty.\n")
		}

		ip.runtimev("logicalnot top stack value.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. ! operator failed. failed to get value: %v\n", err)
		}

		var truthy bool = false
		if s1, ok := v1.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			truthy = len(s1) > 0
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("popped %d\n", n1)
			truthy = n1 != 0
		}

		var result int = 0
		var name string = "false"
		if !truthy {
			result = 1
			name = "true"
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals("~", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. ~ operator failed. stack is empty.\n")
		}

		ip.runtimev("bitwisenot top stack value.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. ~ operator failed. failed to get value: %v\n", err)
		}

		if v1.IsLeft() {
			return ip.runtimeverr("failed to run step. ~ operator failed. cannot bitwise not string.\n")
		}

		n1, ok := v1.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. ~ operator failed. failed to get value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		result := ^n1

		ip.ipush(result)
		ip.runtimev("pushed %d\n", result)
	} else if token.Equals("&", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. & operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("bitwiseand two numbers.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. & operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. & operator failed. failed to get second value: %v\n", err)
		}

		if v2.IsLeft() {
			return ip.runtimeverr("failed to run step. & operator failed. cannot bitwise and string.\n")
		}

		if v1.IsLeft() {
			return ip.runtimeverr("failed to run step. & operator failed. cannot bitwise and string.\n")
		}

		n1, ok := v1.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. & operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. & operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		result := n2 & n1

		ip.ipush(result)
		ip.runtimev("pushed %d\n", result)
	} else if token.Equals("|", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. | operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("bitwiseor two numbers.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. | operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. | operator failed. failed to get second value: %v\n", err)
		}

		if v2.IsLeft() {
			return ip.runtimeverr("failed to run step. | operator failed. cannot bitwise or string.\n")
		}

		if v1.IsLeft() {
			return ip.runtimeverr("failed to run step. | operator failed. cannot bitwise or string.\n")
		}

		n1, ok := v1.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. | operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. | operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		result := n2 | n1

		ip.ipush(result)
		ip.runtimev("pushed %d\n", result)
	} else if token.Equals("^", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. ^ operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("bitwisexor two numbers.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. ^ operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. ^ operator failed. failed to get second value: %v\n", err)
		}

		if v2.IsLeft() {
			return ip.runtimeverr("failed to run step. ^ operator failed. cannot bitwise xor string.\n")
		}

		if v1.IsLeft() {
			return ip.runtimeverr("failed to run step. ^ operator failed. cannot bitwise xor string.\n")
		}

		n1, ok := v1.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. ^ operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Right()
		if !ok {
			return ip.runtimeverr("failed to run step. ^ operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		result := n2 ^ n1

		ip.ipush(result)
		ip.runtimev("pushed %d\n", result)
	} else if token.Equals("&&", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. && operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("logicaland top stack values.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. && operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. && operator failed. failed to get second value: %v\n", err)
		}

		var truthy1 bool = false
		var truthy2 bool = false

		if s1, ok := v1.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			truthy1 = len(s1) > 0
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("popped %d\n", n1)
			truthy1 = n1 != 0
		}

		if s2, ok := v2.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s2)
			truthy2 = len(s2) > 0
		} else if n2, ok := v2.Right(); ok {
			ip.runtimev("popped %d\n", n2)
			truthy2 = n2 != 0
		}

		var result int = 0
		var name string = "false"
		if truthy1 && truthy2 {
			result = 1
			name = "true"
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals("||", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. || operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("logicalor top stack values.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. || operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. || operator failed. failed to get second value: %v\n", err)
		}

		var truthy1 bool = false
		var truthy2 bool = false

		if s1, ok := v1.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			truthy1 = len(s1) > 0
		} else if n1, ok := v1.Right(); ok {
			ip.runtimev("popped %d\n", n1)
			truthy1 = n1 != 0
		}

		if s2, ok := v2.Left(); ok {
			ip.runtimev("popped \"%s\"\n", s2)
			truthy2 = len(s2) > 0
		} else if n2, ok := v2.Right(); ok {
			ip.runtimev("popped %d\n", n2)
			truthy2 = n2 != 0
		}

		var result int = 0
		var name string = "false"
		if truthy1 || truthy2 {
			result = 1
			name = "true"
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else {
		if IsVerboseRuntime() {
			tokenLog := token.Format()
			ip.runtimev("invalid token [%s] at ip %d\n", tokenLog, ip.ip)
		}

		return false, fmt.Errorf("failed to run step. unknown operator at ip %d\n", ip.ip)
	}

	ip.ip++

	return true, nil
}

func (ip *Interpreter) pop() (StackValue, error) {
	var value StackValue
	value, ok := ip.stack.Pop()
	if !ok {
		return value, fmt.Errorf("failed to ipop. stack is empty.")
	}

	ip.runtimev("new stack size %d\n", ip.stack.Len())

	return value, nil
}

func (ip *Interpreter) peek() (StackValue, error) {
	return ip.peekOffset(0)
}

func (ip *Interpreter) peekOffset(offset int) (StackValue, error) {
	var value StackValue
	idx := ip.stack.Len() - 1 - offset
	if idx < 0 {
		return value, fmt.Errorf("failed to ipeek. idx(%d) < 0", idx)
	}
	if idx >= ip.stack.Len() {
		return value, fmt.Errorf("failed to ipeek. idx(%d) >= len(%d)", idx, ip.stack.Len())
	}

	return ip.stack[ip.stack.Len() - 1 - offset], nil
}

func (ip *Interpreter) push(value StackValue) {
	ip.stack.Push(value)
	ip.runtimev("new stack size %d\n", ip.stack.Len())
}

func (ip *Interpreter) ipush(item int) {
	ip.push(util.NewRight[string, int](item))
}

func (ip *Interpreter) spush(item string) {
	ip.push(util.NewLeft[string, int](item))
}

