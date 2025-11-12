package tokenizer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ktnuity/wet/internal/types"
	"github.com/ktnuity/wet/internal/util"
)

func TokenizeCode(input string) ([]types.Token, error) {
	result := make([]types.Token, 0, 8)

	commentLess := stripComments(input)
	code := stripExcessWhitespace(commentLess)

	var scan string = code

	for {
		nextScan, word := nextWord(scan)
		if nextScan == nil {
			break
		}

		tokenType, err := getTokenType(word)
		if err != nil {
			return result, fmt.Errorf("failed to tokenize code: %w", err)
		}

		token := types.Token{
			Value: word,
			Type: tokenType,
		}

		if word == "proc" {
			nextScan, word = nextWord(*nextScan)
			if nextScan == nil {
				break
			}

			status, err := validateNormalName(word)
			if !status {
				if err != nil {
					return result, fmt.Errorf("failed to tokenize code: %w", err)
				} else {
					return result, fmt.Errorf("failed to tokenize code: proc name '%s' is not valid.", word)
				}
			}

			procName := word

			tokenv("proc name: %s\n", procName)

			outTypes := make([]types.ValueType, 0, 8)

			for {
				nextScan, word = nextWord(*nextScan)
				if nextScan == nil {
					return result, fmt.Errorf("failed to tokenize code: proc out type experienced pre-mature exit.")
				}

				tokenv("Out word: %s\n", word)

				if word == "in" || word == "do" {
					break
				}
				
				outType := types.ParseValueType(word)
				if !outType.Any() {
					return result, fmt.Errorf("failed to tokenize code: proc out type '%s' is not valid.", word)
				}

				outTypes = append(outTypes, outType)
			}

			inTypes := make([]types.ValueType, 0, 8)

			if word != "do" {
				for {
					nextScan, word = nextWord(*nextScan)
					if nextScan == nil {
						return result, fmt.Errorf("failed to tokenize code: proc in type experienced pre-mature exit.")
					}

					tokenv("In word: %s\n", word)

					if word == "do" {
						break
					}
					
					inType := types.ParseValueType(word)
					if !inType.Any() {
						return result, fmt.Errorf("failed to tokenize code: proc in type '%s' is not valid.", word)
					}

					inTypes = append(inTypes, inType)
				}
			}

			token.Extra.Proc = &types.TokenExtraProc{
				Name: procName,
				In: inTypes,
				Out: outTypes,
			}
		}

		result = append(result, token)
		scan = *nextScan
	}

	return result, nil
}

func validateNormalName(name string) (bool, error) {
	re, err := regexp.Compile("^[a-z](_?[a-z0-9]+)*$")
	if err != nil {
		return false, fmt.Errorf("failed to validate proc name. proc name regex failed to compile: %w", err)
	}

	if !re.MatchString(name) {
		return false, nil
	}

	return true, nil
}

func LogTokens(tokens []types.Token) error {
	if len(tokens) == 0 {
		return fmt.Errorf("failed to log tokens. no tokens present")
	}

	tokenv("Token Count: %d\n", len(tokens))

	for idx, token := range tokens {
		format := token.Format()
		tokenv("%d : %s\n", idx, format)
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
	var result strings.Builder
	result.Grow(len(str))

	inString := false
	escaped := false
	prevWasSpace := false

	for _, ch := range str {
		if escaped {
			result.WriteRune(ch)
			escaped = false
			prevWasSpace = false
			continue
		}

		if ch == '\\' && inString {
			escaped = true
			result.WriteRune(ch)
			prevWasSpace = false
			continue
		}

		if ch == '"' {
			inString = !inString
			result.WriteRune(ch)
			prevWasSpace = false
			continue
		}

		if inString {
			result.WriteRune(ch)
			prevWasSpace = false
			continue
		}

		// Outside string: collapse whitespace
		isSpace := ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
		if isSpace {
			if !prevWasSpace {
				result.WriteRune(' ')
				prevWasSpace = true
			}
		} else {
			result.WriteRune(ch)
			prevWasSpace = false
		}
	}

	return strings.TrimSpace(result.String())
}

func nextWord(str string) (*string, string) {
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return nil, ""
	}

	if str[0] == '"' {
		escaped := false
		for i := 1; i < len(str); i++ {
			if escaped {
				escaped = false
				continue
			}
			if str[i] == '\\' {
				escaped = true
				continue
			}
			if str[i] == '"' {
				// Found the closing quote
				return util.AsRef(strings.TrimSpace(str[i+1:])), str[:i+1]
			}
		}
		// No closing quote found - return entire string as error token
		return util.AsRef(""), str
	} else if str[0] == '.' || str[0] == '/' || str[0] == ':' {
		escaped := false
		for i := 1; i < len(str); i++ {
			if escaped {
				escaped = false
				continue
			}
			if str[i] == '\\' {
				escaped = true
				continue
			}
			if str[i] == ' ' {
				return util.AsRef(strings.TrimSpace(str[i+1:])), strings.ReplaceAll(str[:i], "\\", "")
			}
		}

		return util.AsRef(""), str
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
	"ret": true, "iret": true, "dret": true,
	"dup": true, "over": true, "swap": true, "2dup": true, "2swap": true, "drop": true, "nop": true,
	"store": true, "iload": true, "sload": true,
	"download": true, "move": true, "copy": true, "exist": true, "touch": true, "mkdir": true, "rm": true, "readfile": true,
	"unzip": true, "lsf": true, "getf": true, "lsd": true, "getd": true,
	"concat": true, "tostring": true, "token": true, "absolute": true, "relative": true,
	"true": true, "false": true,
	"puts": true,
	"int": true, "string": true,
	"proc": true,
	"exit": true,
}

func isKeyword(str string) bool {
	return keywords[str]
}

func isIdentifier(str string) (bool, error) {
	return validateNormalName(str)
}

func isString(str string) bool {
	return len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"'
}

func isPath(str string) bool {
	return strings.HasPrefix(str, "./") || strings.HasPrefix(str, "/") || strings.HasPrefix(str, ":")
}

var symbols = map[string]bool{
	"+": true, "-": true, "/": true, "*": true, "%": true, "++": true, "--": true,
	"&": true, "|": true, "^": true, "~": true,
	"&&": true, "||": true,
	"=": true, "!=": true, "<": true, "<=": true, ">": true, ">=": true, "!": true,
	".": true,
}

func isSymbol(str string) bool {
	return symbols[str]
}

func getTokenType(str string) (types.TokenType, error) {
	if isNumber(str) {
		return types.TokenTypeNumber, nil
	} else if isKeyword(str) {
		return types.TokenTypeKeyword, nil
	}

	identifier, err := isIdentifier(str)
	if err != nil {
		return types.TokenTypeNone, fmt.Errorf("failed to get token type: %w", err)
	}

	if identifier {
		return types.TokenTypeIdentifier, nil
	} else if isSymbol(str) {
		return types.TokenTypeSymbol, nil
	} else if isPath(str) {
		return types.TokenTypePath, nil
	} else if isString(str) {
		return types.TokenTypeString, nil
	} else {
		return types.TokenTypeNone, nil
	}
}
