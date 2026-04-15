package boule

import (
	"bufio"
	"math/big"
	"strconv"
	"strings"
)

// LexerToken pairs a token type with its parsed value.
type LexerToken struct {
	token Token
	value interface{}
}

type lexerTokenWithPosition struct {
	LexerToken
	position int
}

type lexer struct {
	position int
	reader   *bufio.Reader
}

func newLexer(input string) *lexer {
	return &lexer{
		position: 0,
		reader:   bufio.NewReader(strings.NewReader(input)),
	}
}

// Yield scans the string for the next token. It returns the position of the token,
// the token's type, and the literal identifier.
func (l *lexer) Yield() *lexerTokenWithPosition {

	var position int
	var token Token
	var value interface{}

	c, err := l.reader.ReadByte()
	if err != nil {
		return &lexerTokenWithPosition{LexerToken: LexerToken{token: EOF, value: EOF.String()}, position: l.position}
	}

	switch c {
	case '=':
		position = l.position
		token = l.lexEqual()
		value = token.String()

	case '!':
		position = l.position
		token = l.lexExclamation()
		value = token.String()

	case '>':
		position = l.position
		token = l.lexGreater()
		value = token.String()

	case '<':
		position = l.position
		token = l.lexLess()
		value = token.String()

	case '&':
		position = l.position
		token = l.lexAnd()
		value = token.String()

	case '|':
		position = l.position
		token = l.lexOr()
		value = token.String()

	case '"', '\'':
		position = l.position
		token, value = l.lexString()

	case '(':
		position = l.position
		token = OPEN
		value = OPEN.String()

	case ')':
		position = l.position
		token = CLOSE
		value = CLOSE.String()

	default:
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			return l.Yield() // move on to next token

		} else if c >= '0' && c <= '9' {
			position = l.position
			if l.backup() == EOF {
				break
			}
			token, value = l.lexNumber()

		} else if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			position = l.position
			if l.backup() == EOF {
				break
			}
			token, value = l.lexIdent()

		} else {
			position = l.position
			token = ILLEGAL
			value = ILLEGAL.String()
		}
	}

	l.position++

	return &lexerTokenWithPosition{LexerToken: LexerToken{token: token, value: value}, position: position}
}

func (l *lexer) backup() Token {
	l.position--
	if err := l.reader.UnreadByte(); err != nil {
		return EOF
	}
	return Token(-1)
}

func (l *lexer) lexEqual() Token {

	l.position++

	c, err := l.reader.ReadByte()
	if err != nil {
		return EOF
	}

	if c != '=' {
		return ILLEGAL
	}

	return EQUAL // ==
}

func (l *lexer) lexExclamation() Token {

	l.position++

	c, err := l.reader.ReadByte()
	if err != nil {
		return EOF
	}

	if c == '=' { // !=
		return NOT_EQUAL
	}

	if l.backup() == EOF {
		return EOF
	}

	return NOT
}

func (l *lexer) lexGreater() Token {

	l.position++

	c, err := l.reader.ReadByte()
	if err != nil {
		return EOF
	}

	if c == '=' { // >=
		return GREATER_OR_EQUAL
	}

	if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
		return GREATER
	}

	return ILLEGAL
}

func (l *lexer) lexLess() Token {

	l.position++

	c, err := l.reader.ReadByte()
	if err != nil {
		return EOF
	}

	if c == '=' { // <=
		return LESS_OR_EQUAL
	}

	if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
		return LESS
	}

	return ILLEGAL
}

func (l *lexer) lexAnd() Token {

	l.position++

	c, err := l.reader.ReadByte()
	if err != nil {
		return EOF
	}

	if c != '&' {
		return ILLEGAL
	}

	return AND
}

func (l *lexer) lexOr() Token {

	l.position++

	c, err := l.reader.ReadByte()
	if err != nil {
		return EOF
	}

	if c != '|' {
		return ILLEGAL
	}

	return OR
}

func (l *lexer) lexNumber() (Token, interface{}) {
	var dotFound bool
	var b strings.Builder
	for {
		l.position++

		c, err := l.reader.ReadByte()
		if err != nil {
			literal := b.String()
			if dotFound {
				value, err := strconv.ParseFloat(literal, 64)
				if err != nil {
					return ILLEGAL, 0
				}
				return FLOAT, value
			} else {
				integer, ok := (&big.Int{}).SetString(literal, 10)
				if !ok {
					return ILLEGAL, 0
				}
				return INTEGER, integer
			}
		}

		if c < '0' || c > '9' {
			if c == '.' {
				if dotFound {
					return ILLEGAL, 0
				}
				dotFound = true
				b.WriteByte('.')
				continue
			}
			_ = l.backup()
			literal := b.String()
			if dotFound {
				value, err := strconv.ParseFloat(literal, 64)
				if err != nil {
					return ILLEGAL, 0
				}
				return FLOAT, value
			} else {
				integer, ok := (&big.Int{}).SetString(literal, 10)
				if !ok {
					return ILLEGAL, 0
				}
				return INTEGER, integer
			}
		}

		b.WriteByte(c)
	}
}

func (l *lexer) lexString() (Token, string) {
	var b strings.Builder
	for {
		l.position++

		c, err := l.reader.ReadByte()
		if err != nil || c == '"' || c == '\'' {
			return STRING, b.String()
		}

		b.WriteByte(c)
	}
}

func (l *lexer) lexIdent() (Token, string) {
	var b strings.Builder
	for {
		l.position++

		c, err := l.reader.ReadByte()
		if err != nil {
			break
		}
		if !isLetter(c) && c != '_' && c != '.' {
			_ = l.backup()
			break
		}

		b.WriteByte(c)
	}

	return IDENT, b.String()
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
