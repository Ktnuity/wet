package types

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
