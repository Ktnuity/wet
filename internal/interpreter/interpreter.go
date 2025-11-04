package interpreter

import (
	"fmt"
	"strings"

	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

type Interpreter struct {
	stack		util.Stack[int]
	program		[]Instruction
	ip			int
	eop			int
}

func CreateNew(tokens []types.Token) (*Interpreter, error) {
	var stack util.Stack[int] = util.Stack[int]{}

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
	} else if token.Equals(".", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. log operator failed. stack is empty.\n")
		}

		ip.runtimev("logging top value.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. log operator failed. failed to get value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		fmt.Printf("%d\n", n1)
	} else if token.Equals("+", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. + operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("adding two numbers.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. + operator failed. failed to get first number: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. + operator failed. failed to get second number: %v\n", err)
		}
		ip.runtimev("popped %d\n", n2)

		result := n1 + n2

		ip.ipush(result)
		ip.runtimev("pushed %d\n", result)
	} else if token.Equals("-", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. - operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("subtracting two numbers.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. - operator failed. failed to get first number: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. - operator failed. failed to get second number: %v\n", err)
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
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. * operator failed. failed to get first number: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. * operator failed. failed to get second number: %v\n", err)
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
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. / operator failed. failed to get first number: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. / operator failed. failed to get second number: %v\n", err)
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
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. % operator failed. failed to get first number: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. % operator failed. failed to get second number: %v\n", err)
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
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. dup operator failed. failed to get value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)
		ip.ipush(n1)
		ip.runtimev("pushed %d\n", n1)
		ip.ipush(n1)
		ip.runtimev("pushed %d\n", n1)
	} else if token.Equals("drop", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. drop operator failed. stack is empty.\n")
		}

		ip.runtimev("dropping top stack value.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. drop operator failed. failed to pop value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)
	} else if token.Equals("swap", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. swap operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("swapping top two stack values.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. swap operator failed. failed to get first value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. swap operator failed. failed to get second value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n2)

		ip.ipush(n1)
		ip.runtimev("pushed %d\n", n1)
		ip.ipush(n2)
		ip.runtimev("pushed %d\n", n2)
	} else if token.Equals("over", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. over operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("pushing 2nd top-most value onto stack.\n")
		n2, err := ip.ipeekOffset(1)
		if err != nil {
			return ip.runtimeverr("failed to run step. over operator failed. failed to get second value: %v\n", err)
		}
		ip.runtimev("peeked(1) %d\n", n2)
		ip.ipush(n2)
		ip.runtimev("pushed %d\n", n2)
	} else if token.Equals("2dup", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. 2dup operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("two-duping top stack values.\n")
		n1, err := ip.ipeekOffset(1)
		if err != nil {
			return ip.runtimeverr("failed to run step. 2dup operator failed. failed to get first value: %v", err)
		}
		ip.runtimev("peeked(1) %d\n", n1)

		n2, err := ip.ipeekOffset(0)
		if err != nil {
			return ip.runtimeverr("failed to run step. 2dup operator failed. failed to get second value: %v", err)
		}
		ip.runtimev("peeked(0) %d\n", n2)

		ip.ipush(n1)
		ip.runtimev("pushed %d\n", n1)
		ip.ipush(n2)
		ip.runtimev("pushed %d\n", n2)
	} else if token.Equals("2swap", types.TokenTypeKeyword) {
		if ip.stack.Len() < 4 {
			return ip.runtimeverr("failed to run step. 2swap operator failed. stack size is %d. 4 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("two-swapping top stack values.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. 2swap operator failed. failed to get first value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. 2swap operator failed. failed to get second value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n2)

		n3, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. 2swap operator failed. failed to get third value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n3)

		n4, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. 2swap operator failed. failed to get fourth value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n4)

		ip.ipush(n2)
		ip.runtimev("pushed %d\n", n2)
		ip.ipush(n1)
		ip.runtimev("pushed %d\n", n1)
		ip.ipush(n4)
		ip.runtimev("pushed %d\n", n4)
		ip.ipush(n3)
		ip.runtimev("pushed %d\n", n3)
	} else if token.Equals("if", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. if operator failed. stack is empty.\n")
		}

		ip.runtimev("validating if-condition.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. if operator failed. failed to get condition value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		if n1 == 0 {
			ip.ip = int(inst.Next)
			return true, nil
		}
	} else if token.Equals("else", types.TokenTypeKeyword) {
		if inst.Next != -1 {
			ip.ip = int(inst.Next)
			return true, nil
		}
	} else if token.Equals("end", types.TokenTypeKeyword) {
		if inst.Next != -1 {
			ip.ip = int(inst.Next)
			return true, nil
		}
	} else if token.Equals("=", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. = operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("equality-check top stack values.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. = operator failed. failed to get first value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. = operator failed. failed to get second value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n2)

		var result int = 0
		var name string = "false"
		if n1 == n2 {
			result = 1
			name = "true"
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals("!=", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. != operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("inequality-check top stack values.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. != operator failed. failed to get first value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. != operator failed. failed to get second value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n2)

		var result int = 0
		var name string = "false"
		if n1 != n2 {
			result = 1
			name = "true"
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals("<", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. < operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("less-check top stack values.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. < operator failed. failed to get first value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. < operator failed. failed to get second value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n2)

		var result int = 0
		var name string = "false"
		if n2 < n1 {
			result = 1
			name = "true"
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals(">", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. > operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("greater-check top stack values.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. > operator failed. failed to get first value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. > operator failed. failed to get second value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n2)

		var result int = 0
		var name string = "false"
		if n2 > n1 {
			result = 1
			name = "true"
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals("<=", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. <= operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("lessorequal-check top stack values.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. <= operator failed. failed to get first value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. <= operator failed. failed to get second value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n2)

		var result int = 0
		var name string = "false"
		if n2 <= n1 {
			result = 1
			name = "true"
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals(">=", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. >= operator failed. stack size is %d. 2 is reqired.\n", ip.stack.Len())
		}

		ip.runtimev("greaterorequal-check top stack values.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. >= operator failed. failed to get first value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. >= operator failed. failed to get second value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n2)

		var result int = 0
		var name string = "false"
		if n2 >= n1 {
			result = 1
			name = "true"
		}

		ip.ipush(result)
		ip.runtimev("pushed %s\n", name)
	} else if token.Equals("!", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. ! operator failed. stack is empty.\n")
		}

		ip.runtimev("logicalnot top stack value.\n")
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. ! operator failed. failed to get value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		var result int = 0
		var name string = "false"
		if n1 == 0 {
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
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. ~ operator failed. failed to get value: %v\n", err)
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
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. & operator failed. failed to get first number: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. & operator failed. failed to get second number: %v\n", err)
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
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. | operator failed. failed to get first number: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. | operator failed. failed to get second number: %v\n", err)
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
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. ^ operator failed. failed to get first number: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. ^ operator failed. failed to get second number: %v\n", err)
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
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. && operator failed. failed to get first value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. && operator failed. failed to get second value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n2)

		var result int = 0
		var name string = "false"
		if n1 != 0 && n2 != 0 {
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
		n1, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. || operator failed. failed to get first value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n1)

		n2, err := ip.ipop()
		if err != nil {
			return ip.runtimeverr("failed to run step. || operator failed. failed to get second value: %v\n", err)
		}
		ip.runtimev("popped %d\n", n2)

		var result int = 0
		var name string = "false"
		if n1 != 0 || n2 != 0 {
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

func (ip *Interpreter) ipop() (int, error) {
	result, ok := ip.stack.Pop()
	if !ok {
		return 0, fmt.Errorf("failed to ipop. stack is empty.")
	}

	ip.runtimev("new stack size %d\n", ip.stack.Len())

	return result, nil
}

func (ip *Interpreter) ipeek() (int, error) {
	return ip.ipeekOffset(0)
}

func (ip *Interpreter) ipeekOffset(offset int) (int, error) {
	idx := ip.stack.Len() - 1 - offset
	if idx < 0 {
		return 0, fmt.Errorf("failed to ipeek. idx(%d) < 0", idx)
	}
	if idx >= ip.stack.Len() {
		return 0, fmt.Errorf("failed to ipeek. idx(%d) >= len(%d)", idx, ip.stack.Len())
	}

	return ip.stack[ip.stack.Len() - 1 - offset], nil
}

func (ip *Interpreter) ipush(item int) {
	ip.stack.Push(item)
	ip.runtimev("new stack size %d\n", ip.stack.Len())
}


