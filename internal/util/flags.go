package util

func HasFlag(flags, flag uint8) bool {
	return (flags & flag) == flag
}

func HasAllFlags(flags uint8, test...uint8) bool {
	for _, flag := range test {
		if !HasFlag(flags, flag) {
			return false
		}
	}

	return true
}

func HasAnyFlags(flags uint8, test...uint8) bool {
	for _, flag := range test {
		if HasFlag(flags, flag) {
			return true
		}
	}

	return false
}

func HasBitFlag(flags, flag uint8) bool {
	return (flags & flag) == flag
}

func HasAnyBitFlags(flags, test uint8) bool {
	return (flags & test) != 0
}

func HasAllBitFlags(flags, test uint8) bool {
	return (flags & test) == test
}
