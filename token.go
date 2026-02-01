package json

// TokenType represents the type of a JSON token
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenLeftBrace    // {
	TokenRightBrace   // }
	TokenLeftBracket  // [
	TokenRightBracket // ]
	TokenColon        // :
	TokenComma        // ,
	TokenString
	TokenNumber
	TokenTrue
	TokenFalse
	TokenNull
	TokenError
)

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenLeftBrace:
		return "{"
	case TokenRightBrace:
		return "}"
	case TokenLeftBracket:
		return "["
	case TokenRightBracket:
		return "]"
	case TokenColon:
		return ":"
	case TokenComma:
		return ","
	case TokenString:
		return "STRING"
	case TokenNumber:
		return "NUMBER"
	case TokenTrue:
		return "TRUE"
	case TokenFalse:
		return "FALSE"
	case TokenNull:
		return "NULL"
	case TokenError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Token represents a JSON token
type Token struct {
	Type    TokenType
	Value   string
	Line    int
	Column  int
}
