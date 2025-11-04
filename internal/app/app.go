package app

import (
	"fmt"

	"github.com/ktnuity/wet/internal/interpreter"
	"github.com/ktnuity/wet/internal/tokenizer"
	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

func EntryPoint(src string, args *types.WetArgs) error {
	if util.HasFlag(args.Flags, types.WetFlagVerboseRuntime) {
		fmt.Printf("Code:\n%s\n", src)
		fmt.Printf("Tokenizing code...\n")
	}
	tokens := tokenizer.TokenizeCode(src)
	if util.HasFlag(args.Flags, types.WetFlagVerboseTokenize) {
		fmt.Printf("Tokens:\n")
		tokenizer.LogTokens(tokens)
	}

	interpreter.SubmitFlags(args.Flags)

	intr, err := interpreter.CreateNew(tokens)
	if err != nil {
		fmt.Printf("Failed to init interpreter: %v\n", err)
	}

	status, err := intr.Run()
	if err != nil {
		return fmt.Errorf("error running wet: %v", err)
	}

	if !status {
		return fmt.Errorf("error running wet.")
	}

	return nil
}
