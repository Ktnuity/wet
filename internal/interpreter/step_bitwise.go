package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepBitwise(d *StepData) *StepResult {
	if d.token.Equals("~", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. ~ operator failed. stack is empty.\n")
		}

		ip.runtimev("bitwisenot top stack value.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. ~ operator failed. failed to get value: %v\n", err)
		}

		if v1.IsString() {
			return ip.runtimeverr("failed to run step. ~ operator failed. cannot bitwise not string.\n")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. ~ operator failed. failed to get value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		result := ^n1

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. ~ operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d\n", result)
	} else if d.token.Equals("&", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. & operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("bitwiseand two numbers.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. & operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. & operator failed. failed to get second value: %v\n", err)
		}

		if v2.IsString() {
			return ip.runtimeverr("failed to run step. & operator failed. cannot bitwise and string.\n")
		}

		if v1.IsString() {
			return ip.runtimeverr("failed to run step. & operator failed. cannot bitwise and string.\n")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. & operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. & operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		result := n2 & n1

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. & operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d\n", result)
	} else if d.token.Equals("|", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. | operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("bitwiseor two numbers.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. | operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. | operator failed. failed to get second value: %v\n", err)
		}

		if v2.IsString() {
			return ip.runtimeverr("failed to run step. | operator failed. cannot bitwise or string.\n")
		}

		if v1.IsString() {
			return ip.runtimeverr("failed to run step. | operator failed. cannot bitwise or string.\n")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. | operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. | operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		result := n2 | n1

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. | operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d\n", result)
	} else if d.token.Equals("^", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr("failed to run step. ^ operator failed. stack size is %d. 2 is required.\n", ip.stack.Len())
		}

		ip.runtimev("bitwisexor two numbers.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. ^ operator failed. failed to get first value: %v\n", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. ^ operator failed. failed to get second value: %v\n", err)
		}

		if v2.IsString() {
			return ip.runtimeverr("failed to run step. ^ operator failed. cannot bitwise xor string.\n")
		}

		if v1.IsString() {
			return ip.runtimeverr("failed to run step. ^ operator failed. cannot bitwise xor string.\n")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. ^ operator failed. failed to get first value type.\n")
		}
		ip.runtimev("popped %d\n", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. ^ operator failed. failed to get second value type.\n")
		}
		ip.runtimev("popped %d\n", n2)

		result := n2 ^ n1

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. ^ operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d\n", result)
	} else {
		return nil
	}

	return StepOk()
}
