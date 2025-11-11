package interpreter

import (
	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

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
	} else if d.token.Equals("ret", types.TokenTypeKeyword) {
		ip.runtimev("ret.\n")
		ret, ok := ip.calls.Pop()
		if !ok {
			return ip.runtimeverr("failed to run step. ret operator failed. call stack is empty.\n")
		}

		ip.ip = ret
		return StepOkSkip()
	} else if d.token.Equals("dret", types.TokenTypeKeyword) {
		ip.runtimev("dret.\n")
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. dret operator failed. stack is empty.\n")
		}

		ret, ok := ip.calls.Pop()
		if !ok {
			return ip.runtimeverr("failed to run step. dret operator failed. call stack is empty.\n")
		}

		pop1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. dret operator failed. failed to get count value: %v\n", err)
		}

		popi, ok := pop1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. dret operator failed. count value must be int.\n")
		}

		if ip.stack.Len() < popi {
			return ip.runtimeverr("failed to run step. dret operator failed. %d dret requires %d values on stack, had %d.\n", popi, popi + 1, ip.stack.Len() + 1)
		}

		for i := range popi {
			vp, err := ip.pop()
			if err != nil {
				return ip.runtimeverr("failed to run step. dret operator failed. failed to get value #%d: %v\n", i, err)
			}

			if pi, ok := vp.Int(); ok {
				ip.runtimev("popped %d\n", pi)
			} else if ps, ok := vp.String(); ok {
				ip.runtimev("popped \"%s\"\n", ps)
			} else if pp, ok := vp.Path(); ok {
				ip.runtimev("popped %s\n", pp)
			}
		}

		ip.ip = ret
		return StepOkSkip()
	} else if d.token.Equals("iret", types.TokenTypeKeyword) {
		ip.runtimev("iret.\n")
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. iret operator failed. stack need 2 elements .\n")
		}

		ret, ok := ip.calls.Pop()
		if !ok {
			return ip.runtimeverr("failed to run step. iret operator failed. call stack is empty.\n")
		}

		ret1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. iret operator failed. failed to get ret count value: %v\n", err)
		}

		pop1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. iret operator failed. failed to get drop count value: %v\n", err)
		}

		reti, ok := ret1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. iret operator failed. ret count value must be int.\n")
		}

		popi, ok := pop1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. iret operator failed. drop count value must be int.\n")
		}

		if ip.stack.Len() < popi + reti {
			return ip.runtimeverr("failed to run step. iret operator failed. %d %d iret requires %d values on stack, had %d.\n", popi, reti, popi + reti + 2, ip.stack.Len() + 2)
		}

		retStack := util.Stack[types.DataType]{}

		for i := range reti {
			retp, err := ip.pop()
			if err != nil {
				return ip.runtimeverr("failed to run step. iret operator failed. failed to get ret value #%d: %v\n", i, err)
			}

			if pi, ok := retp.Int(); ok {
				ip.runtimev("popped %d\n", pi)
			} else if ps, ok := retp.String(); ok {
				ip.runtimev("popped \"%s\"\n", ps)
			} else if pp, ok := retp.Path(); ok {
				ip.runtimev("popped %s\n", pp)
			}

			retStack.Push(retp)
		}

		for i := range popi {
			vp, err := ip.pop()
			if err != nil {
				return ip.runtimeverr("failed to run step. iret operator failed. failed to get drop value #%d: %v\n", i, err)
			}

			if pi, ok := vp.Int(); ok {
				ip.runtimev("dropped %d\n", pi)
			} else if ps, ok := vp.String(); ok {
				ip.runtimev("dropped \"%s\"\n", ps)
			} else if pp, ok := vp.Path(); ok {
				ip.runtimev("dropped %s\n", pp)
			}
		}

		for retStack.Len() > 0 {
			retp, ok := retStack.Pop()
			if !ok {
				return ip.runtimeverr("failed to run step. iret operator failed. failed to repush ret value #%d\n", retStack.Len())
			}

			if pi, ok := retp.Int(); ok {
				ip.runtimev("pushed %d\n", pi)
			} else if ps, ok := retp.String(); ok {
				ip.runtimev("pushed \"%s\"\n", ps)
			} else if pp, ok := retp.Path(); ok {
				ip.runtimev("pushed %s\n", pp)
			}

			ip.push(retp)
		}

		ip.ip = ret
		return StepOkSkip()
	} else {
		return nil
	}

	return StepOk()
}
