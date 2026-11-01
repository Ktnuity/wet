package interpreter

import (
	"fmt"
	"strings"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepLiteral(d *StepData) *StepResult {
	if d.token.Equals("", types.TokenTypeNumber) {
		num, ok := d.token.GetNumberValue()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. token at index %d is has invalid number", ip.ip)
		}

		ip.runtimev("pushing %d", num)
		err := ip.ipush(num)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. number push operator failed. failure pushing value: %v", err))
		}
	} else if d.token.Equals("", types.TokenTypeString) {
		str, ok := d.token.GetStringValue()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. token at index %d has invalid string", ip.ip)
		}

		ip.runtimev("pushing \"%s\"", escapeNewLines(str))
		err := ip.spush(str)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. string push operator failed. failure pushing value: %v", err))
		}
	} else if d.token.Equals("", types.TokenTypePath) {
		str, ok := d.token.GetPathValue()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. token at index %d has invalid path", ip.ip)
		}

		ip.runtimev("pushing path(%s)", str)
		err := ip.ppush(str)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. path push operator failed. failure pushing value: %v", err))
		}
	} else {
		return nil
	}

	return StepOk()
}

func escapeNewLines(str string) string {
	return strings.ReplaceAll(str, "\n", "\\n")
}
