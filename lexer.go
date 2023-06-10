package boule

import (
	"bufio"
	"strconv"
	"strings"
	"unicode"
)

type lexerToken struct {
	position int
	token    Token
	value    interface{}
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
func (l *lexer) Yield() *lexerToken {

	var position int
	var token Token
	var value interface{}

	r, _, err := l.reader.ReadRune()
	if err != nil {
		return &lexerToken{l.position, EOF, EOF.String()}
	}

	switch r {
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

	case '"':
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
		if unicode.IsSpace(r) {
			return l.Yield() // move on to next token

		} else if unicode.IsDigit(r) {
			position = l.position
			if l.backup() == EOF {
				break
			}
			token, value = l.lexNumber()

		} else if unicode.IsLetter(r) {
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

	return &lexerToken{position, token, value}
}

func (l *lexer) backup() Token {
	l.position--
	if err := l.reader.UnreadRune(); err != nil {
		return EOF
	}
	return Token(-1)
}

func (l *lexer) lexEqual() Token {

	l.position++

	r, _, err := l.reader.ReadRune()
	if err != nil {
		return EOF
	}

	if r != '=' {
		return ILLEGAL
	}

	return EQUAL // ==
}

func (l *lexer) lexExclamation() Token {

	l.position++

	r, _, err := l.reader.ReadRune()
	if err != nil {
		return EOF
	}

	if r == '=' { // !=
		return NOT_EQUAL
	}

	if l.backup() == EOF {
		return EOF
	}

	return NOT
}

func (l *lexer) lexGreater() Token {

	l.position++

	r, _, err := l.reader.ReadRune()
	if err != nil {
		return EOF
	}

	if r == '=' { // >=
		return GREATER_OR_EQUAL
	}

	if unicode.IsSpace(r) {
		return GREATER
	}

	return ILLEGAL
}

func (l *lexer) lexLess() Token {

	l.position++

	r, _, err := l.reader.ReadRune()
	if err != nil {
		return EOF
	}

	if r == '=' { // <=
		return LESS_OR_EQUAL
	}

	if unicode.IsSpace(r) {
		return LESS
	}

	return ILLEGAL
}

func (l *lexer) lexAnd() Token {

	l.position++

	r, _, err := l.reader.ReadRune()
	if err != nil {
		return EOF
	}

	if r != '&' {
		return ILLEGAL
	}

	return AND
}

func (l *lexer) lexOr() Token {

	l.position++

	r, _, err := l.reader.ReadRune()
	if err != nil {
		return EOF
	}

	if r != '|' {
		return ILLEGAL
	}

	return OR
}

func (l *lexer) lexNumber() (Token, int) {
	var dotFound bool
	var literal string
	for {
		l.position++

		r, _, err := l.reader.ReadRune()
		if err != nil {
			_ = l.backup()
			value, _ := strconv.Atoi(literal)
			return NUMBER, value
		}

		if !unicode.IsDigit(r) {
			if r == '.' {
				if dotFound {
					return ILLEGAL, 0
				}
				dotFound = true
				literal += "."
				continue
			}
			_ = l.backup()
			value, _ := strconv.Atoi(literal)
			return NUMBER, value
		}

		literal += string(r)
	}
}

func (l *lexer) lexString() (Token, string) {
	var literal string
	for {
		l.position++

		r, _, err := l.reader.ReadRune()
		if err != nil || r == '"' {
			return STRING, literal
		}

		literal += string(r)
	}
}

func (l *lexer) lexIdent() (Token, string) {
	var s string
	for {
		l.position++

		r, _, err := l.reader.ReadRune()
		if err != nil || !unicode.IsLetter(r) {
			_ = l.backup()
			break
		}

		s += string(r)
	}

	return IDENT, s
}
