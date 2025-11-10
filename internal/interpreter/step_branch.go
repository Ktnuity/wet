package interpreter

import "github.com/ktnuity/wet/internal/types"

func (ip *Interpreter) StepBranch(d *StepData) *StepResult {
	if d.token.Equals("if", types.TokenTypeKeyword) {
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
			ip.ip = int(d.inst.Next)
			return StepOkSkip()
		}
	} else if d.token.Equals("unless", types.TokenTypeKeyword) {
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
			ip.ip = int(d.inst.Next)
			return StepOkSkip()
		}
	} else if d.token.Equals("else", types.TokenTypeKeyword) {
		if d.inst.Next != -1 {
			ip.ip = int(d.inst.Next)
			return StepOkSkip()
		}
	} else if d.token.Equals("do", types.TokenTypeKeyword) {
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

		if d.inst.Mode == DoModeUntil {
			truthy = !truthy
		}

		if !truthy {
			ip.ip = int(d.inst.Next)
			return StepOkSkip()
		}
	} else if d.token.Equals("end", types.TokenTypeKeyword) {
		if d.inst.Mode == EndModeProc {
			ret, ok := ip.calls.Pop()
			if !ok {
				return ip.runtimeverr("failed to run step. end operator failed. call stack is empty.\n")
			}

			ip.ip = ret
			return StepOkSkip()
		} else {
			if d.inst.Next != -1 {
				ip.ip = int(d.inst.Next)
				return StepOkSkip()
			}
		}
	} else if d.token.Equals("while", types.TokenTypeKeyword) {
		ip.runtimev("while (do nothing).\n")
	} else if d.token.Equals("until", types.TokenTypeKeyword) {
		ip.runtimev("until (do nothing).\n")
	} else if d.token.Equals("proc", types.TokenTypeKeyword) {
		ip.runtimev("proc (skip to end).\n")
		if d.inst.Next == -1 {
			return ip.runtimeverr("failed to run step. proc operator failed. skip is undefined.\n")
		}

		ip.ip = int(d.inst.Next)
		return StepOkSkip()
	} else {
		return nil
	}

	return StepOk()
}
