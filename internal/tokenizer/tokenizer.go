package tokenizer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

func TokenizeCode(input string) ([]types.Token) {
	result := make([]types.Token, 0, 8)

	commentLess := stripComments(input)
	code := stripExcessWhitespace(commentLess)

	var scan string = code

	for {
		nextScan, word := nextWord(scan)
		if nextScan == nil {
			break
		}

		tokenType := getTokenType(word)
		result = append(result, types.Token{
			Value: word,
			Type: tokenType,
		})

		scan = *nextScan
	}

	return result
}

func LogTokens(tokens []types.Token) error {
	if len(tokens) == 0 {
		return fmt.Errorf("failed to log tokens. no tokens present")
	}

	fmt.Printf("Token Count: %d\n", len(tokens))

	for idx, token := range tokens {
		format := token.Format()
		fmt.Printf("%d : %s\n", idx, format)
	}

	return nil
}

func stripComments(str string) string {
	lines := strings.Split(str, "\n")
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "//") {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func stripExcessWhitespace(str string) string {
	return strings.Join(strings.Fields(str), " ")
}

func nextWord(str string) (*string, string) {
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return nil, ""
	}

	if str[0] == '"' {
		end := strings.Index(str[1:], "\"")
		if end == -1 {
			return &str, ""
		}
		return util.AsRef(strings.TrimSpace(str[end+2:])), str[:end+2]
	}

	parts := strings.Fields(str)
	if len(parts) == 0 {
		return nil, ""
	}

	word := parts[0]
	rest := strings.TrimSpace(str[len(word):])
	return &rest, word
}

var numberRegex = regexp.MustCompile(`^-?(\d+\.?\d*|\.\d+)$`)

func isNumber(str string) bool {
	return numberRegex.MatchString(str)
}

var keywords = map[string]bool{
	"while": true, "until": true, "do": true, "end": true,
	"if": true, "unless": true, "else": true,
	"dup": true, "over": true, "swap": true, "2dup": true, "2swap": true, "drop": true, "nop": true,
	"store": true, "load": true,
	"download": true, "move": true, "copy": true, "exist": true, "touch": true, "rm": true,
	"unzip": true, "lsf": true, "getf": true, "lsd": true, "getd": true,
	"concat": true, "tostring": true, "token": true, "absolute": true, "relative": true,
	"true": true, "false": true,
}

func isKeyword(str string) bool {
	return keywords[str]
}

func isString(str string) bool {
	return len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"'
}

func isPath(str string) bool {
	return strings.HasPrefix(str, "./") || strings.HasPrefix(str, "/")
}

var symbols = map[string]bool{
	"+": true, "-": true, "/": true, "*": true, "%": true, "++": true, "--": true,
	"&": true, "|": true, "^": true, "~": true,
	"==": true, "!=": true, "<": true, "<=": true, ">": true, ">=": true, "!": true,
	".": true,
}

func isSymbol(str string) bool {
	return symbols[str]
}

func getTokenType(str string) types.TokenType {
	if isNumber(str) {
		return types.TokenTypeNumber
	} else if isKeyword(str) {
		return types.TokenTypeKeyword
	} else if isSymbol(str) {
		return types.TokenTypeSymbol
	} else if isPath(str) {
		return types.TokenTypePath
	} else if isString(str) {
		return types.TokenTypeString
	} else {
		return types.TokenTypeNone
	}
}
