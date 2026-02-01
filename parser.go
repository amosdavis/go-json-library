package json

import (
	"fmt"
	"strconv"
)

// Parser parses JSON tokens into Go values
type Parser struct {
	lexer   *Lexer
	current Token
}

// NewParser creates a new parser for the given input
func NewParser(input string) *Parser {
	p := &Parser{
		lexer: NewLexer(input),
	}
	p.advance()
	return p
}

func (p *Parser) advance() {
	p.current = p.lexer.NextToken()
}

func (p *Parser) error(msg string) error {
	return fmt.Errorf("parse error at line %d, column %d: %s", p.current.Line, p.current.Column, msg)
}

// Parse parses the JSON input and returns the resulting value
func (p *Parser) Parse() (interface{}, error) {
	value, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	if p.current.Type != TokenEOF {
		return nil, p.error("unexpected token after value")
	}

	return value, nil
}

func (p *Parser) parseValue() (interface{}, error) {
	switch p.current.Type {
	case TokenLeftBrace:
		return p.parseObject()
	case TokenLeftBracket:
		return p.parseArray()
	case TokenString:
		value := p.current.Value
		p.advance()
		return value, nil
	case TokenNumber:
		return p.parseNumber()
	case TokenTrue:
		p.advance()
		return true, nil
	case TokenFalse:
		p.advance()
		return false, nil
	case TokenNull:
		p.advance()
		return nil, nil
	case TokenError:
		return nil, p.error(p.current.Value)
	case TokenEOF:
		return nil, p.error("unexpected end of input")
	default:
		return nil, p.error(fmt.Sprintf("unexpected token: %s", p.current.Type))
	}
}

func (p *Parser) parseObject() (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	p.advance() // consume '{'

	if p.current.Type == TokenRightBrace {
		p.advance()
		return obj, nil
	}

	for {
		// Expect string key
		if p.current.Type != TokenString {
			return nil, p.error("expected string key in object")
		}
		key := p.current.Value
		p.advance()

		// Expect colon
		if p.current.Type != TokenColon {
			return nil, p.error("expected ':' after object key")
		}
		p.advance()

		// Parse value
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		obj[key] = value

		// Check for comma or end of object
		if p.current.Type == TokenRightBrace {
			p.advance()
			return obj, nil
		}
		if p.current.Type != TokenComma {
			return nil, p.error("expected ',' or '}' in object")
		}
		p.advance()
	}
}

func (p *Parser) parseArray() ([]interface{}, error) {
	arr := make([]interface{}, 0)

	p.advance() // consume '['

	if p.current.Type == TokenRightBracket {
		p.advance()
		return arr, nil
	}

	for {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)

		if p.current.Type == TokenRightBracket {
			p.advance()
			return arr, nil
		}
		if p.current.Type != TokenComma {
			return nil, p.error("expected ',' or ']' in array")
		}
		p.advance()
	}
}

func (p *Parser) parseNumber() (float64, error) {
	value, err := strconv.ParseFloat(p.current.Value, 64)
	if err != nil {
		return 0, p.error(fmt.Sprintf("invalid number: %s", p.current.Value))
	}
	p.advance()
	return value, nil
}
