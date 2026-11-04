package types

import (
	"fmt"
	"strconv"
	"strings"
)

type Source struct {
	Name			string
	Snippets		[]*SourceSnippet
}

type SourceSnippet struct {
	Name			string
	Start			int
	End				int
	Parent			*SourceParent
	Lines			[]*SourceLine
}

type SourceParent struct {
	Snippet			*SourceSnippet
	Line			int
}

type SourceLine struct {
	Parent			*SourceSnippet
	Content			string
	Line			int
}

func (s *Source) Lines() []*SourceLine {
	if s == nil || s.Snippets == nil {
		return []*SourceLine{}
	}

	result := make([]*SourceLine, 0, 8)

	for _, snippet := range s.Snippets {
		if snippet == nil {
			continue
		}

		for _, line := range snippet.Lines {
			if line == nil {
				continue
			}

			result = append(result, line)
		}
	}

	return result
}

func (s *Source) LineCount() int {
	count := 0
	for _, snippet := range s.Snippets {
		count += len(snippet.Lines)
	}

	return count
}

func (s *Source) LineSpan() int {
	max := 0
	for _, snippet := range s.Snippets {
		count := len(strconv.Itoa(len(snippet.Lines)))
		if count > max {
			max = count
		}
	}

	return max
}

func (s *Source) Unwrap() string {
	lineCount := s.LineCount()

	lines := make([]string, 0, lineCount)

	for _, snippet := range s.Snippets {
		for _, line := range snippet.Lines {
			lines = append(lines, line.Content)
		}
	}

	return strings.Join(lines, "\n")
}

func (s *Source) String() string {
	lineCount := s.LineCount()
	lineSpan := s.LineSpan() + 1

	// post fix line
	pfLine := func(l int) string {
		str := strconv.Itoa(l)
		return str + strings.Repeat(" ", lineSpan - len(str))
	}

	lines := make([]string, 0, lineCount)

	for _, snippet := range s.Snippets {
		lines = append(lines, fmt.Sprintf("\033[32m%s[%d:%d]\033[0m", snippet.Name, snippet.Start, snippet.End))

		for _, line := range snippet.Lines {
			lines = append(lines, fmt.Sprintf(" \033[36m%s\033[0m: %s", pfLine(line.Line), line.Content))
		}
	}

	return strings.Join(lines, "\n")
}
