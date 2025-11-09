package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ktnuity/wet/internal/tools"
	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)


type StackValue struct {
	tstring *string
	tint *int
	tpath *string
}

func StackString(val string) StackValue {return StackValue{tstring: &val}}
func StackInt(val int) StackValue {return StackValue{tint: &val}}
func StackPath(val string) StackValue {return StackValue{tpath: &val}}

func (sv StackValue) IsString() bool { return sv.tstring != nil }
func (sv StackValue) IsInt() bool { return sv.tint != nil }
func (sv StackValue) IsPath() bool { return sv.tpath != nil }

func (sv StackValue) IsPrimary() bool { return sv.tstring != nil || sv.tint != nil }
func (sv StackValue) IsAny() bool { return sv.tstring != nil || sv.tint != nil || sv.tpath != nil }

func (sv StackValue) String() (string, bool) {
	if sv.tstring != nil {
		return *sv.tstring, true
	}

	return "", false
}

func (sv StackValue) Int() (int, bool) {
	if sv.tint != nil {
		return *sv.tint, true
	}

	return 0, false
}

func (sv StackValue) Path() (string, bool) {
	if sv.tpath != nil {
		return *sv.tpath, true
	}

	return "", false
}

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
		return nil, fmt.Errorf("failed to create interpreter: %v", err)
	}

	return &Interpreter{
		stack: stack,
		program: program,
		memory: make(map[string]StackValue),
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
	fmt.Printf(sb.String(), args...)

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
		err := ip.ipush(num)
		if err != nil {
			return false, fmt.Errorf("failed to run step. number push operator failed. failure pushing value: %v", err)
		}
	} else if token.Equals("", types.TokenTypeString) {
		str, ok := token.GetStringValue()
		if !ok {
			return ip.runtimeverr("failed to run step. token at index %d has invalid string\n", ip.ip)
		}

		ip.runtimev("pushing \"%s\"\n", strings.ReplaceAll(str, "\n", "\\n"))
		err := ip.spush(str)
		if err != nil {
			return false, fmt.Errorf("failed to run step. string push operator failed. failure pushing value: %v", err)
		}
	} else if token.Equals("", types.TokenTypePath) {
		str, ok := token.GetPathValue()
		if !ok {
			return ip.runtimeverr("failed to run step. token at index %d has invalid path\n", ip.ip)
		}

		ip.runtimev("pushing path(%s)\n", str)
		err := ip.ppush(str)
		if err != nil {
			return false, fmt.Errorf("failed to run step. path push operator failed. failure pushing value: %v", err)
		}
	} else if token.Equals(".", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. log operator failed. stack is empty.\n")
		}

		ip.runtimev("logging top value.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. log operator failed. failed to get value: %v\n", err)
		}

		n1, ok := v1.Int()
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

		s1, ok := v1.String()
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

		s1, ok := v1.String()
		if !ok {
			return ip.runtimeverr("failed to run step. int operator failed. value is not a string.\n")
		}
		ip.runtimev("popped \"%s\"\n", s1)

		parsed, err := strconv.Atoi(s1)
		if err != nil {
			ip.runtimev("parse failed, pushing 0\n")
			err := ip.ipush(0)
			if err != nil {
				return false, fmt.Errorf("failed to run step. int operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed %d\n", 0)
			err = ip.ipush(0)
			if err != nil {
				return false, fmt.Errorf("failed to run step. int operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed %d (false)\n", 0)
		} else {
			err := ip.ipush(parsed)
			if err != nil {
				return false, fmt.Errorf("failed to run step. int operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed %d\n", parsed)
			err = ip.ipush(1)
			if err != nil {
				return false, fmt.Errorf("failed to run step. int operator failed. failure pushing value: %v", err)
			}
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

		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			err := ip.spush(s1)
			if err != nil {
				return false, fmt.Errorf("failed to run step. string operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed \"%s\"\n", s1)
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d\n", n1)
			result := strconv.Itoa(n1)
			err := ip.spush(result)
			if err != nil {
				return false, fmt.Errorf("failed to run step. string operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed \"%s\"\n", result)
		}
	} else if token.Equals("store", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. store operator failed. stack is empty.\n")
		}

		ip.runtimev("storing value in memory.\n")
		vName, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. store operator failed. failed to get name: %v\n", err)
		}

		name, ok := vName.String()
		if !ok {
			return ip.runtimeverr("failed to run step. store operator failed. failed to get name value.\n")
		}

		vValue, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. store operator failed. failed to get value: %v\n", err)
		}

		ip.store(name, vValue)
	} else if token.Equals("load", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. load operator failed. stack is empty.\n")
		}

		ip.runtimev("loading value from memory.\n")
		vName, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. load operator failed. failed to get name: %v\n", err)
		}

		name, ok := vName.String()
		if !ok {
			return ip.runtimeverr("failed to run step. load operator failed. failed to get name value.\n")
		}

		value, err := ip.load(name)
		if err != nil {
			return ip.runtimeverr("failed to run step. load operator failed. failed to load memory: %v\n", err)
		}

		ip.push(value)

		left, okLeft := value.String()
		right, okRight := value.Int()
		if okLeft {
			ip.runtimev("pushed \"%s\"\n", left)
		} else if okRight {
			ip.runtimev("pushed %d\n", right)
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

		if s2, ok := v2.String(); ok {
			if s1, ok := v1.String(); ok {
				ip.runtimev("popped \"%s\"\n", s1)
				ip.runtimev("popped \"%s\"\n", s2)

				result := s2 + s1

				err := ip.spush(result)
				if err != nil {
					return false, fmt.Errorf("failed to run step. + operator failed. failure pushing value: %v", err)
				}
				ip.runtimev("pushed \"%s\"\n", result)
			} else if n1, ok := v1.Int(); ok {
				ip.runtimev("popped %d\n", n1)
				ip.runtimev("popped \"%s\"\n", s2)

				result := fmt.Sprintf("%s%d", s2, n1)

				err := ip.spush(result)
				if err != nil {
					return false, fmt.Errorf("failed to run step. + operator failed. failure pushing value: %v", err)
				}
				ip.runtimev("pushed \"%s\"\n", result)
			} else {
				return ip.runtimeverr("failed to run step. + operator failed. failed to get first value type.\n")
			}
		} else if n2, ok := v2.Int(); ok {
			if n1, ok := v1.Int(); ok {
				ip.runtimev("popped %d\n", n1)
				ip.runtimev("popped %d\n", n2)

				result := n2 + n1

				err := ip.ipush(result)
				if err != nil {
					return false, fmt.Errorf("failed to run step. + operator failed. failure pushing value: %v", err)
				}
				ip.runtimev("pushed %d\n", result)
			} else if v1.IsString() {
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

		if v2.IsString() {
			return ip.runtimeverr("failed to run step. - operator failed. cannot subtract from string.\n")
		}

		if v1.IsString() {
			return ip.runtimeverr("failed to run step. - operator failed. cannot subtract string.\n")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. - operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. - operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		result := n2 - n1

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. - operator failed. failure pushing value: %v", err)
		}
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

		if v2.IsString() {
			return ip.runtimeverr("failed to run step. * operator failed. cannot multiply string.\n")
		}

		if v1.IsString() {
			return ip.runtimeverr("failed to run step. * operator failed. cannot multiply string.\n")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. * operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. * operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		result := n1 * n2

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. * operator failed. failure pushing value: %v", err)
		}
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

		if v2.IsString() {
			return ip.runtimeverr("failed to run step. / operator failed. cannot divide string.\n")
		}

		if v1.IsString() {
			return ip.runtimeverr("failed to run step. / operator failed. cannot divide string.\n")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. / operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. / operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		var result int
		if n1 != 0 {
			result = n2 / n1
		}

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. / operator failed. failure pushing value: %v", err)
		}
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

		if v2.IsString() {
			return ip.runtimeverr("failed to run step. % operator failed. cannot modulo string.\n")
		}

		if v1.IsString() {
			return ip.runtimeverr("failed to run step. % operator failed. cannot modulo string.\n")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. % operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. % operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		var result int
		if n1 != 0 {
			result = n2 % n1
		}

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. %% operator failed. failure pushing value: %v", err)
		}
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

		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			err := ip.spush(s1)
			if err != nil {
				return false, fmt.Errorf("failed to run step. dup operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed \"%s\"\n", s1)
			err = ip.spush(s1)
			if err != nil {
				return false, fmt.Errorf("failed to run step. dup operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed \"%s\"\n", s1)
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d\n", n1)
			err := ip.ipush(n1)
			if err != nil {
				return false, fmt.Errorf("failed to run step. dup operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed %d\n", n1)
			err = ip.ipush(n1)
			if err != nil {
				return false, fmt.Errorf("failed to run step. dup operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed %d\n", n1)
		} else if p1, ok := v1.Path(); ok {
			ip.runtimev("popped %s\n", p1)
			err := ip.ppush(p1)
			if err != nil {
				return false, fmt.Errorf("failed to run step. dup operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed %s\n", p1)
			err = ip.ppush(p1)
			if err != nil {
				return false, fmt.Errorf("failed to run step. dup operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed %s\n", p1)
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

		if s1, ok := v1.String(); ok {
			ip.runtimev("dropped \"%s\"\n", s1)
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("dropped %d\n", n1)
		} else if p1, ok := v1.Path(); ok {
			ip.runtimev("dropped %s\n", p1)
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

		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d\n", n1)
		} else if p1, ok := v1.Path(); ok {
			ip.runtimev("dropped %s\n", p1)
		}

		if s2, ok := v2.String(); ok {
			ip.runtimev("popped \"%s\"\n", s2)
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("popped %d\n", n2)
		} else if p2, ok := v2.Path(); ok {
			ip.runtimev("dropped %s\n", p2)
		}

		err = ip.push(v1)
		if err != nil {
			return false, fmt.Errorf("failed to run step. swap operator failed. failure pushing value: %v", err)
		}
		if s1, ok := v1.String(); ok {
			ip.runtimev("pushed \"%s\"\n", s1)
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("pushed %d\n", n1)
		} else if p1, ok := v1.Path(); ok {
			ip.runtimev("pushed %s\n", p1)
		}

		err = ip.push(v2)
		if err != nil {
			return false, fmt.Errorf("failed to run step. swap operator failed. failure pushing value: %v", err)
		}
		if s2, ok := v2.String(); ok {
			ip.runtimev("pushed \"%s\"\n", s2)
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("pushed %d\n", n2)
		} else if p2, ok := v2.Path(); ok {
			ip.runtimev("pushed %s\n", p2)
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

		if s2, ok := v2.String(); ok {
			ip.runtimev("peeked(1) \"%s\"\n", s2)
			err := ip.spush(s2)
			if err != nil {
				return false, fmt.Errorf("failed to run step. over operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed \"%s\"\n", s2)
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("peeked(1) %d\n", n2)
			err := ip.ipush(n2)
			if err != nil {
				return false, fmt.Errorf("failed to run step. over operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed %d\n", n2)
		} else if p2, ok := v2.Path(); ok {
			ip.runtimev("peeked(1) %d\n", p2)
			err := ip.ppush(p2)
			if err != nil {
				return false, fmt.Errorf("failed to run step. over operator failed. failure pushing value: %v", err)
			}
			ip.runtimev("pushed %s\n", p2)
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

		if s1, ok := v1.String(); ok {
			ip.runtimev("peeked(1) \"%s\"\n", s1)
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("peeked(1) %d\n", n1)
		} else if p1, ok := v1.Path(); ok {
			ip.runtimev("peeked(1) %s\n", p1)
		}

		if s2, ok := v2.String(); ok {
			ip.runtimev("peeked(0) \"%s\"\n", s2)
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("peeked(0) %d\n", n2)
		} else if p2, ok := v2.Path(); ok {
			ip.runtimev("peeked(0) %s\n", p2)
		}

		err = ip.push(v1)
		if err != nil {
			return false, fmt.Errorf("failed to run step. 2dup operator failed. failure pushing value: %v", err)
		}
		if s1, ok := v1.String(); ok {
			ip.runtimev("pushed \"%s\"\n", s1)
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("pushed %d\n", n1)
		} else if p1, ok := v1.Path(); ok {
			ip.runtimev("pushed %s\n", p1)
		}

		err = ip.push(v2)
		if err != nil {
			return false, fmt.Errorf("failed to run step. 2dup operator failed. failure pushing value: %v", err)
		}
		if s2, ok := v2.String(); ok {
			ip.runtimev("pushed \"%s\"\n", s2)
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("pushed %d\n", n2)
		} else if p2, ok := v2.Path(); ok {
			ip.runtimev("pushed %s\n", p2)
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

		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d\n", n1)
		} else if p1, ok := v1.Path(); ok {
			ip.runtimev("popped %s\n", p1)
		}

		if s2, ok := v2.String(); ok {
			ip.runtimev("popped \"%s\"\n", s2)
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("popped %d\n", n2)
		} else if p2, ok := v2.Path(); ok {
			ip.runtimev("popped %s\n", p2)
		}

		if s3, ok := v3.String(); ok {
			ip.runtimev("popped \"%s\"\n", s3)
		} else if n3, ok := v3.Int(); ok {
			ip.runtimev("popped %d\n", n3)
		} else if p3, ok := v3.Path(); ok {
			ip.runtimev("popped %s\n", p3)
		}

		if s4, ok := v4.String(); ok {
			ip.runtimev("popped \"%s\"\n", s4)
		} else if n4, ok := v4.Int(); ok {
			ip.runtimev("popped %d\n", n4)
		} else if p4, ok := v4.Path(); ok {
			ip.runtimev("popped %s\n", p4)
		}

		err = ip.push(v2)
		if err != nil {
			return false, fmt.Errorf("failed to run step. 2swap operator failed. failure pushing value: %v", err)
		}
		if s2, ok := v2.String(); ok {
			ip.runtimev("pushed \"%s\"\n", s2)
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("pushed %d\n", n2)
		} else if p2, ok := v2.Path(); ok {
			ip.runtimev("pushed %s\n", p2)
		}

		err = ip.push(v1)
		if err != nil {
			return false, fmt.Errorf("failed to run step. 2swap operator failed. failure pushing value: %v", err)
		}
		if s1, ok := v1.String(); ok {
			ip.runtimev("pushed \"%s\"\n", s1)
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("pushed %d\n", n1)
		} else if p1, ok := v1.Path(); ok {
			ip.runtimev("pushed %s\n", p1)
		}

		err = ip.push(v4)
		if err != nil {
			return false, fmt.Errorf("failed to run step. 2swap operator failed. failure pushing value: %v", err)
		}
		if s4, ok := v4.String(); ok {
			ip.runtimev("pushed \"%s\"\n", s4)
		} else if n4, ok := v4.Int(); ok {
			ip.runtimev("pushed %d\n", n4)
		} else if p4, ok := v4.Path(); ok {
			ip.runtimev("pushed %s\n", p4)
		}

		err = ip.push(v3)
		if err != nil {
			return false, fmt.Errorf("failed to run step. 2swap operator failed. failure pushing value: %v", err)
		}
		if s3, ok := v3.String(); ok {
			ip.runtimev("pushed \"%s\"\n", s3)
		} else if n3, ok := v3.Int(); ok {
			ip.runtimev("pushed %d\n", n3)
		} else if p3, ok := v3.Path(); ok {
			ip.runtimev("pushed %s\n", p3)
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
		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			truthy = len(s1) > 0
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d\n", n1)
			truthy = n1 != 0
		}

		if !truthy {
			ip.ip = int(inst.Next)
			return true, nil
		}
	} else if token.Equals("unless", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. unless operator failed. stack is empty.\n")
		}

		ip.runtimev("validating unless-condition.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. unless operator failed. failed to get condition value: %v\n", err)
		}

		var truthy bool = false
		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			truthy = len(s1) > 0
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d\n", n1)
			truthy = n1 != 0
		}

		if truthy {
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
		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			truthy = len(s1) > 0
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d\n", n1)
			truthy = n1 != 0
		}

		if inst.Mode == DoModeUntil {
			truthy = !truthy
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
	} else if token.Equals("until", types.TokenTypeKeyword) {
		ip.runtimev("until (do nothing).\n")
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

		if s1, ok := v1.String(); ok {
			if s2, ok := v2.String(); ok {
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
		} else if n1, ok := v1.Int(); ok {
			if n2, ok := v2.Int(); ok {
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

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. = operator failed. failure pushing value: %v", err)
		}
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

		if s1, ok := v1.String(); ok {
			if s2, ok := v2.String(); ok {
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
		} else if n1, ok := v1.Int(); ok {
			if n2, ok := v2.Int(); ok {
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

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. != operator failed. failure pushing value: %v", err)
		}
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

		if s1, ok := v1.String(); ok {
			if s2, ok := v2.String(); ok {
				ip.runtimev("popped \"%s\"\n", s1)
				ip.runtimev("popped \"%s\"\n", s2)
				if s2 < s1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr("failed to run step. < operator failed. cannot compare string and int.\n")
			}
		} else if n1, ok := v1.Int(); ok {
			if n2, ok := v2.Int(); ok {
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

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. < operator failed. failure pushing value: %v", err)
		}
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

		if s1, ok := v1.String(); ok {
			if s2, ok := v2.String(); ok {
				ip.runtimev("popped \"%s\"\n", s1)
				ip.runtimev("popped \"%s\"\n", s2)
				if s2 > s1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr("failed to run step. > operator failed. cannot compare string and int.\n")
			}
		} else if n1, ok := v1.Int(); ok {
			if n2, ok := v2.Int(); ok {
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

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. > operator failed. failure pushing value: %v", err)
		}
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

		if s1, ok := v1.String(); ok {
			if s2, ok := v2.String(); ok {
				ip.runtimev("popped \"%s\"\n", s1)
				ip.runtimev("popped \"%s\"\n", s2)
				if s2 <= s1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr("failed to run step. <= operator failed. cannot compare string and int.\n")
			}
		} else if n1, ok := v1.Int(); ok {
			if n2, ok := v2.Int(); ok {
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

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. <= operator failed. failure pushing value: %v", err)
		}
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

		if s1, ok := v1.String(); ok {
			if s2, ok := v2.String(); ok {
				ip.runtimev("popped \"%s\"\n", s1)
				ip.runtimev("popped \"%s\"\n", s2)
				if s2 >= s1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr("failed to run step. >= operator failed. cannot compare string and int.\n")
			}
		} else if n1, ok := v1.Int(); ok {
			if n2, ok := v2.Int(); ok {
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

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. >= operator failed. failure pushing value: %v", err)
		}
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
		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			truthy = len(s1) > 0
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d\n", n1)
			truthy = n1 != 0
		}

		var result int = 0
		var name string = "false"
		if !truthy {
			result = 1
			name = "true"
		}

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. ! operator failed. failure pushing value: %v", err)
		}
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

		if v1.IsString() {
			return ip.runtimeverr("failed to run step. ~ operator failed. cannot bitwise not string.\n")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. ~ operator failed. failed to get value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		result := ^n1

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. ~ operator failed. failure pushing value: %v", err)
		}
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

		if v2.IsString() {
			return ip.runtimeverr("failed to run step. & operator failed. cannot bitwise and string.\n")
		}

		if v1.IsString() {
			return ip.runtimeverr("failed to run step. & operator failed. cannot bitwise and string.\n")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. & operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. & operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		result := n2 & n1

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. & operator failed. failure pushing value: %v", err)
		}
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

		if v2.IsString() {
			return ip.runtimeverr("failed to run step. | operator failed. cannot bitwise or string.\n")
		}

		if v1.IsString() {
			return ip.runtimeverr("failed to run step. | operator failed. cannot bitwise or string.\n")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. | operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. | operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		result := n2 | n1

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. | operator failed. failure pushing value: %v", err)
		}
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

		if v2.IsString() {
			return ip.runtimeverr("failed to run step. ^ operator failed. cannot bitwise xor string.\n")
		}

		if v1.IsString() {
			return ip.runtimeverr("failed to run step. ^ operator failed. cannot bitwise xor string.\n")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. ^ operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. ^ operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		result := n2 ^ n1

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. ^ operator failed. failure pushing value: %v", err)
		}
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

		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			truthy1 = len(s1) > 0
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d\n", n1)
			truthy1 = n1 != 0
		}

		if s2, ok := v2.String(); ok {
			ip.runtimev("popped \"%s\"\n", s2)
			truthy2 = len(s2) > 0
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("popped %d\n", n2)
			truthy2 = n2 != 0
		}

		var result int = 0
		var name string = "false"
		if truthy1 && truthy2 {
			result = 1
			name = "true"
		}

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. && operator failed. failure pushing value: %v", err)
		}
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

		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"\n", s1)
			truthy1 = len(s1) > 0
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d\n", n1)
			truthy1 = n1 != 0
		}

		if s2, ok := v2.String(); ok {
			ip.runtimev("popped \"%s\"\n", s2)
			truthy2 = len(s2) > 0
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("popped %d\n", n2)
			truthy2 = n2 != 0
		}

		var result int = 0
		var name string = "false"
		if truthy1 || truthy2 {
			result = 1
			name = "true"
		}

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. || operator failed. failure pushing value: %v", err)
		}
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals("download", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. download command failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("download command.\n")
		vDst, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. download command failed. failed to get destination value: %v\n", err)
		}

		vUrl, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. download command failed. failed to get url value: %v\n", err)
		}

		pDst, okDst := vDst.Path()
		if !okDst {
			return ip.runtimeverr("failed to run step. download command failed. failed to get destination path.\n")
		}

		sUrl, okUrl := vUrl.String()
		if !okUrl {
			return ip.runtimeverr("failed to run step. download command failed. failed to get url string.\n")
		}


		var result int = 1
		err = tools.ToolDownload(sUrl, pDst)
		if err != nil {
			ip.runtimev("failed to use download tool: %v\n", err)
			result = 0
		}

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. download command failed. failure pushing value: %v", err)
		}
	} else if token.Equals("readfile", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. readfile command failed. stack is empty.\n", ip.stack.Len())
		}

		ip.runtimev("readfile command.\n")
		vSrc, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. readfile command failed. failed to get source value: %v\n", err)
		}

		pSrc, okSrc := vSrc.Path()
		if !okSrc {
			return ip.runtimeverr("failed to run step. readfile command failed. failed to get source path.\n")
		}

		data, err := tools.ToolReadfile(pSrc)
		if err != nil {
			ip.runtimev("failed to use readfile too: %v\n", err)
			err = ip.ipush(0)
			if err != nil {
				return false, fmt.Errorf("failed to run step. readfile command failed. failure pushing value: %v", err)
			}
		} else {
			err = ip.spush(data)
			if err != nil {
				return false, fmt.Errorf("failed to run step. readfile command failed. failure pushing value: %v", err)
			}

			err = ip.ipush(1)
			if err != nil {
				return false, fmt.Errorf("failed to run step. readfile command failed. failure pushing value: %v", err)
			}
		}
	} else if token.Equals("move", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. move command failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("move command.\n")
		vDst, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. move command failed. failed to get destination value: %v\n", err)
		}

		vSrc, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. move command failed. failed to get source value: %v\n", err)
		}

		pDst, okDst := vDst.Path()
		if !okDst {
			return ip.runtimeverr("failed to run step. move command failed. failed to get destination path.\n")
		}

		pSrc, okSrc := vSrc.Path()
		if !okSrc {
			return ip.runtimeverr("failed to run step. move command failed. failed to get source path.\n")
		}

		var result int = 1
		err = tools.ToolMoveFile(pSrc, pDst)
		if err != nil {
			ip.runtimev("failed to use move tool: %v\n", err)
			result = 0
		}

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. move command failed. failure pushing value: %v", err)
		}
	} else if token.Equals("copy", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. copy command failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("copy command.\n")
		vDst, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. copy command failed. failed to get destination value: %v\n", err)
		}

		vSrc, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. copy command failed. failed to get source value: %v\n", err)
		}

		pDst, okDst := vDst.Path()
		if !okDst {
			return ip.runtimeverr("failed to run step. copy command failed. failed to get destination path.\n")
		}

		pSrc, okSrc := vSrc.Path()
		if !okSrc {
			return ip.runtimeverr("failed to run step. copy command failed. failed to get source path.\n")
		}

		occupied, err := tools.ToolCopyFile(pSrc, pDst)
		if err != nil {
			ip.runtimev("failed to use copy tool: %v\n", err)
			if occupied {
				// File already exists: push true false
				err = ip.ipush(1)
				if err != nil {
					return false, fmt.Errorf("failed to run step. copy command failed. failure pushing value: %v", err)
				}
				err = ip.ipush(0)
				if err != nil {
					return false, fmt.Errorf("failed to run step. copy command failed. failure pushing value: %v", err)
				}
			} else {
				// Other failure: push false false
				err = ip.ipush(0)
				if err != nil {
					return false, fmt.Errorf("failed to run step. copy command failed. failure pushing value: %v", err)
				}
				err = ip.ipush(0)
				if err != nil {
					return false, fmt.Errorf("failed to run step. copy command failed. failure pushing value: %v", err)
				}
			}
		} else {
			// Success: push true
			err = ip.ipush(1)
			if err != nil {
				return false, fmt.Errorf("failed to run step. copy command failed. failure pushing value: %v", err)
			}
		}
	} else if token.Equals("exist", types.TokenTypeKeyword) {
		if ip.stack.Len() == 0 {
			return ip.runtimeverr("failed to run step. exist command failed. stack is empty.\n", ip.stack.Len())
		}

		ip.runtimev("exist command.\n")
		vRes, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. exist command failed. failed to get path value: %v\n", err)
		}

		pRes, okRes := vRes.Path()
		if !okRes {
			return ip.runtimeverr("failed to run step. exist command failed. failed to get path.\n")
		}

		err = tools.ToolExistFile(pRes)
		if err != nil {
			ip.runtimev("failed to use exist tool: %v\n", err)
			err = ip.ipush(0)
			if err != nil {
				return false, fmt.Errorf("failed to run step. exist command failed. failure pushing value: %v", err)
			}
		} else {
			err = ip.ipush(1)
			if err != nil {
				return false, fmt.Errorf("failed to run step. exist command failed. failure pushing value: %v", err)
			}
		}
	} else if token.Equals("touch", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. touch command failed. stack is empty.\n", ip.stack.Len())
		}

		ip.runtimev("touch command.\n")
		vPath, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. touch command failed. failed to get path value: %v\n", err)
		}

		pPath, okPath := vPath.Path()
		if !okPath {
			return ip.runtimeverr("failed to run step. touch command failed. failed to get path.\n")
		}

		var result int = 1
		err = tools.ToolTouchFile(pPath)
		if err != nil {
			ip.runtimev("failed to use touch tool: %v\n", err)
			result = 0
		}

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. touch command failed. failure pushing value: %v", err)
		}
	} else if token.Equals("mkdir", types.TokenTypeKeyword) {
		if ip.stack.Len() == 0 {
			return ip.runtimeverr("failed to run step. mkdir command failed. stack is empty.\n", ip.stack.Len())
		}

		ip.runtimev("mkdir command.\n")
		vPath, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. mkdir command failed. failed to get path value: %v\n", err)
		}

		pPath, okPath := vPath.Path()
		if !okPath {
			return ip.runtimeverr("failed to run step. mkdir command failed. failed to get path.\n")
		}

		var result int = 1
		err = tools.ToolMakeDirectory(pPath)
		if err != nil {
			ip.runtimev("failed to use mkdir tool: %v\n", err)
			result = 0
		}

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. mkdir command failed. failure pushing value: %v", err)
		}
	} else if token.Equals("rm", types.TokenTypeKeyword) {
		if ip.stack.Len() == 0 {
			return ip.runtimeverr("failed to run step. rm command failed. stack is empty.\n", ip.stack.Len())
		}

		ip.runtimev("rm command.\n")
		vRes, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. rm command failed. failed to get path value: %v\n", err)
		}

		pRes, okRes := vRes.Path()
		if !okRes {
			return ip.runtimeverr("failed to run step. rm command failed. failed to get path.\n")
		}

		err = tools.ToolRemoveFile(pRes)
		if err != nil {
			ip.runtimev("failed to use rm tool: %v\n", err)
			err = ip.ipush(0)
			if err != nil {
				return false, fmt.Errorf("failed to run step. rm command failed. failure pushing value: %v", err)
			}
		} else {
			err = ip.ipush(1)
			if err != nil {
				return false, fmt.Errorf("failed to run step. rm command failed. failure pushing value: %v", err)
			}
		}
	} else if token.Equals("unzip", types.TokenTypeKeyword) {
		ip.runtimev("unzip command.\n")
		vDst, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. unzip command failed. failed to get dst value: %v\n", err)
		}

		vRes, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. unzip command failed. failed to get res value: %v\n", err)
		}

		pDst, okDst := vDst.Path()
		if !okDst {
			return ip.runtimeverr("failed to run step. unzip command failed. failed to get dst path.\n")
		}

		pRes, okRes := vRes.Path()
		if !okRes {
			return ip.runtimeverr("failed to run step. unzip command failed. failed to get res path.\n")
		}

		result, err := tools.ToolUnzipFile(pDst, pRes)
		if err != nil {
			ip.runtimev("failed to use unzip tool: %v\n", err)
			err = ip.ipush(-1)
			if err != nil {
				return false, fmt.Errorf("failed to run step. unzip command failed. failure pushing error value: %v", err)
			}
		} else {
			err = ip.ipush(int(result.DirCount))
			if err != nil {
				return false, fmt.Errorf("failed to run step. unzip command failed. failure pushing dir count: %v", err)
			}
			err = ip.ipush(int(result.FileCount))
			if err != nil {
				return false, fmt.Errorf("failed to run step. unzip command failed. failure pushing file count: %v", err)
			}
		}
	} else if token.Equals("lsf", types.TokenTypeKeyword) {
		ip.runtimev("lsf command.\n")
		vDir, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. lsf command failed. failed to get dir value: %v\n", err)
		}

		pDir, ok := vDir.Path()
		if !ok {
			return ip.runtimeverr("failed to run step. lsf command failed. failed to get dir path.\n")
		}

		var result int = 0
		count, err := tools.ToolLsf(pDir)
		if err != nil {
			ip.runtimev("failed to use lsf tool: %v\n", err)
		} else {
			result = int(count)
		}

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. lsf command failed. failure pushing count: %v", err)
		}
	} else if token.Equals("getf", types.TokenTypeKeyword) {
		ip.runtimev("getf command.\n")
		vDir, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. getf command failed. failed to get dir value: %v\n", err)
		}

		pDir, ok := vDir.Path()
		if !ok {
			return ip.runtimeverr("failed to run step. getf command failed. failed to get dir path.\n")
		}

		vIdx, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. getf command failed. failed to get idx value: %v\n", err)
		}

		idx, ok := vIdx.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. getf command failed. failed to get idx int.\n")
		}

		var result string = ""
		name, err := tools.ToolGetf(idx, pDir)
		if err != nil {
			ip.runtimev("failed to use getf tool: %v\n", err)
		} else {
			result = name
		}

		err = ip.spush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. getf command failed. failure pushing name: %v", err)
		}
	} else if token.Equals("lsd", types.TokenTypeKeyword) {
		ip.runtimev("lsd command.\n")
		vDir, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. lsd command failed. failed to get dir value: %v\n", err)
		}

		pDir, ok := vDir.Path()
		if !ok {
			return ip.runtimeverr("failed to run step. lsd command failed. failed to get dir path.\n")
		}

		var result int = 0
		count, err := tools.ToolLsd(pDir)
		if err != nil {
			ip.runtimev("failed to use lsd tool: %v\n", err)
		} else {
			result = int(count)
		}

		err = ip.ipush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. lsd command failed. failure pushing count: %v", err)
		}
	} else if token.Equals("getd", types.TokenTypeKeyword) {
		ip.runtimev("getd command.\n")
		vDir, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. getd command failed. failed to get dir value: %v\n", err)
		}

		pDir, ok := vDir.Path()
		if !ok {
			return ip.runtimeverr("failed to run step. getd command failed. failed to get dir path.\n")
		}

		vIdx, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. getd command failed. failed to get index value: %v\n", err)
		}

		iIdx, ok := vIdx.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. getd command failed. failed to get index int.\n")
		}

		result, err := tools.ToolGetd(iIdx, pDir)
		if err != nil {
			return ip.runtimeverr("failed to run step. getd command failed. %v\n", err)
		}

		err = ip.spush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. getd command failed. failure pushing result: %v", err)
		}
	} else if token.Equals("concat", types.TokenTypeKeyword) {
		ip.runtimev("concat command.\n")
		vB, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. concat command failed. failed to get second value: %v\n", err)
		}

		vA, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. concat command failed. failed to get first value: %v\n", err)
		}

		if pA, okA := vA.Path(); okA {
			if sB, okB := vB.String(); okB {
				result := tools.ToolConcatPath(pA, sB)
				err = ip.ppush(result)
				if err != nil {
					return false, fmt.Errorf("failed to run step. concat command failed. failure pushing path result: %v", err)
				}
			} else {
				return ip.runtimeverr("failed to run step. concat command failed. path concat requires string as second value.\n")
			}
		} else if sA, okA := vA.String(); okA {
			if sB, okB := vB.String(); okB {
				result := tools.ToolConcatString(sA, sB)
				err = ip.spush(result)
				if err != nil {
					return false, fmt.Errorf("failed to run step. concat command failed. failure pushing string result: %v", err)
				}
			} else {
				return ip.runtimeverr("failed to run step. concat command failed. string concat requires string as second value.\n")
			}
		} else {
			return ip.runtimeverr("failed to run step. concat command failed. first value must be path or string.\n")
		}
	} else if token.Equals("tostring", types.TokenTypeKeyword) {
		ip.runtimev("tostring command.\n")
		v, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. tostring command failed. failed to get value: %v\n", err)
		}

		var result string
		if i, ok := v.Int(); ok {
			result = tools.ToolToStringInt(i)
		} else if s, ok := v.String(); ok {
			result = tools.ToolToStringString(s)
		} else if p, ok := v.Path(); ok {
			result = tools.ToolToStringPath(p)
		} else {
			result = ""
		}

		err = ip.spush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. tostring command failed. failure pushing string result: %v", err)
		}
	} else if token.Equals("token", types.TokenTypeKeyword) {
		ip.runtimev("token command.\n")
		v, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. token command failed. failed to get value: %v\n", err)
		}

		var result string
		if s, ok := v.String(); ok {
			result = tools.ToolToToken(s)
		} else {
			result = ""
		}

		err = ip.ppush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. token command failed. failure pushing path result: %v", err)
		}
	} else if token.Equals("absolute", types.TokenTypeKeyword) {
		ip.runtimev("absolute command.\n")
		v, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. absolute command failed. failed to get value: %v\n", err)
		}

		var result string
		if s, ok := v.String(); ok {
			result = tools.ToolToAbsolute(s)
		} else {
			result = ""
		}

		err = ip.ppush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. absolute command failed. failure pushing path result: %v", err)
		}
	} else if token.Equals("relative", types.TokenTypeKeyword) {
		ip.runtimev("relative command.\n")
		v, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. relative command failed. failed to get value: %v\n", err)
		}

		var result string
		if s, ok := v.String(); ok {
			result = tools.ToolToRelative(s)
		} else {
			result = ""
		}

		err = ip.ppush(result)
		if err != nil {
			return false, fmt.Errorf("failed to run step. relative command failed. failure pushing path result: %v", err)
		}
	} else if token.Equals("exit", types.TokenTypeKeyword) {
		ip.runtimev("exit command. exiting...\n")
		ip.ip = ip.eop
		return true, nil
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

func (ip *Interpreter) load(name string) (StackValue, error) {
	var value StackValue
	value, ok := ip.memory[name]
	if !ok {
		return value, fmt.Errorf("failed to load(%s). value not found.", name)
	}

	left, ok := value.String()
	if ok {
		ip.runtimev("loaded \"%s\"\n", left)
		return value, nil
	}

	right, ok := value.Int()
	if ok {
		ip.runtimev("loaded %d\n", right)
		return value, nil
	}

	return value, fmt.Errorf("failed to load(%s). value found, but unknown type.", name)
}

func (ip *Interpreter) store(name string, value StackValue) error {
	left, okLeft := value.String()
	right, okRight := value.Int()
	if okLeft {
		ip.runtimev("stored \"%s\"\n", left)
	} else if okRight {
		ip.runtimev("stored %d\n", right)
	} else {
		return fmt.Errorf("failed to store(%s). value given, but unknown type.", name)
	}

	ip.memory[name] = value
	return nil
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

func (ip *Interpreter) push(value StackValue) error {
	if !value.IsAny() {
		return fmt.Errorf("failed to push stack. value given, but unknown type.")
	}

	ip.stack.Push(value)
	ip.runtimev("new stack size %d\n", ip.stack.Len())

	return nil
}

func (ip *Interpreter) ipush(item int) error {
	return ip.push(StackInt(item))
}

func (ip *Interpreter) spush(item string) error {
	return ip.push(StackString(item))
}

func (ip *Interpreter) ppush(item string) error {
	return ip.push(StackPath(item))
}

