package interpreter

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepPrint(d *StepData) *StepResult {
	if d.token.Equals(".", types.TokenTypeSymbol) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. log operator failed. stack is empty.\n")
		}

		ip.runtimev("logging top value.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. log operator failed. failed to get value: %v\n", err)
		}

		n1, ok := v1.Int()
		if !ok {
			return ip.runtimeverr("failed to run step. log operator failed. failed to get int value.")
		}
		ip.runtimev("popped %d\n", n1)

		fmt.Printf("%d\n", n1)
	} else if d.token.Equals("puts",types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr("failed to run step. puts operator failed. stack is empty.\n")
		}

		ip.runtimev("printing top string value.\n")
		v1, err := ip.pop()
		if err != nil {
			return ip.runtimeverr("failed to run step. puts operator failed. failed to get value: %v\n", err)
		}

		s1, ok := v1.String()
		if !ok {
			return ip.runtimeverr("failed to run step. puts operator failed. value is not a string.\n")
		}
		ip.runtimev("popped \"%s\"\n", s1)

		fmt.Printf("%s", s1)
	} else if d.token.Equals(",", types.TokenTypeSymbol) {
		var result *int = nil
		for result == nil {
			ip.runtimev("reading value.\n")
			line, err := ip.reader.ReadString('\n')
			if err != nil && err != io.EOF {
				return ip.runtimeverr("Failed to run step. readint operator failed. failed to read input: %v\n", err)
			}

			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			i, err := strconv.ParseInt(line, 10, 32)
			if err != nil {
				continue
			}

			tmp := int(i)
			result = &tmp
		}

		ip.ipush(*result)
		ip.runtimev("pushed %d\n", *result)
	} else {
		return nil
	}

	return StepOk()
}
