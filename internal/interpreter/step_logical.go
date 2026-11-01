package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepLogical(d *StepData) *StepResult {
	if d.token.Equals("=", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. = operator failed. stack size is %d. 2 is reqired.", ip.stack.Len())
		}

		ip.runtimev("equality-check top stack values.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. = operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. = operator failed. failed to get second value: %v", err)
		}

		var result int = 0
		var name string = "false"

		if s1, ok := v1.String(); ok {
			if s2, ok := v2.String(); ok {
				ip.runtimev("popped \"%s\"", s1)
				ip.runtimev("popped \"%s\"", s2)
				if s1 == s2 {
					result = 1
					name = "true"
				}
			} else {
				ip.runtimev("popped int")
				ip.runtimev("popped \"%s\"", s1)
			}
		} else if n1, ok := v1.Int(); ok {
			if n2, ok := v2.Int(); ok {
				ip.runtimev("popped %d", n1)
				ip.runtimev("popped %d", n2)
				if n1 == n2 {
					result = 1
					name = "true"
				}
			} else {
				ip.runtimev("popped string")
				ip.runtimev("popped %d", n1)
			}
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. = operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %s", name)
	} else if d.token.Equals("!=", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. != operator failed. stack size is %d. 2 is reqired.", ip.stack.Len())
		}

		ip.runtimev("inequality-check top stack values.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. != operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. != operator failed. failed to get second value: %v", err)
		}

		var result int = 1
		var name string = "true"

		if s1, ok := v1.String(); ok {
			if s2, ok := v2.String(); ok {
				ip.runtimev("popped \"%s\"", s1)
				ip.runtimev("popped \"%s\"", s2)
				if s1 == s2 {
					result = 0
					name = "false"
				}
			} else {
				ip.runtimev("popped int")
				ip.runtimev("popped \"%s\"", s1)
			}
		} else if n1, ok := v1.Int(); ok {
			if n2, ok := v2.Int(); ok {
				ip.runtimev("popped %d", n1)
				ip.runtimev("popped %d", n2)
				if n1 == n2 {
					result = 0
					name = "false"
				}
			} else {
				ip.runtimev("popped string")
				ip.runtimev("popped %d", n1)
			}
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. != operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %s", name)
	} else if d.token.Equals("<", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. < operator failed. stack size is %d. 2 is reqired.", ip.stack.Len())
		}

		ip.runtimev("less-check top stack values.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. < operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. < operator failed. failed to get second value: %v", err)
		}

		var result int = 0
		var name string = "false"

		if s1, ok := v1.String(); ok {
			if s2, ok := v2.String(); ok {
				ip.runtimev("popped \"%s\"", s1)
				ip.runtimev("popped \"%s\"", s2)
				if s2 < s1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr(d, "failed to run step. < operator failed. cannot compare string and int.")
			}
		} else if n1, ok := v1.Int(); ok {
			if n2, ok := v2.Int(); ok {
				ip.runtimev("popped %d", n1)
				ip.runtimev("popped %d", n2)
				if n2 < n1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr(d, "failed to run step. < operator failed. cannot compare int and string.")
			}
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. < operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %s", name)
	} else if d.token.Equals(">", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. > operator failed. stack size is %d. 2 is reqired.", ip.stack.Len())
		}

		ip.runtimev("greater-check top stack values.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. > operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. > operator failed. failed to get second value: %v", err)
		}

		var result int = 0
		var name string = "false"

		if s1, ok := v1.String(); ok {
			if s2, ok := v2.String(); ok {
				ip.runtimev("popped \"%s\"", s1)
				ip.runtimev("popped \"%s\"", s2)
				if s2 > s1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr(d, "failed to run step. > operator failed. cannot compare string and int.")
			}
		} else if n1, ok := v1.Int(); ok {
			if n2, ok := v2.Int(); ok {
				ip.runtimev("popped %d", n1)
				ip.runtimev("popped %d", n2)
				if n2 > n1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr(d, "failed to run step. > operator failed. cannot compare int and string.")
			}
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. > operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %s", name)
	} else if d.token.Equals("<=", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. <= operator failed. stack size is %d. 2 is reqired.", ip.stack.Len())
		}

		ip.runtimev("lessorequal-check top stack values.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. <= operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. <= operator failed. failed to get second value: %v", err)
		}

		var result int = 0
		var name string = "false"

		if s1, ok := v1.String(); ok {
			if s2, ok := v2.String(); ok {
				ip.runtimev("popped \"%s\"", s1)
				ip.runtimev("popped \"%s\"", s2)
				if s2 <= s1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr(d, "failed to run step. <= operator failed. cannot compare string and int.")
			}
		} else if n1, ok := v1.Int(); ok {
			if n2, ok := v2.Int(); ok {
				ip.runtimev("popped %d", n1)
				ip.runtimev("popped %d", n2)
				if n2 <= n1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr(d, "failed to run step. <= operator failed. cannot compare int and string.")
			}
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. <= operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %s", name)
	} else if d.token.Equals(">=", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. >= operator failed. stack size is %d. 2 is reqired.", ip.stack.Len())
		}

		ip.runtimev("greaterorequal-check top stack values.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. >= operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. >= operator failed. failed to get second value: %v", err)
		}

		var result int = 0
		var name string = "false"

		if s1, ok := v1.String(); ok {
			if s2, ok := v2.String(); ok {
				ip.runtimev("popped \"%s\"", s1)
				ip.runtimev("popped \"%s\"", s2)
				if s2 >= s1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr(d, "failed to run step. >= operator failed. cannot compare string and int.")
			}
		} else if n1, ok := v1.Int(); ok {
			if n2, ok := v2.Int(); ok {
				ip.runtimev("popped %d", n1)
				ip.runtimev("popped %d", n2)
				if n2 >= n1 {
					result = 1
					name = "true"
				}
			} else {
				return ip.runtimeverr(d, "failed to run step. >= operator failed. cannot compare int and string.")
			}
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. >= operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %s", name)
	} else if d.token.Equals("!", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr(d, "failed to run step. ! operator failed. stack is empty.")
		}

		ip.runtimev("logicalnot top stack value.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. ! operator failed. failed to get value: %v", err)
		}

		var truthy bool = false
		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"", s1)
			truthy = len(s1) > 0
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d", n1)
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
			return StepBad(fmt.Errorf("failed to run step. ! operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %s", name)
	} else {
		return nil
	}

	return StepOk()
}
