package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepArithmetic(d *StepData) *StepResult {
	if d.token.Equals("+", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. + operator failed. stack size is %d. 2 is required.", ip.stack.Len())
		}

		ip.runtimev("adding two numbers, concating two strings, or concating string + number.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. + operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. + operator failed. failed to get second value: %v", err)
		}

		if s2, ok := v2.String(); ok {
			if s1, ok := v1.String(); ok {
				ip.runtimev("popped \"%s\"", s1)
				ip.runtimev("popped \"%s\"", s2)

				result := s2 + s1

				err := ip.spush(result)
				if err != nil {
					return StepBad(fmt.Errorf("failed to run step. + operator failed. failure pushing value: %v", err))
				}
				ip.runtimev("pushed \"%s\"", result)
			} else if n1, ok := v1.Int(); ok {
				ip.runtimev("popped %d", n1)
				ip.runtimev("popped \"%s\"", s2)

				result := fmt.Sprintf("%s%d", s2, n1)

				err := ip.spush(result)
				if err != nil {
					return StepBad(fmt.Errorf("failed to run step. + operator failed. failure pushing value: %v", err))
				}
				ip.runtimev("pushed \"%s\"", result)
			} else {
				return ip.runtimeverr(d, "failed to run step. + operator failed. failed to get first value type.")
			}
		} else if n2, ok := v2.Int(); ok {
			if n1, ok := v1.Int(); ok {
				ip.runtimev("popped %d", n1)
				ip.runtimev("popped %d", n2)

				result := n2 + n1

				err := ip.ipush(result)
				if err != nil {
					return StepBad(fmt.Errorf("failed to run step. + operator failed. failure pushing value: %v", err))
				}
				ip.runtimev("pushed %d", result)
			} else if v1.IsString() {
				return ip.runtimeverr(d, "failed to run step. + operator failed. second value is int, thus first value must be int.")
			} else {
				return ip.runtimeverr(d, "failed to run step. + operator failed. failed to get first value type.")
			}
		} else {
			return ip.runtimeverr(d, "failed to run step. + operator failed. failed to get second value type.")
		}
	} else if d.token.Equals("-", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. - operator failed. stack size is %d. 2 is required.", ip.stack.Len())
		}

		ip.runtimev("subtracting two numbers.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. - operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. - operator failed. failed to get second value: %v", err)
		}

		if v2.IsString() {
			return ip.runtimeverr(d, "failed to run step. - operator failed. cannot subtract from string.")
		}

		if v1.IsString() {
			return ip.runtimeverr(d, "failed to run step. - operator failed. cannot subtract string.")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. - operator failed. failed to get first value type.")
		}
		ip.runtimev("popped %d", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. - operator failed. failed to get second value type.")
		}
		ip.runtimev("popped %d", n2)

		result := n2 - n1

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. - operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d", result)
	} else if d.token.Equals("*", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. * operator failed. stack size is %d. 2 is required.", ip.stack.Len())
		}

		ip.runtimev("multiplying two numbers.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. * operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. * operator failed. failed to get second value: %v", err)
		}

		if v2.IsString() {
			return ip.runtimeverr(d, "failed to run step. * operator failed. cannot multiply string.")
		}

		if v1.IsString() {
			return ip.runtimeverr(d, "failed to run step. * operator failed. cannot multiply string.")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. * operator failed. failed to get first value type.")
		}
		ip.runtimev("popped %d", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. * operator failed. failed to get second value type.")
		}
		ip.runtimev("popped %d", n2)

		result := n1 * n2

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. * operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d", result)
	} else if d.token.Equals("/", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. / operator failed. stack size is %d. 2 is required.", ip.stack.Len())
		}

		ip.runtimev("dividing two numbers.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. / operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. / operator failed. failed to get second value: %v", err)
		}

		if v2.IsString() {
			return ip.runtimeverr(d, "failed to run step. / operator failed. cannot divide string.")
		}

		if v1.IsString() {
			return ip.runtimeverr(d, "failed to run step. / operator failed. cannot divide string.")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. / operator failed. failed to get first value type.")
		}
		ip.runtimev("popped %d", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. / operator failed. failed to get second value type.")
		}
		ip.runtimev("popped %d", n2)

		var result int
		if n1 != 0 {
			result = n2 / n1
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. / operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d", result)
	} else if d.token.Equals("%", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. %% operator failed. stack size is %d. 2 is required.", ip.stack.Len())
		}

		ip.runtimev("modulo two numbers.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. %% operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. %% operator failed. failed to get second value: %v", err)
		}

		if v2.IsString() {
			return ip.runtimeverr(d, "failed to run step. %% operator failed. cannot modulo string.")
		}

		if v1.IsString() {
			return ip.runtimeverr(d, "failed to run step. %% operator failed. cannot modulo string.")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. %% operator failed. failed to get first value type.")
		}
		ip.runtimev("popped %d", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. %% operator failed. failed to get second value type.")
		}
		ip.runtimev("popped %d", n2)

		var result int
		if n1 != 0 {
			result = n2 % n1
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. %% operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d", result)
	} else if d.token.Equals("++", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr(d, "failed to run step. ++ operator failed. stack is empty.")
		}

		ip.runtimev("increment top stack value.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. ++ operator failed. failed to get value: %v", err)
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. ++ operator failed. failed tog et value type.")
		}
		ip.runtimev("popped %d", n1)

		result := n1 + 1

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. ++ operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d", result)
	} else if d.token.Equals("--", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr(d, "failed to run step. -- operator failed. stack is empty.")
		}

		ip.runtimev("decrement top stack value.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. -- operator failed. failed to get value: %v", err)
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. -- operator failed. failed tog et value type.")
		}
		ip.runtimev("popped %d", n1)

		result := n1 - 1

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. -- operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d", result)
	} else {
		return nil
	}

	return StepOk()
}
