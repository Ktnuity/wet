package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
)

type StepResult struct {
	status		bool
	result		error
	skipInc		bool
}

type StepData struct {
	inst		*Instruction
	token		*types.Token
}

func StepOk() *StepResult {
	return &StepResult{ true, nil, false }
}

func StepOkSkip() *StepResult {
	return &StepResult{ true, nil, true }
}

func StepBad(err error) *StepResult {
	return &StepResult{ false, err, false }
}

func (ip *Interpreter) newStep() (StepData, *StepResult) {
	if ip == nil {
		return StepData{}, &StepResult{
			status: false,
			result: fmt.Errorf("failed to step interpreter. instance is nil."),
		}
	}

	if ip.ip >= ip.eop {
		return StepData{}, ip.runtimeverr("failed to run step. ip >= eop.")
	}

	inst := &ip.program[ip.ip]
	token := inst.Token

	if token == nil {
		return StepData{}, ip.runtimeverr("failed to run step. token at index %d is nil\n", ip.ip)
	}

	data := StepData{
		inst: inst,
		token: token,
	}

	return data, nil
}

func (ip *Interpreter) validateStep(result *StepResult) bool {
	return result != nil
}

func (ip *Interpreter) defaultStep(result StepResult) (bool, error) {
	if result.status && !result.skipInc {
		ip.ip++
	}

	return result.status, result.result
}
