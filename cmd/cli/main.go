package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ktnuity/wet/internal/app"
	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

func main() {
	args, err := getCommandArguments()
	if err != nil {
		if v, ok := err.(*ArgError); ok {
			fmt.Printf("%s\n", v.Message)
		} else {
			fmt.Printf("Error: %v\n", err)
		}

		return
	}

	if args.Path == nil {
		fmt.Printf("Usage: %s [options] <file>\n", args.Bin.Name)
		return
	}

	sourcePath := *args.Path


	source, err := loadFile(sourcePath)
	ExitWithError(err, util.AsRef(fmt.Sprintf("Failed to load file: %s", *args.Path)))

	dir := sourcePath[:strings.LastIndex(strings.ReplaceAll(sourcePath, "\\", "/"), "/")]
	if dir != "" {
		err = os.Chdir(dir)
		ExitWithError(err, util.AsRef("Failed to change directory"))
	}

	err = app.EntryPoint(source, args)
	ExitWithError(err, util.AsRef("Runtime failure"))
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

func loadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getCommandArguments() (*types.WetArgs, error) {
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
			}
		} else if args.Path == nil {
			args.Path = util.AsRef(argv[argi])
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
