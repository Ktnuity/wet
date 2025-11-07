package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/ktnuity/wet/internal/types"
)

func AsRef[T any](value T) *T {
	return &value
}

func ExitWithError(err error, message *string) {
	if err != nil {
		if message != nil {
			fmt.Printf("%s\n", *message)
		}
		fmt.Printf("Failure: %v\n", err)
		os.Exit(1)
	}
}

func GetCommandArguments() (*types.WetArgs, error) {
	argv := os.Args
	argc := len(argv)

	binPath := os.Args[0]
	tmpPath := strings.ReplaceAll(binPath, "\\", "/")
	binName := tmpPath[strings.LastIndex(tmpPath, "/")+1:]

	args := types.WetArgs{
		Bin: types.WetBin{ Path: binPath, Name: binName, },
		Flags: 0,
		Path: nil,
	}

	for argi := 1; argi < argc; argi++ {
		if strings.HasPrefix(argv[argi], "--") {
			switch argv[argi] {
			case "--verbose-runtime":
				args.Flags |= types.WetFlagVerboseRuntime
			case "--verbose-tokenize":
				args.Flags |= types.WetFlagVerboseTokenize
			case "--verbose":
				args.Flags |= types.WetFlagVerbose
			case "--dev":
				args.Flags |= types.WetFlagDev
			case "--help":
				args.Flags |= types.WetFlagHelp
			case "--version":
				args.Flags |= types.WetFlagVersion
			case "--license":
				args.Flags |= types.WetFlagLicense
			}
		} else if args.Path == nil {
			args.Path = AsRef(argv[argi])
		} else {
			return nil, &ArgError{
				Message: fmt.Sprintf("More than one file provided.\nUsage: %s [options] <file>", args.Bin.Name),
			}
		}
	}

	return &args, nil
}

type ArgError struct {
	Message string
}

func (ae *ArgError) Error() string {
	first, _, _ := strings.Cut(ae.Message, "\n")
	return strings.ToLower(first)
}
