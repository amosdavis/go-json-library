package json

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lexer tokenizes JSON input
type Lexer struct {
	input   string
	pos     int
	line    int
	column  int
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
	}
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

func (l *Lexer) advance() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	r, size := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += size
	if r == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return r
}

func (l *Lexer) skipWhitespace() {
	for {
		r := l.peek()
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			l.advance()
		} else {
			break
		}
	}
}

func (l *Lexer) makeToken(typ TokenType, value string, line, col int) Token {
	return Token{Type: typ, Value: value, Line: line, Column: col}
}

func (l *Lexer) errorToken(msg string) Token {
	return Token{Type: TokenError, Value: msg, Line: l.line, Column: l.column}
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	line := l.line
	col := l.column

	if l.pos >= len(l.input) {
		return l.makeToken(TokenEOF, "", line, col)
	}

	r := l.peek()

	switch r {
	case '{':
		l.advance()
		return l.makeToken(TokenLeftBrace, "{", line, col)
	case '}':
		l.advance()
		return l.makeToken(TokenRightBrace, "}", line, col)
	case '[':
		l.advance()
		return l.makeToken(TokenLeftBracket, "[", line, col)
	case ']':
		l.advance()
		return l.makeToken(TokenRightBracket, "]", line, col)
	case ':':
		l.advance()
		return l.makeToken(TokenColon, ":", line, col)
	case ',':
		l.advance()
		return l.makeToken(TokenComma, ",", line, col)
	case '"':
		return l.scanString()
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return l.scanNumber()
	case 't':
		return l.scanKeyword("true", TokenTrue)
	case 'f':
		return l.scanKeyword("false", TokenFalse)
	case 'n':
		return l.scanKeyword("null", TokenNull)
	default:
		return l.errorToken(fmt.Sprintf("unexpected character: %c", r))
	}
}

func (l *Lexer) scanString() Token {
	line := l.line
	col := l.column
	var sb strings.Builder

	l.advance() // consume opening quote

	for {
		r := l.peek()
		if r == 0 {
			return l.errorToken("unterminated string")
		}
		if r == '"' {
			l.advance()
			return l.makeToken(TokenString, sb.String(), line, col)
		}
		if r == '\\' {
			l.advance()
			escaped := l.peek()
			switch escaped {
			case '"':
				sb.WriteRune('"')
			case '\\':
				sb.WriteRune('\\')
			case '/':
				sb.WriteRune('/')
			case 'b':
				sb.WriteRune('\b')
			case 'f':
				sb.WriteRune('\f')
			case 'n':
				sb.WriteRune('\n')
			case 'r':
				sb.WriteRune('\r')
			case 't':
				sb.WriteRune('\t')
			case 'u':
				l.advance()
				hex := ""
				for i := 0; i < 4; i++ {
					h := l.peek()
					if !isHexDigit(h) {
						return l.errorToken("invalid unicode escape sequence")
					}
					hex += string(h)
					l.advance()
				}
				var codepoint rune
				fmt.Sscanf(hex, "%x", &codepoint)
				sb.WriteRune(codepoint)
				continue
			default:
				return l.errorToken(fmt.Sprintf("invalid escape character: %c", escaped))
			}
			l.advance()
		} else if r < 0x20 {
			return l.errorToken("control character in string")
		} else {
			sb.WriteRune(r)
			l.advance()
		}
	}
}

func (l *Lexer) scanNumber() Token {
	line := l.line
	col := l.column
	start := l.pos

	// Optional minus
	if l.peek() == '-' {
		l.advance()
	}

	// Integer part
	if l.peek() == '0' {
		l.advance()
	} else if isDigit19(l.peek()) {
		l.advance()
		for isDigit(l.peek()) {
			l.advance()
		}
	} else {
		return l.errorToken("invalid number")
	}

	// Fractional part
	if l.peek() == '.' {
		l.advance()
		if !isDigit(l.peek()) {
			return l.errorToken("invalid number: expected digit after decimal point")
		}
		for isDigit(l.peek()) {
			l.advance()
		}
	}

	// Exponent part
	if l.peek() == 'e' || l.peek() == 'E' {
		l.advance()
		if l.peek() == '+' || l.peek() == '-' {
			l.advance()
		}
		if !isDigit(l.peek()) {
			return l.errorToken("invalid number: expected digit in exponent")
		}
		for isDigit(l.peek()) {
			l.advance()
		}
	}

	return l.makeToken(TokenNumber, l.input[start:l.pos], line, col)
}

func (l *Lexer) scanKeyword(keyword string, tokenType TokenType) Token {
	line := l.line
	col := l.column

	for _, expected := range keyword {
		if l.peek() != expected {
			return l.errorToken(fmt.Sprintf("expected '%s'", keyword))
		}
		l.advance()
	}

	// Make sure keyword is not followed by alphanumeric
	if unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek()) {
		return l.errorToken(fmt.Sprintf("unexpected character after '%s'", keyword))
	}

	return l.makeToken(tokenType, keyword, line, col)
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isDigit19(r rune) bool {
	return r >= '1' && r <= '9'
}

func isHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}
