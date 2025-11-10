package interpreter

import (
	"fmt"
	"strconv"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepPrimary(d *StepData) *StepResult {
	if d.token.Equals("int", types.TokenTypeKeyword) {
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
				return StepBad(fmt.Errorf("failed to run step. int operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %d\n", 0)
			err = ip.ipush(0)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. int operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %d (false)\n", 0)
		} else {
			err := ip.ipush(parsed)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. int operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %d\n", parsed)
			err = ip.ipush(1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. int operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %d (true)\n", 1)
		}
	} else if d.token.Equals("string", types.TokenTypeKeyword) {
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
				return StepBad(fmt.Errorf("failed to run step. string operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed \"%s\"\n", s1)
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d\n", n1)
			result := strconv.Itoa(n1)
			err := ip.spush(result)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. string operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed \"%s\"\n", result)
		}
	} else {
		return nil
	}

	return StepOk()
}
