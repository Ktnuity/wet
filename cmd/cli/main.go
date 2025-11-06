package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
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

	var source string
	if args.Flags & types.WetFlagHelp == types.WetFlagHelp {
		source, err = processSource("@include std.wet\nhelp\n", 4, args)
		ExitWithError(err, util.AsRef("Failed to load help file"))
	} else {
		if args.Path == nil {
			fmt.Printf("Usage: %s [options] <file>\n", args.Bin.Name)
			return
		}

		sourcePath := *args.Path

		source, err = loadSource(sourcePath, 4, args)
		ExitWithError(err, util.AsRef(fmt.Sprintf("Failed to load file: %s", *args.Path)))

		dir := sourcePath[:strings.LastIndex(strings.ReplaceAll(sourcePath, "\\", "/"), "/")]
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

func loadSource(path string, maxDepth int, args *types.WetArgs) (string, error) {
	file, err := loadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load source. loadFile failed: %v", err)
	}

	file = "@include std.wet\n" + file

	return processSource(file, maxDepth, args)
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
	} else if matched, _ := regexp.MatchString(`^[a-z]+\.wet$`, fileName); matched {
		stdPath := "./std/" + fileName
		if _, statErr := os.Stat(stdPath); statErr == nil && args.Flags & types.WetFlagDev == types.WetFlagDev {
			recSource, err = loadFile(stdPath)
		} else {
			url := "https://raw.githubusercontent.com/Ktnuity/wet/refs/heads/master/std/" + fileName
			resp, httpErr := http.Get(url)
			if httpErr != nil {
				err = fmt.Errorf("failed to download from %s: %v", url, httpErr)
			} else {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					body, readErr := io.ReadAll(resp.Body)
					if readErr != nil {
						err = fmt.Errorf("failed to read response from %s: %v", url, readErr)
					} else {
						recSource = string(body)
					}
				} else {
					err = fmt.Errorf("failed to download from %s: status %d", url, resp.StatusCode)
				}
			}
		}
	} else {
		return "", fmt.Errorf("failed to load source. include directive '%s' yielded invalid result.", fileName)
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
			case "--help":
				args.Flags |= types.WetFlagHelp
			case "--dev":
				args.Flags |= types.WetFlagDev
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
