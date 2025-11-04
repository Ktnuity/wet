package interpreter

import (
	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

var argFlags uint8

func SubmitFlags(flags uint8) {
	argFlags = flags
}

func IsVerbose() bool {
	return util.HasAnyBitFlags(argFlags, types.WetFlagVerbose)
}

func IsVerboseRuntime() bool {
	return util.HasBitFlag(argFlags, types.WetFlagVerboseRuntime)
}
