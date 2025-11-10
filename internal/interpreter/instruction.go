package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

type Instruction struct {
	Token		*types.Token
	Next		int64
	Mode		uint8
}

type ProcessTokensResult struct {
	inst		[]Instruction
	proc		map[string]Proc
}

type DoMode	= uint8
const (
	DoModeWhile DoMode = 0
	DoModeUntil DoMode = 1
)

type EndMode = uint8
const (
	EndModeNormal EndMode = 0
	EndModeProc	EndMode = 1
)

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

func scanMacros(tokens []types.Token) ([]types.Token, map[string][]types.Token, error) {
	result := make([]types.Token, 0, len(tokens))

	resultMap := make(map[string][]types.Token)

	for idx := 0; idx < len(tokens); idx++ {
		token := &tokens[idx]

		if token.Value == "macro" {
			idx++
			if idx >= len(tokens) {
				return nil, nil, fmt.Errorf("failed to detect macro name at index %d. reached eof early.", idx)
			}

			macroName := tokens[idx].Value

			idx++
			macroStart := idx
			if idx >= len(tokens) {
				return nil, nil, fmt.Errorf("failed to detect macro start at index %d for macro name '%s'. reached eof early.", idx, macroName)
			}

			end := -1
			scopes := 0

			for idx < len(tokens) {
				endToken := &tokens[idx]

				if endToken.Equals("end", types.TokenTypeNone) {
					if scopes == 0 {
						end = idx
						break
					} else {
						scopes--
					}
				} else if endToken.Value == "if" || endToken.Value == "while" || endToken.Value == "unless" || endToken.Value == "until" {
					scopes++
				} else if endToken.Equals("macro", types.TokenTypeNone) {
					return nil, nil, fmt.Errorf("failed to parse macro body. detected unsupported nested macro at index %d for macro name '%s'.", idx, macroName)
				}

				idx++
			}

			if end == -1 {
				return nil, nil, fmt.Errorf("failed to find end for macro name '%s' definition at index %d.", macroName, idx)
			}

			body := make([]types.Token, end - macroStart)

			for itemIndex := macroStart; itemIndex < end; itemIndex++ {
				body[itemIndex - macroStart] = tokens[itemIndex]
			}

			resultMap[macroName] = body

			idx = end
		} else {
			result = append(result, *token)
		}
	}

	return result, resultMap, nil
}

func expandMacros(tokens []types.Token) ([]types.Token, error) {
	tokens, macroMap, err := scanMacros(tokens)

	if err != nil {
		return nil, fmt.Errorf("failed to expand macros. scan macros failed: %v", err)
	}

	deadLimit := 12
	dirty := true

	for dirty && deadLimit > 0 {
		deadLimit--
		dirty = false

		newTokens := make([]types.Token, 0, len(tokens) * 3 / 2)

		for idx := range len(tokens) {
			token := &tokens[idx]

			body, exists := macroMap[token.Value]
			if exists {
				for _, item := range body {
					newTokens = append(newTokens, item)
				}
				dirty = true
				continue
			} else {
				newTokens = append(newTokens, *token)
			}
		}

		tokens = newTokens
	}

	if dirty {
		return nil, fmt.Errorf("failed to expand macro. dead limit reached.")
	}

	return tokens, nil
}

func ProcessTokens(tokens []types.Token) (*ProcessTokensResult, error) {
	tokens, err := expandMacros(tokens)

	if err != nil {
		return nil, fmt.Errorf("failed to process tokens: %v.", err)
	}

	ipStack := &util.Stack[int64]{}

	procs := make(map[string]Proc)

	instructions := CreateInstructions(tokens)
	for idx := range int64(len(instructions)) {
		instruction := &instructions[idx]

		if instruction.Token.Equals("if", types.TokenTypeKeyword) {
			ipStack.Push(idx)
		} else if instruction.Token.Equals("unless", types.TokenTypeKeyword) {
			ipStack.Push(idx)
		} else if instruction.Token.Equals("proc", types.TokenTypeKeyword) {
			ipStack.Push(idx)
		} else if instruction.Token.Equals("else", types.TokenTypeKeyword) {
			ip, ok := ipStack.Pop()
			if !ok {
				return nil, fmt.Errorf("failed to process instruction. else reached without ip-stack at index %d", idx)
			}

			other := &instructions[ip]
			if !other.Token.Equals("if", types.TokenTypeKeyword) && !other.Token.Equals("unless", types.TokenTypeKeyword) {
				return nil, fmt.Errorf("failed to process instruction. else reached without if at index %d", idx)
			}

			other.Next = idx + 1

			ipStack.Push(idx)
		} else if instruction.Token.Equals("while", types.TokenTypeKeyword) || instruction.Token.Equals("until", types.TokenTypeKeyword) {
			ipStack.Push(idx)
		} else if instruction.Token.Equals("do", types.TokenTypeKeyword) {
			ip, ok := ipStack.Pop()
			if !ok {
				return nil, fmt.Errorf("failed to process instruction. do reached without ip-stack at index %d", idx)
			}

			other := &instructions[ip]
			if !other.Token.Equals("while", types.TokenTypeKeyword) && !other.Token.Equals("until", types.TokenTypeKeyword){
				return nil, fmt.Errorf("failed to process instruction. do reached without while or until at index %d", idx)
			}

			ipStack.Push(idx);
			instruction.Next = ip

			switch other.Token.Value {
			case "while": instruction.Mode = DoModeWhile
			case "until": instruction.Mode = DoModeUntil
			}
		} else if instruction.Token.Equals("end", types.TokenTypeKeyword) {
			ip, ok := ipStack.Pop()
			if !ok {
				return nil, fmt.Errorf("failed to process instruction. end reached without ip-stack at index %d", idx)
			}

			instruction.Mode = EndModeNormal

			other := &instructions[ip]
			if other.Token.Equals("if", types.TokenTypeKeyword) || other.Token.Equals("unless", types.TokenTypeKeyword) {
				other.Next = idx + 1
			} else if other.Token.Equals("else", types.TokenTypeKeyword) {
				other.Next = idx + 1
			} else if other.Token.Equals("do", types.TokenTypeKeyword) {
				doIp := other.Next
				other.Next = idx + 1
				instruction.Next = doIp
			} else if other.Token.Equals("proc", types.TokenTypeKeyword) {
				instruction.Next = other.Next
				instruction.Mode = EndModeProc
				other.Next = idx + 1

				procs[other.Token.Extra] = Proc{
					other.Token.Extra,
					ip + 1,
					idx,
				}
			} else {
				return nil, fmt.Errorf("failed to process instruction. end reached without if or else at index %d", idx)
			}
		}
	}

	if ipStack.Len() > 0 {
		return nil, fmt.Errorf("failed to process instruction. process ended with ip-stack size %d. empty required.", ipStack.Len())
	}

	return &ProcessTokensResult{
		instructions,
		procs,
	}, nil
}
