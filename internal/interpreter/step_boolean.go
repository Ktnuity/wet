package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepBoolean(d *StepData) *StepResult {
	if d.token.Equals("&&", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. && operator failed. stack size is %d. 2 is reqired.", ip.stack.Len())
		}

		ip.runtimev("logicaland top stack values.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. && operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. && operator failed. failed to get second value: %v", err)
		}

		var truthy1 bool = false
		var truthy2 bool = false

		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"", s1)
			truthy1 = len(s1) > 0
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d", n1)
			truthy1 = n1 != 0
		}

		if s2, ok := v2.String(); ok {
			ip.runtimev("popped \"%s\"", s2)
			truthy2 = len(s2) > 0
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("popped %d", n2)
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
			return StepBad(fmt.Errorf("failed to run step. && operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %s", name)
	} else if d.token.Equals("||", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. || operator failed. stack size is %d. 2 is reqired.", ip.stack.Len())
		}

		ip.runtimev("logicalor top stack values.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. || operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. || operator failed. failed to get second value: %v", err)
		}

		var truthy1 bool = false
		var truthy2 bool = false

		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"", s1)
			truthy1 = len(s1) > 0
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d", n1)
			truthy1 = n1 != 0
		}

		if s2, ok := v2.String(); ok {
			ip.runtimev("popped \"%s\"", s2)
			truthy2 = len(s2) > 0
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("popped %d", n2)
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
			return StepBad(fmt.Errorf("failed to run step. || operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %s", name)
	} else {
		return nil
	}

	return StepOk()
}
