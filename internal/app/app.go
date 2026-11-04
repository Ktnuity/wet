package app

import (
	"fmt"

	"github.com/ktnuity/wet/internal/interpreter"
	"github.com/ktnuity/wet/internal/tokenizer"
	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

func EntryPoint(src *types.Source, args *types.WetArgs) error {
	if util.HasFlag(args.Flags, types.WetFlagVerboseRuntime) {
		fmt.Printf("Code:\n%s\n", src)
		fmt.Printf("Tokenizing code...\n")
	}

	tokenizer.SubmitFlags(args.Flags)
	tokens, err := tokenizer.TokenizeCode(src)
	if err != nil {
		return types.MakeError(fmt.Errorf("failed to init wet: %w", err))
	}

	if util.HasFlag(args.Flags, types.WetFlagVerboseTokenize) {
		fmt.Printf("Tokens:\n")
		tokenizer.LogTokens(tokens)
	}

	interpreter.SubmitFlags(args.Flags)

	intr, err := interpreter.CreateNew(tokens)
	if err != nil {
		return types.MakeError(err)
	}

	status, err := intr.Run()
	if err != nil {
		return &types.RuntimeError{ Err: err }
	}

	if !status {
		return types.MakeError(fmt.Errorf("error running wet."))
	}

	return nil
}
