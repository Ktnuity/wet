package interpreter

import (
	"fmt"
	"strconv"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepPrimary(d *StepData) *StepResult {
	if d.token.Equals("int", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr(d, "failed to run step. int operator failed. stack is empty.")
		}

		ip.runtimev("parsing string to int.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. int operator failed. failed to get value: %v", err)
		}

		s1, ok := v1.String()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. int operator failed. value is not a string.")
		}
		ip.runtimev("popped \"%s\"", s1)

		parsed, err := strconv.Atoi(s1)
		if err != nil {
			ip.runtimev("parse failed, pushing 0")
			err := ip.ipush(0)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. int operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %d", 0)
			err = ip.ipush(0)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. int operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %d (false)", 0)
		} else {
			err := ip.ipush(parsed)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. int operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %d", parsed)
			err = ip.ipush(1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. int operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed %d (true)", 1)
		}
	} else if d.token.Equals("string", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr(d, "failed to run step. string operator failed. stack is empty.")
		}

		ip.runtimev("converting value to string.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. string operator failed. failed to get value: %v", err)
		}

		if s1, ok := v1.String(); ok {
			ip.runtimev("popped \"%s\"", s1)
			err := ip.spush(s1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. string operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed \"%s\"", s1)
		} else if n1, ok := v1.Int(); ok {
			ip.runtimev("popped %d", n1)
			result := strconv.Itoa(n1)
			err := ip.spush(result)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. string operator failed. failure pushing value: %v", err))
			}
			ip.runtimev("pushed \"%s\"", result)
		}
	} else {
		return nil
	}

	return StepOk()
}
