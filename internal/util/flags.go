package util

import "github.com/ktnuity/wet/internal/types"

func HasFlag(flags, flag types.WetFlag) bool {
	return (flags & flag) == flag
}

func HasAllFlags(flags types.WetFlag, test...types.WetFlag) bool {
	for _, flag := range test {
		if !HasFlag(flags, flag) {
			return false
		}
	}

	return true
}

func HasAnyFlags(flags types.WetFlag, test...types.WetFlag) bool {
	for _, flag := range test {
		if HasFlag(flags, flag) {
			return true
		}
	}

	return false
}

func HasBitFlag(flags, flag types.WetFlag) bool {
	return (flags & flag) == flag
}

func HasAnyBitFlags(flags, test types.WetFlag) bool {
	return (flags & test) != 0
}

func HasAllBitFlags(flags, test types.WetFlag) bool {
	return (flags & test) == test
}
