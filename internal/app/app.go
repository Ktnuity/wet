package app

import (
	"log"

	"github.com/ktnuity/wet/internal/tokenizer"
	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

func EntryPoint(src string, args *types.WetArgs) error {
	if util.HasFlag(args.Flags, types.WetFlagVerboseRuntime) {
		log.Printf("Code:\n%s\n", src)
		log.Printf("Tokenizing code...\n")
	}
	tokens := tokenizer.TokenizeCode(src)
	if util.HasFlag(args.Flags, types.WetFlagVerboseTokenize) {
		log.Printf("Tokens:\n")
		tokenizer.LogTokens(tokens)
	}



	return nil
}
