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

func singleLine(name, content string) (string, *types.SourceSnippet) {
	lines := strings.Split(content, "\n")
	snippet := &types.SourceSnippet{
		Name: name,
		Start: 1,
		End: len(lines),
		Lines: make([]*types.SourceLine, 0, len(lines)),
	}

	for idx, line := range lines {
		snippet.Lines = append(snippet.Lines, &types.SourceLine{
			Parent: snippet,
			Line: idx + 1,
			Content: line,
		})
	}

	return name, snippet
}

func Load(args *types.WetArgs) (*types.Source, ExitCallback) {
	std, err := stdlib.GetContent()
	util.ExitWithError(err, util.AsRef("Failed to load STD Lib"))

	var exit ExitCallback = func() {}

	var sourcePath string

	var inputSource *types.SourceSnippet
	var source []*types.SourceSnippet

	if args.Flags.Is(types.WetFlagHelp) {
		sourcePath, inputSource = singleLine("help.wet", "wet_cmd_help\n")
	} else if args.Flags.Is(types.WetFlagVersion) {
		sourcePath, inputSource = singleLine("version.wet", "wet_cmd_version\n")
	} else if args.Flags.Is(types.WetFlagLicense) {
		sourcePath, inputSource = singleLine("license.wet", "wet_cmd_license\n")
	} else if args.Path == nil {
		exit = func() {
			usage(args.Bin.Name)
		}
		sourcePath, inputSource = singleLine("version_usage.wet", "wet_cmd_version\n")
	} else {
		sourcePath = *args.Path

		lastIndex := strings.LastIndex(strings.ReplaceAll(sourcePath, "\\", "/"), "/")
		dir := sourcePath[:lastIndex]
		if dir != "" {
			err = os.Chdir(dir)
			util.ExitWithError(err, util.AsRef("Failed to change directory"))
		}

		inputSource, err = loadFile(sourcePath[lastIndex+1:])
		util.ExitWithError(err, util.AsRef("Failed to load input source"))
	}

	err = stdlib.Attach(inputSource, std)
	util.ExitWithError(err, util.AsRef(fmt.Sprintf("Failed to load file: %s", sourcePath)))

	source, err = processSource(inputSource, 4, args)
	util.ExitWithError(err, util.AsRef(fmt.Sprintf("Failed to load file: %s", sourcePath)))

	snippets := make([]*types.SourceSnippet, 0, len(std))
	snippets = append(snippets, std...)
	snippets = append(snippets, source...)

	return &types.Source{
		Name: sourcePath,
		Snippets: snippets,
	}, exit
}

func usage(name string) {
	fmt.Printf("Usage: %s [options] <file>\n", name)
}

func loadFile(path string) (*types.SourceSnippet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")

	snippet := &types.SourceSnippet{
		Name: path,
		Start: 1,
		End: len(lines),
		Lines: make([]*types.SourceLine, 0, len(lines)),
	}

	for idx, line := range lines {
		snippet.Lines = append(snippet.Lines, &types.SourceLine{
			Parent: snippet,
			Content: line,
			Line: idx + 1,
		})
	}

	return snippet, nil
}

func attach(snippet *types.SourceSnippet, snippets []*types.SourceSnippet, line int) {
	for _, item := range snippets {
		item.Parent = &types.SourceParent{
			Snippet: snippet,
			Line: line,
		}
	}
}

func processSource(snippet *types.SourceSnippet, maxDepth int, args *types.WetArgs) ([]*types.SourceSnippet, error) {
	if maxDepth <= 0 {
		return nil, fmt.Errorf("failed to load source. max depth reached.")
	}

	snippets := make([]*types.SourceSnippet, 0, 8)

	for _, line := range snippet.Lines {
		if strings.HasPrefix(line.Content, "@include ") {
			fileName := line.Content[9:]

			recSource, err := includeSource(fileName, args)
			if err != nil {
				return nil, fmt.Errorf("failed to load source: %v", err)
			}

			procSource, err := processSource(recSource, maxDepth - 1, args)
			if err != nil {
				return nil, fmt.Errorf("failed to process source: %v", err)
			}

			if snippet.Start < line.Line {
				snippets = append(snippets, &types.SourceSnippet{
					Name: snippet.Name,
					Start: snippet.Start,
					End: line.Line-1,
					Lines: snippet.Lines[:(line.Line-snippet.Start)],
				})
			}

			attach(snippet, procSource, line.Line)

			snippets = append(snippets, procSource...)

			snippet.Lines = snippet.Lines[(line.Line-snippet.Start+1):]
			snippet.Start = line.Line+1
		}
	}

	if len(snippet.Lines) != 0 {
		snippets = append(snippets, snippet)
	}

	return snippets, nil
}

func includeSource(fileName string, args *types.WetArgs) (*types.SourceSnippet, error) {
	if !strings.HasSuffix(fileName, ".wet") {
		return nil, fmt.Errorf("failed to include source. file '%s' has invalid suffix.", fileName)
	}

	_, err := os.Stat(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to include source. file '%s' not found: %v", fileName, err)
	}

	source, err := loadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to include source. failed to load file '%s': %v", fileName, err)
	}

	return source, nil
}
