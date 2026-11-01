package interpreter

import (
	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepMemory(d *StepData) *StepResult {
	if d.token.Equals("store", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. store operator failed. stack is empty.")
		}

		ip.runtimev("storing value in memory.")
		vName, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. store operator failed. failed to get name: %v", err)
		}

		name, ok := vName.String()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. store operator failed. failed to get name value.")
		}

		vValue, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. store operator failed. failed to get value: %v", err)
		}

		ip.store(name, vValue)
	} else if d.token.Equals("iload", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr(d, "failed to run step. iload operator failed. stack is empty.")
		}

		ip.runtimev("loading value from memory.")
		vName, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. iload operator failed. failed to get name: %v", err)
		}

		name, ok := vName.String()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. iload operator failed. failed to get name value.")
		}

		value, err := ip.load(name)
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. iload operator failed. failed to load memory: %v", err)
		}

		ivalue, ok := value.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. iload operator failed. failed to load memory. int expected.")
		}

		ip.push(value)
		ip.runtimev("pushed %d", ivalue)
	} else if d.token.Equals("sload", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr(d, "failed to run step. sload operator failed. stack is empty.")
		}

		ip.runtimev("loading value from memory.")
		vName, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. sload operator failed. failed to get name: %v", err)
		}

		name, ok := vName.String()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. sload operator failed. failed to get name value.")
		}

		value, err := ip.load(name)
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. sload operator failed. failed to load memory: %v", err)
		}

		svalue, ok := value.String()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. sload operator failed. failed to load memory. string expected.")
		}

		ip.push(value)
		ip.runtimev("pushed \"%s\"", svalue)
	} else {
		return nil
	}

	return StepOk()
}
