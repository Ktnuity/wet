package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepStack(d *StepData) *StepResult {
	if d.token.Equals("dup", types.TokenTypeKeyword) {
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
				return StepBad(fmt.Errorf("failed to run step. dup operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed \"%s\"\n", s1)
			err = ip.spush(s1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. dup operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed \"%s\"\n", s1)
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d\n", n1)
			err := ip.ipush(n1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. dup operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %d\n", n1)
			err = ip.ipush(n1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. dup operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %d\n", n1)
		} else if p1, ok := v1.Path(); ok {
			ip.runtimev("popped %s\n", p1)
			err := ip.ppush(p1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. dup operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %s\n", p1)
			err = ip.ppush(p1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. dup operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %s\n", p1)
		}
	} else if d.token.Equals("drop", types.TokenTypeKeyword) {
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
	} else if d.token.Equals("swap", types.TokenTypeKeyword) {
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
			return StepBad(fmt.Errorf("failed to run step. swap operator failed. failure pushing value: %v", err))
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
			return StepBad(fmt.Errorf("failed to run step. swap operator failed. failure pushing value: %v", err))
		}
		if s2, ok := v2.String(); ok {
			ip.runtimev("pushed \"%s\"\n", s2)
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("pushed %d\n", n2)
		} else if p2, ok := v2.Path(); ok {
			ip.runtimev("pushed %s\n", p2)
		}
	} else if d.token.Equals("over", types.TokenTypeKeyword) {
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
				return StepBad(fmt.Errorf("failed to run step. over operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed \"%s\"\n", s2)
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("peeked(1) %d\n", n2)
			err := ip.ipush(n2)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. over operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %d\n", n2)
		} else if p2, ok := v2.Path(); ok {
			ip.runtimev("peeked(1) %d\n", p2)
			err := ip.ppush(p2)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. over operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %s\n", p2)
		}
	} else if d.token.Equals("2dup", types.TokenTypeKeyword) {
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
			return StepBad(fmt.Errorf("failed to run step. 2dup operator failed. failure pushing value: %v", err))
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
			return StepBad(fmt.Errorf("failed to run step. 2dup operator failed. failure pushing value: %v", err))
		}
		if s2, ok := v2.String(); ok {
			ip.runtimev("pushed \"%s\"\n", s2)
		} else if n2, ok := v2.Int(); ok {
			ip.runtimev("pushed %d\n", n2)
		} else if p2, ok := v2.Path(); ok {
			ip.runtimev("pushed %s\n", p2)
		}
	} else if d.token.Equals("2swap", types.TokenTypeKeyword) {
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
			return StepBad(fmt.Errorf("failed to run step. 2swap operator failed. failure pushing value: %v", err))
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
			return StepBad(fmt.Errorf("failed to run step. 2swap operator failed. failure pushing value: %v", err))
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
			return StepBad(fmt.Errorf("failed to run step. 2swap operator failed. failure pushing value: %v", err))
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
			return StepBad(fmt.Errorf("failed to run step. 2swap operator failed. failure pushing value: %v", err))
		}
		if s3, ok := v3.String(); ok {
			ip.runtimev("pushed \"%s\"\n", s3)
		} else if n3, ok := v3.Int(); ok {
			ip.runtimev("pushed %d\n", n3)
		} else if p3, ok := v3.Path(); ok {
			ip.runtimev("pushed %s\n", p3)
		}
	} else {
		return nil
	}

	return StepOk()
}
