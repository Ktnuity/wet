package interpreter

import (
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

func IsVerboseRuntime() bool {
	return util.HasBitFlag(argFlags, types.WetFlagVerboseRuntime)
}

func IsVerboseTypeCheck() bool {
	return util.HasBitFlag(argFlags, types.WetFlagVerboseTypeCheck)
}
