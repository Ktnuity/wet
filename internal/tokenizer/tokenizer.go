package tokenizer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ktnuity/wet/internal/types"
)

type Scan struct {
	scan		*types.SourceLine
	next		[]*types.SourceLine
}

func TokenizeCode(input *types.Source) ([]types.Token, error) {
	result := make([]types.Token, 0, 8)

	stripComments(input)
	stripExcessWhitespace(input)

	scan := &Scan{
		scan: nil,
		next: input.Lines(),
	}

	entry := 1

	for {
		word := nextWord(scan, &entry)
		if word == nil {
			break
		}

		tokenType, err := getTokenType(word)
		if err != nil {
			return result, fmt.Errorf("failed to tokenize code: %w", err)
		}

		token := types.Token{
			Word: word,
			Type: tokenType,
		}

		if word.Equals("proc") {
			word = nextWord(scan, &entry)
			if word == nil {
				break
			}

			status, err := validateNormalName(word)
			if !status {
				if err != nil {
					return result, fmt.Errorf("failed to tokenize code: %w", err)
				} else {
					return result, fmt.Errorf("failed to tokenize code: proc name '%s' is not valid. %s.", word.UnwrapName(), word.InlineTrace())
				}
			}

			procName := word

			tokenv("proc name: %s\n", procName.UnwrapName())

			outTypes := make([]types.ValueType, 0, 8)

			for {
				word = nextWord(scan, &entry)
				if word == nil {
					return result, fmt.Errorf("failed to tokenize code: proc out type experienced pre-mature exit.")
				}

				tokenv("Out word: %s\n", word.UnwrapName())

				if word.Any("in", "do") {
					break
				}

				outType := types.ParseValueType(word)
				if !outType.Any() {
					return result, fmt.Errorf("failed to tokenize code: proc out type '%s' is not valid. %s.", word.UnwrapName(), word.InlineTrace())
				}

				outTypes = append(outTypes, outType)
			}

			inTypes := make([]types.ValueType, 0, 8)

			if !word.Equals("do") {
				for {
					word = nextWord(scan, &entry)
					if word == nil {
						return result, fmt.Errorf("failed to tokenize code: proc in type experienced pre-mature exit.")
					}

					tokenv("In word: %s\n", word.UnwrapName())

					if word.Equals("do") {
						break
					}

					inType := types.ParseValueType(word)
					if !inType.Any() {
						return result, fmt.Errorf("failed to tokenize code: proc in type '%s' is not valid. %s.", word.UnwrapName(), word.InlineTrace())
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
	}

	return result, nil
}

func validateNormalName(word *types.Word) (bool, error) {
	re, err := regexp.Compile("^[a-z](_?[a-z0-9]+)*$")
	if err != nil {
		return false, fmt.Errorf("failed to validate proc name. proc name regex failed to compile %s: %w", word.InlineTrace(), err)
	}

	if !re.MatchString(word.UnwrapName()) {
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

func stripComments(src *types.Source) {

	for _, snippet := range src.Snippets {
		filtered := make([]*types.SourceLine, 0, len(snippet.Lines))
		for _, line := range snippet.Lines {
			if !strings.HasPrefix(line.Content, "#") && !strings.HasPrefix(line.Content, "//") {
				filtered = append(filtered, line)
			}
		}
		snippet.Lines = filtered
	}

}

func stripExcessWhitespace(src *types.Source) {
	for _, snippet := range src.Snippets {
		filtered := make([]*types.SourceLine, 0, len(snippet.Lines))

		for _, line := range snippet.Lines {
			stripped := stripLineWhitespace(line.Content)
			if len(stripped) > 0 {
				line.Content = stripped
				filtered = append(filtered, line)
			}
		}

		snippet.Lines = filtered
	}
}

func stripLineWhitespace(str string) string {
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

func scanValid(scan *Scan, entry *int) bool {
	if scan.scan == nil {
		if len(scan.next) == 0 {
			return false
		}

		scan.scan = scan.next[0]
		scan.next = scan.next[1:]
		*entry = 1
	}

	scan.scan.Content = strings.TrimSpace(scan.scan.Content)
	for len(scan.scan.Content) == 0 {
		if len(scan.next) == 0 {
			return false
		}

		scan.scan = scan.next[0]
		scan.next = scan.next[1:]
		scan.scan.Content = strings.TrimSpace(scan.scan.Content)
		*entry = 1
	}

	return true
}

func nextWord(scan *Scan, entry *int) *types.Word {
	if !scanValid(scan, entry) {
		return nil
	}

	str := scan.scan.Content

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
				scan.scan.Content = strings.TrimSpace(str[i+1:])
				(*entry)++
				return types.NewWord(str[:i+1], scan.scan, *entry - 1)
			}
		}
		// No closing quote found - return entire string as error token
		scan.scan.Content = ""
		return types.NewWord(str, scan.scan, *entry)
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
				scan.scan.Content = strings.TrimSpace(str[i+1:])
				(*entry)++
				return types.NewWord(strings.ReplaceAll(str[:i], "\\", ""), scan.scan, *entry - 1)
			}
		}

		scan.scan.Content = ""
		(*entry)++
		return types.NewWord(str, scan.scan, *entry - 1)
	}

	parts := strings.Fields(str)
	if len(parts) == 0 {
		scan.scan.Content = ""
		return nil
	}

	word := parts[0]
	scan.scan.Content = strings.TrimSpace(str[len(word):])
	(*entry)++
	return types.NewWord(word, scan.scan, *entry - 1)
}

var numberRegex = regexp.MustCompile(`^-?(\d+\.?\d*|\.\d+)$`)

func isNumber(word *types.Word) bool {
	return numberRegex.MatchString(word.UnwrapName())
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

func isKeyword(word *types.Word) bool {
	return keywords[word.UnwrapName()]
}

func isIdentifier(word *types.Word) (bool, error) {
	return validateNormalName(word)
}

func isString(word *types.Word) bool {
	name := word.UnwrapName()
	return len(name) >= 2 && name[0] == '"' && name[len(name)-1] == '"'
}

func isPath(word *types.Word) bool {
	return strings.HasPrefix(word.UnwrapName(), "./") || strings.HasPrefix(word.UnwrapName(), "/") || strings.HasPrefix(word.UnwrapName(), ":")
}

var symbols = map[string]bool{
	"+": true, "-": true, "/": true, "*": true, "%": true, "++": true, "--": true,
	"&": true, "|": true, "^": true, "~": true,
	"&&": true, "||": true,
	"=": true, "!=": true, "<": true, "<=": true, ">": true, ">=": true, "!": true,
	".": true,
}

func isSymbol(word *types.Word) bool {
	return symbols[word.UnwrapName()]
}

func getTokenType(word *types.Word) (types.TokenType, error) {
	if isNumber(word) {
		return types.TokenTypeNumber, nil
	} else if isKeyword(word) {
		return types.TokenTypeKeyword, nil
	}

	identifier, err := isIdentifier(word)
	if err != nil {
		return types.TokenTypeNone, fmt.Errorf("failed to get token type: %w", err)
	}

	if identifier {
		return types.TokenTypeIdentifier, nil
	} else if isSymbol(word) {
		return types.TokenTypeSymbol, nil
	} else if isPath(word) {
		return types.TokenTypePath, nil
	} else if isString(word) {
		return types.TokenTypeString, nil
	} else {
		return types.TokenTypeNone, nil
	}
}
