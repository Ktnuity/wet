package interpreter

import (
	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepMemory(d *StepData) *StepResult {
	if d.token.Equals("store", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. store operator failed. stack is empty.\n")
		}

		ip.runtimev("storing value in memory.\n")
		vName, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. store operator failed. failed to get name: %v\n", err)
		}

		name, ok := vName.String()
		if !ok {
			return ip.runtimeverr("failed to run step. store operator failed. failed to get name value.\n")
		}

		vValue, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. store operator failed. failed to get value: %v\n", err)
		}

		ip.store(name, vValue)
	} else if d.token.Equals("iload", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. iload operator failed. stack is empty.\n")
		}

		ip.runtimev("loading value from memory.\n")
		vName, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. iload operator failed. failed to get name: %v\n", err)
		}

		name, ok := vName.String()
		if !ok {
			return ip.runtimeverr("failed to run step. iload operator failed. failed to get name value.\n")
		}

		value, err := ip.load(name)
		if err != nil {
			return ip.runtimeverr("failed to run step. iload operator failed. failed to load memory: %v\n", err)
		}

		ivalue, ok := value.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. iload operator failed. failed to load memory. int expected.\n")
		}

		ip.push(value)
		ip.runtimev("pushed %d\n", ivalue)
	} else if d.token.Equals("sload", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. sload operator failed. stack is empty.\n")
		}

		ip.runtimev("loading value from memory.\n")
		vName, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. sload operator failed. failed to get name: %v\n", err)
		}

		name, ok := vName.String()
		if !ok {
			return ip.runtimeverr("failed to run step. sload operator failed. failed to get name value.\n")
		}

		value, err := ip.load(name)
		if err != nil {
			return ip.runtimeverr("failed to run step. sload operator failed. failed to load memory: %v\n", err)
		}

		svalue, ok := value.String()
		if !ok {
			return ip.runtimeverr("failed to run step. sload operator failed. failed to load memory. string expected.\n")
		}

		ip.push(value)
		ip.runtimev("pushed \"%s\"\n", svalue)
	} else {
		return nil
	}

	return StepOk()
}
