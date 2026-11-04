package tokenizer

import (
	"fmt"

	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

var argFlags types.WetFlag
func SubmitFlags(flags types.WetFlag) {
	argFlags = flags
}

func IsVerbose() bool {
	return util.HasAnyBitFlags(argFlags, types.WetFlagVerbose)
}

func IsVerboseTokenizer() bool {
	return util.HasBitFlag(argFlags, types.WetFlagVerboseTokenize)
}

func tokenv(format string, args...any) {
	if !IsVerboseTokenizer() || len(format) <= 0 {
		return
	}

	colorYellow := "\033[33m"
	colorReset := "\033[0m"

	fmt.Printf("%sTokenizer:%s ", colorYellow, colorReset)
	fmt.Printf(format, args...)

	if format[len(format)-1] != '\n' {
		fmt.Println()
	}
}
