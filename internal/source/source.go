package source

import (
	"fmt"
	"os"
	"strings"

	"github.com/ktnuity/wet/internal/stdlib"
	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

type ExitCallback = func()

func Load(args *types.WetArgs) (string, ExitCallback) {
	std, err := stdlib.GetContent()
	util.ExitWithError(err, util.AsRef("Failed to load STD Lib"))

	var exit ExitCallback = func() {}

	var sourcePath *string

	var inputSource string
	var source string

	if args.Flags.Is(types.WetFlagHelp) {
		inputSource = "wet_cmd_help\n"
	} else if args.Flags.Is(types.WetFlagVersion) {
		inputSource = "wet_cmd_version\n"
	} else if args.Flags.Is(types.WetFlagLicense) {
		inputSource = "wet_cmd_license\n"
	} else if args.Path == nil {
		exit = func() {
			usage(args.Bin.Name)
		}
		inputSource = "wet_cmd_version\n"
	} else {
		sourcePath = args.Path

		lastIndex := strings.LastIndex(strings.ReplaceAll(*sourcePath, "\\", "/"), "/")
		dir := (*sourcePath)[:lastIndex]
		if dir != "" {
			err = os.Chdir(dir)
			util.ExitWithError(err, util.AsRef("Failed to change directory"))
		}

		inputSource, err = loadFile((*sourcePath)[lastIndex+1:])
		util.ExitWithError(err, util.AsRef("Failed to load input source"))
	}

	inputSourceWithStd := fmt.Sprintf("%s\n%s", std, inputSource)

	source, err = processSource(inputSourceWithStd, 4, args)
	if sourcePath != nil {
		util.ExitWithError(err, util.AsRef(fmt.Sprintf("Failed to load file: %s", *sourcePath)))
	} else {
		util.ExitWithError(err, util.AsRef("Failed to process source"))
	}

	return source, exit
}

func usage(name string) {
	fmt.Printf("Usage: %s [options] <file>\n", name)
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
	if !strings.HasSuffix(fileName, ".wet") {
		return "", fmt.Errorf("failed to include source. file '%s' has invalid suffix.", fileName)
	}

	_, err := os.Stat(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to include source. file '%s' not found: %v", fileName, err)
	}

	source, err := loadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to include source. failed to load file '%s': %v", fileName, err)
	}

	return source, nil
}
