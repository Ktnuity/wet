package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepBoolean(d *StepData) *StepResult {
	if d.token.Equals("&&", types.TokenTypeSymbol) {
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
			return StepBad(fmt.Errorf("failed to run step. && operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %s\n", name)
	} else if d.token.Equals("||", types.TokenTypeSymbol) {
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
			return StepBad(fmt.Errorf("failed to run step. || operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %s\n", name)
	} else {
		return nil
	}

	return StepOk()
}
