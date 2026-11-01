package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepBitwise(d *StepData) *StepResult {
	if d.token.Equals("~", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr(d, "failed to run step. ~ operator failed. stack is empty.")
		}

		ip.runtimev("bitwisenot top stack value.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. ~ operator failed. failed to get value: %v", err)
		}

		if v1.IsString() {
			return ip.runtimeverr(d, "failed to run step. ~ operator failed. cannot bitwise not string.")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. ~ operator failed. failed to get value type.")
		}
		ip.runtimev("popped %d", n1)

		result := ^n1

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. ~ operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d", result)
	} else if d.token.Equals("&", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. & operator failed. stack size is %d. 2 is required.", ip.stack.Len())
		}

		ip.runtimev("bitwiseand two numbers.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. & operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. & operator failed. failed to get second value: %v", err)
		}

		if v2.IsString() {
			return ip.runtimeverr(d, "failed to run step. & operator failed. cannot bitwise and string.")
		}

		if v1.IsString() {
			return ip.runtimeverr(d, "failed to run step. & operator failed. cannot bitwise and string.")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. & operator failed. failed to get first value type.")
		}
		ip.runtimev("popped %d", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. & operator failed. failed to get second value type.")
		}
		ip.runtimev("popped %d", n2)

		result := n2 & n1

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. & operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d", result)
	} else if d.token.Equals("|", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. | operator failed. stack size is %d. 2 is required.", ip.stack.Len())
		}

		ip.runtimev("bitwiseor two numbers.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. | operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. | operator failed. failed to get second value: %v", err)
		}

		if v2.IsString() {
			return ip.runtimeverr(d, "failed to run step. | operator failed. cannot bitwise or string.")
		}

		if v1.IsString() {
			return ip.runtimeverr(d, "failed to run step. | operator failed. cannot bitwise or string.")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. | operator failed. failed to get first value type.")
		}
		ip.runtimev("popped %d", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. | operator failed. failed to get second value type.")
		}
		ip.runtimev("popped %d", n2)

		result := n2 | n1

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. | operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d", result)
	} else if d.token.Equals("^", types.TokenTypeSymbol) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. ^ operator failed. stack size is %d. 2 is required.", ip.stack.Len())
		}

		ip.runtimev("bitwisexor two numbers.")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. ^ operator failed. failed to get first value: %v", err)
		}

		v2, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. ^ operator failed. failed to get second value: %v", err)
		}

		if v2.IsString() {
			return ip.runtimeverr(d, "failed to run step. ^ operator failed. cannot bitwise xor string.")
		}

		if v1.IsString() {
			return ip.runtimeverr(d, "failed to run step. ^ operator failed. cannot bitwise xor string.")
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. ^ operator failed. failed to get first value type.")
		}
		ip.runtimev("popped %d", n1)

		n2, ok := v2.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. ^ operator failed. failed to get second value type.")
		}
		ip.runtimev("popped %d", n2)

		result := n2 ^ n1

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. ^ operator failed. failure pushing value: %v", err))
		}
		ip.runtimev("pushed %d", result)
	} else {
		return nil
	}

	return StepOk()
}
