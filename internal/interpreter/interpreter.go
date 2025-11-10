package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) Run() (bool, error) {
	if ip == nil {
		return false, fmt.Errorf("failed to run interpreter. instance is nil.")
	}

	for ip.ip < ip.eop {
		status, err := ip.Step()
		if err != nil {
			return false, fmt.Errorf("failed to run interpreter. interpreter step failed: %v", err)
		}

		if !status {
			return false, nil
		}
	}

	return true, nil
}

func (ip *Interpreter) Step() (bool, error) {
	data, result := ip.newStep()
	if result != nil {
		return result.status, result.result
	}

	ip.runtimev("current stack size %d\n", ip.stack.Len())

	if IsVerboseRuntime() {
		tokenLog := data.token.Format()
		ip.runtimev("current token [%s] at ip %d\n", tokenLog, ip.ip)
	}

	result = ip.StepLiteral(&data)
	if result != nil {
		return ip.defaultStep(*result)
	}

	result = ip.StepPrint(&data)
	if result != nil {
		return ip.defaultStep(*result)
	}

	result = ip.StepPrimary(&data)
	if result != nil {
		return ip.defaultStep(*result)
	}

	result = ip.StepMemory(&data)
	if result != nil {
		return ip.defaultStep(*result)
	}

	result = ip.StepArithmetic(&data)
	if result != nil {
		return ip.defaultStep(*result)
	}

	result = ip.StepStack(&data)
	if result != nil {
		return ip.defaultStep(*result)
	}

	result = ip.StepBranch(&data)
	if result != nil {
		return ip.defaultStep(*result)
	}

	result = ip.StepLogical(&data)
	if result != nil {
		return ip.defaultStep(*result)
	}

	result = ip.StepBitwise(&data)
	if result != nil {
		return ip.defaultStep(*result)
	}

	result = ip.StepBoolean(&data)
	if result != nil {
		return ip.defaultStep(*result)
	}

	result = ip.StepTools(&data)
	if result != nil {
		return ip.defaultStep(*result)
	}

	proc, ok := ip.procs[data.token.Value]
	if ok {
		ip.calls.Push(ip.ip + 1)
		ip.ip = int(proc.Start)
		return true, nil
	}

	if data.token.Equals("exit", types.TokenTypeKeyword) {
		ip.runtimev("exit command. exiting...\n")
		ip.ip = ip.eop
		return true, nil
	}

	if IsVerboseRuntime() {
		tokenLog := data.token.Format()
		ip.runtimev("invalid token [%s] at ip %d\n", tokenLog, ip.ip)
	}

	return false, fmt.Errorf("failed to run step. unknown operator at ip %d\n", ip.ip)
}
