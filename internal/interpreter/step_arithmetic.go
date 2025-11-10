package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepArithmetic(d *StepData) *StepResult {
	if d.token.Equals("+", types.TokenTypeSymbol) {
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
					return StepBad(fmt.Errorf("failed to run step. + operator failed. failure pushing value: %v", err))
				}
				ip.runtimev("pushed \"%s\"\n", result)
			} else if n1, ok := v1.Int(); ok {
				ip.runtimev("popped %d\n", n1)
				ip.runtimev("popped \"%s\"\n", s2)

				result := fmt.Sprintf("%s%d", s2, n1)

				err := ip.spush(result)
				if err != nil {
					return StepBad(fmt.Errorf("failed to run step. + operator failed. failure pushing value: %v", err))
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
					return StepBad(fmt.Errorf("failed to run step. + operator failed. failure pushing value: %v", err))
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
	} else if d.token.Equals("-", types.TokenTypeSymbol) {
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
			return StepBad(fmt.Errorf("failed to run step. - operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d\n", result)
	} else if d.token.Equals("*", types.TokenTypeSymbol) {
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
			return StepBad(fmt.Errorf("failed to run step. * operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d\n", result)
	} else if d.token.Equals("/", types.TokenTypeSymbol) {
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
			return StepBad(fmt.Errorf("failed to run step. / operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d\n", result)
	} else if d.token.Equals("%", types.TokenTypeSymbol) {
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
			return StepBad(fmt.Errorf("failed to run step. %% operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d\n", result)
	} else if d.token.Equals("++", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. ++ operator failed. stack is empty.\n")
		}

		ip.runtimev("increment top stack value.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. ++ operator failed. failed to get value: %v\n", err)
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. ++ operator failed. failed tog et value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		result := n1 + 1

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. ++ operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d\n", result)
	} else if d.token.Equals("--", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. -- operator failed. stack is empty.\n")
		}

		ip.runtimev("decrement top stack value.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. -- operator failed. failed to get value: %v\n", err)
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. -- operator failed. failed tog et value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		result := n1 - 1

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. -- operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d\n", result)
	} else {
		return nil
	}

	return StepOk()
}
