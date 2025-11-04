package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

type Instruction struct {
	Token		*types.Token
	Next		int64
}

func CreateInstruction(token *types.Token) Instruction {
	return Instruction{
		Token: token,
		Next: -1,
	}
}

func CreateInstructions(tokens []types.Token) []Instruction {
	result := make([]Instruction, len(tokens))

	if len(result) == 0 {
		return result
	}

	for idx, token := range tokens {
		result[idx] = CreateInstruction(&token)
	}

	return result
}

func ProcessTokens(tokens []types.Token) ([]Instruction, error) {
	ipStack := &util.Stack[int64]{}

	instructions := CreateInstructions(tokens)
	for idx := range int64(len(instructions)) {
		instruction := &instructions[idx]

		if instruction.Token.Equals("if", types.TokenTypeKeyword) {
			ipStack.Push(idx)
		} else if instruction.Token.Equals("else", types.TokenTypeKeyword) {
			ip, ok := ipStack.Pop()
			if !ok {
				return nil, fmt.Errorf("failed to process instruction. else reached without ip-stack at index %d", idx)
			}

			other := &instructions[ip]
			if !other.Token.Equals("if", types.TokenTypeKeyword) {
				return nil, fmt.Errorf("failed to process instruction. else reached without if at index %d", idx)
			}

			other.Next = idx + 1

			ipStack.Push(idx + 1)
		} else if instruction.Token.Equals("end", types.TokenTypeKeyword) {
			ip, ok := ipStack.Pop()
			if !ok {
				return nil, fmt.Errorf("failed to process instruction. end reached without ip-stack at index %d", idx)
			}

			other := &instructions[ip]
			if !other.Token.Equals("if", types.TokenTypeKeyword) && !other.Token.Equals("else", types.TokenTypeKeyword) {
				return nil, fmt.Errorf("failed to process instruction. end reached without if or else at index %d", idx)
			}

			other.Next = idx + 1
		}
	}

	if ipStack.Len() > 0 {
		return nil, fmt.Errorf("failed to process instruction. process ended with ip-stack size %d. empty required.", ipStack.Len())
	}

	return instructions, nil
}
