package types

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type TokenType uint8
const (
	TokenTypeNone TokenType = iota
	TokenTypeNumber
	TokenTypeKeyword
	TokenTypeIdentifier
	TokenTypePath
	TokenTypeSymbol
	TokenTypeString
)

type TokenExtraProc struct {
	Name		*Word
	In			[]ValueType
	Out			[]ValueType
}

type TokenExtra struct {
	Proc		*TokenExtraProc
}

type Word struct {
	word		string
	line		*SourceLine
	entry		int
}

type Token struct {
	Word		*Word
	Type		TokenType
	Extra		TokenExtra
}

func GetTokenTypeName(tokenType TokenType) string {
	switch tokenType {
	case TokenTypeNone: return "None"
	case TokenTypeNumber: return "Number"
	case TokenTypeKeyword: return "Keyword"
	case TokenTypeIdentifier: return "Identifier"
	case TokenTypePath: return "Path"
	case TokenTypeSymbol: return "Symbol"
	case TokenTypeString: return "String"
	default: return "Unknown"
	}
}

var escapePathMap = map[rune]bool{
	'\\': true, ' ': true,
}
func EscapePath(path string) string {
	var sb strings.Builder
	sb.Grow(len(path))

	for _, ch := range path {
		if escapePathMap[ch] {
			sb.WriteRune('\\')
		}

		sb.WriteRune(ch)
	}

	return sb.String()
}


func (token *Token) Format() string {
	if token == nil {
		return "<nil>"
	}

	typeName := GetTokenTypeName(token.Type)

	switch token.Type {
	case TokenTypeNumber, TokenTypeKeyword, TokenTypeSymbol, TokenTypeIdentifier:
		return fmt.Sprintf("%s(%s)", typeName, token.Word.UnwrapName())
	case TokenTypePath:
		return EscapePath(token.Word.UnwrapName())
	case TokenTypeNone:
		return fmt.Sprintf("%s[[%s]]", typeName, token.Word.UnwrapName())
	default:
		return fmt.Sprintf("%s", token.Word.UnwrapName())
	}
}

func (t *Token) Equals(value string, ttype TokenType) bool {
	if t == nil {
		return false
	}

	if t.Type == TokenTypeNone {
		return value == "" && ttype == TokenTypeNone
	}

	if ttype == TokenTypeNone {
		if value == "" {
			return true
		}
	} else if value == "" {
		return ttype == t.Type
	} else {
		if ttype != t.Type {
			return false
		}

		if t.Word.Empty() {
			return false
		}
	}

	return t.Word.Equals(value)
}

func (t *Token) GetNumberValue() (int, bool) {
	if t == nil {
		return 0, false
	}

	if !t.Equals("", TokenTypeNumber) {
		return 0, false
	}

	if t.Word.Empty() {
		return 0, false
	}

	i, err := strconv.ParseInt(t.Word.UnwrapName(), 10, 64)
	if err == nil {
		return int(i), true
	}

	return 0, false
}

func (t *Token) GetStringValue() (string, bool) {
	if t == nil {
		return "", false
	}

	if !t.Equals("", TokenTypeString) {
		return "", false
	}

	if t.Word.Empty() {
		return "", false
	}

	return UnescapeString(t.Word), true
}

func (t *Token) GetPathValue() (string, bool) {
	if t == nil {
		return "", false
	}

	if !t.Equals("", TokenTypePath) {
		return "", false
	}

	if t.Word.Empty() {
		return "", false
	}

	return UnescapeString(t.Word), true
}

func UnescapeString(word *Word) string {
	if len(word.word) < 2 || word.word[0] != '"' || word.word[len(word.word)-1] != '"' {
		return word.word
	}

	// Remove quotes
	content := word.word[1 : len(word.word)-1]

	var result strings.Builder
	result.Grow(len(content))

	escaped := false
	for _, ch := range content {
		if escaped {
			switch ch {
			case 'n':
				result.WriteRune('\n')
			case 't':
				result.WriteRune('\t')
			case 'r':
				result.WriteRune('\r')
			case '\\':
				result.WriteRune('\\')
			case '"':
				result.WriteRune('"')
			default:
				// Unknown escape sequence, keep as-is
				result.WriteRune('\\')
				result.WriteRune(ch)
			}
			escaped = false
		} else if ch == '\\' {
			escaped = true
		} else {
			result.WriteRune(ch)
		}
	}

	return result.String()
}

func NewWord(word string, line *SourceLine, entry int) *Word {
	return &Word{word, line, entry}
}

func (w *Word) Any(names...string) bool {
	return slices.Contains(names, w.word)
}

func (w *Word) Equals(name string) bool {
	return w.word == name
}

func (w *Word) Empty() bool {
	return w.word == ""
}

func (w *Word) UnwrapName() string {
	return w.word
}

func (w *Word) LineNumber() int {
	return w.line.Line
}

func (w *Word) LineFile() string {
	return w.line.Parent.Name
}

func (w *Word) UnwrapLine() *SourceLine {
	return w.line
}

func (w *Word) InlineTrace() string {
	result := make([]string, 0, 8)

	result = append(result, fmt.Sprintf("at %s:%d", w.LineFile(), w.LineNumber()))

	snippet := w.line.Parent.Parent

	for snippet != nil {
		result = append(result, fmt.Sprintf("in %s", snippet.Snippet.Name))
		snippet = snippet.Snippet.Parent
	}

	return strings.Join(result, " ")
}

func (w *Word) EntryNumber() int {
	return w.entry
}
