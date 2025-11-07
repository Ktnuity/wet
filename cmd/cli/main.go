package main

//go:generate rm -rf ./std
//go:generate cp -r ../../wetstd ./std

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"

	"github.com/ktnuity/wet/internal/app"
	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

//go:embed std/*
var stdFS embed.FS

func getStdFiles() ([]string, error) {
	if _, err := fs.Stat(stdFS, "std/std.txt"); err == nil {
		data, err := fs.ReadFile(stdFS, "std/std.txt")
		if err != nil {
			return []string{}, fmt.Errorf("failed to stat std lib index: %v", err)
		}

		input := strings.Split(string(data), "\n")

		result := make([]string, 0, len(input))

		for _, line := range input {
			fileName := strings.TrimSpace(line)
			if len(fileName) == 0 {
				continue
			}

			if !strings.HasSuffix(fileName, ".wet") {
				continue
			}

			result = append(result, fileName)
		}

		return result, nil
	} else {
		return []string{}, fmt.Errorf("failed to stat std lib index.")
	}
}

func getStdContent() (string, error) {
	stdFiles, err := getStdFiles()

	parts := make([]string, 0, len(stdFiles))

	if err != nil {
		return "", fmt.Errorf("failed to load std: %v", err)
	}

	for _, file := range stdFiles {
		content, success, err := getStdFile(file)
		if success {
			parts = append(parts, content)
		} else if err != nil {
			return "", fmt.Errorf("failed to load std file '%s': %v", file, err)
		}
	}

	return strings.Join(parts, "\n"), nil
}

func getStdFile(fileName string) (string, bool, error) {
	if matched, _ := regexp.MatchString(`^[a-z]+\.wet$`, fileName); matched {
		stdPath := "./std/" + fileName
		embedPath := stdPath[2:] // Remove "./" prefix

		if _, err := fs.Stat(stdFS, embedPath); err == nil {
			data, err := fs.ReadFile(stdFS, embedPath)

			if err != nil {
				return "", false, fmt.Errorf("failed to open resource '%s': %v", stdPath, err)
			}

			return string(data), true, nil
		} else {
			return "", false, nil
		}
	} else {
		return "", false, nil
	}
}

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

	var sourcePath *string

	var inputSource string
	var source string
	if args.Flags.Is(types.WetFlagHelp) {
		inputSource = "help\n"
	} else if args.Flags.Is(types.WetFlagVersion) {
		inputSource = "version\n"
	} else if args.Flags.Is(types.WetFlagLicense) {
		inputSource = "license\n"
	} else {
		if args.Path == nil {
			defer fmt.Printf("Usage: %s [options] <file>\n", args.Bin.Name)
			inputSource = "version\n"
		} else {
			sourcePath = args.Path

			inputSource, err = loadFile(*sourcePath)
			ExitWithError(err, util.AsRef("Failed to load input source"))
		}
	}

	std, err := getStdContent()
	ExitWithError(err, util.AsRef("Failed to load STD Lib"))

	inputSourceWithStd := fmt.Sprintf("%s\n%s", std, inputSource)

	source, err = processSource(inputSourceWithStd, 4, args)
	if sourcePath != nil {
		ExitWithError(err, util.AsRef(fmt.Sprintf("Failed to load file: %s", *sourcePath)))
	} else {
		ExitWithError(err, util.AsRef("Failed to process source"))
	}


	if sourcePath != nil {
		dir := (*sourcePath)[:strings.LastIndex(strings.ReplaceAll(*sourcePath, "\\", "/"), "/")]
		if dir != "" {
			err = os.Chdir(dir)
			ExitWithError(err, util.AsRef("Failed to change directory"))
		}
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

func processSource(file string, maxDepth int, args *types.WetArgs) (string, error) {
	if maxDepth <= 0 {
		return "", fmt.Errorf("failed to load source. max depth reached.")
	}

	lines := strings.Split(file, "\n")
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		if strings.HasPrefix(line, "@include ") {
			fileName := line[9:]

			recSource, err := includeSource(fileName, args)
			if err != nil {
				return "", fmt.Errorf("failed to load source: %v", err)
			}

			procSource, err := processSource(recSource, maxDepth - 1, args)
			if err != nil {
				return "", fmt.Errorf("failed to process source: %v", err)
			}

			result = append(result, procSource)
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n"), nil
}

func includeSource(fileName string, args *types.WetArgs) (string, error) {
	var err error
	var recSource string
	if _, statErr := os.Stat(fileName); statErr == nil && fileName != "std.wet" {
		recSource, err = loadFile(fileName)
	} else {
		stdFile, success, err := getStdFile(fileName)
		if success {
			recSource = stdFile
		} else if err != nil {
			return "", fmt.Errorf("failed to load std file: %v", err)
		} else {
			return "", fmt.Errorf("failed to load std file. std file not found.")
		}
	}

	if err != nil {
		return "", fmt.Errorf("failed to load source: %v", err)
	}

	return recSource, nil
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
