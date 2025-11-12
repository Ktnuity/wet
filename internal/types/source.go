package types

type Source struct {
	Name			string
	Snippets		[]*SourceSnippet
}

type SourceSnippet struct {
	Name			string
	Start			int
	End				int
	Lines			[]*SourceLine
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
