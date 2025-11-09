package types

import (
	"fmt"
	"strconv"
	"strings"
)

type TokenType uint8
const (
	TokenTypeNone TokenType = iota
	TokenTypeNumber
	TokenTypeKeyword
	TokenTypePath
	TokenTypeSymbol
	TokenTypeString
)

type Token struct {
	Value		string
	Type		TokenType
}

func GetTokenTypeName(tokenType TokenType) string {
	switch tokenType {
	case TokenTypeNone: return "None"
	case TokenTypeNumber: return "Number"
	case TokenTypeKeyword: return "Keyword"
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
	case TokenTypeNumber, TokenTypeKeyword, TokenTypeSymbol:
		return fmt.Sprintf("%s(%s)", typeName, token.Value)
	case TokenTypePath:
		return EscapePath(token.Value)
	case TokenTypeNone:
		return fmt.Sprintf("%s[[%s]]", typeName, token.Value)
	default:
		return fmt.Sprintf("%s", token.Value)
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

		if t.Value == "" {
			return false
		}
	}

	return t.Value == value
}

func (t *Token) GetNumberValue() (int, bool) {
	if t == nil {
		return 0, false
	}

	if !t.Equals("", TokenTypeNumber) {
		return 0, false
	}

	if t.Value == "" {
		return 0, false
	}

	i, err := strconv.ParseInt(t.Value, 10, 64)
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

	if t.Value == "" {
		return "", false
	}

	return UnescapeString(t.Value), true
}

func (t *Token) GetPathValue() (string, bool) {
	if t == nil {
		return "", false
	}
	
	if !t.Equals("", TokenTypePath) {
		return "", false
	}

	if t.Value == "" {
		return "", false
	}

	return UnescapeString(t.Value), true
}

func UnescapeString(str string) string {
	if len(str) < 2 || str[0] != '"' || str[len(str)-1] != '"' {
		return str
	}

	// Remove quotes
	content := str[1 : len(str)-1]

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
