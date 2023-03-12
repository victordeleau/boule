package boule

import (
	"errors"
	"io"
)

var ErrInvalidSyntax = errors.New("invalid syntax")

type AST struct {
	program Expression
	lexer   *lexer
	current *lexerToken
	peek    *lexerToken
}

func newAST(lexer *lexer) (*AST, error) {

	ast := &AST{
		lexer: lexer,
		current: &lexerToken{
			token: OPEN,
		},
		peek: &lexerToken{
			token: OPEN,
		},
	}

	var err error
	if err = ast.next(); err != nil {
		return nil, err
	}
	if err = ast.next(); err != nil {
		return nil, err
	}

	ast.program, err = ast.expression()
	return ast, err
}

func (a *AST) next() error {

	a.current = a.peek
	a.peek = a.lexer.Yield()

	if a.current.token == EOF {
		return io.EOF
	}

	return nil
}

func (a *AST) expression() (Expression, error) {

	var err error
	if a.current.token == OPEN {

		position := a.current.position

		if err = a.next(); err != nil {
			return nil, ErrInvalidSyntax
		}

		expression, err := a.expression()
		if err != nil {
			return nil, err
		}

		if err = a.next(); err != nil {
			return nil, ErrInvalidSyntax
		}

		if a.current.token != CLOSE {
			return nil, ErrInvalidSyntax
		}

		return &GroupingExpression{&Grouping{
			openPosition:  position,
			Expression:    expression,
			closePosition: a.current.position,
		}}, nil
	}

	return a.suffixExpression()
}

func (a *AST) suffixExpression() (Expression, error) {

	var err error
	if a.current.token.Literal() {

		var prefixExpression Expression
		switch a.current.token {
		case INTEGER:
			prefixExpression = &LiteralExpression{&LiteralInteger{
				value:    a.current.value.(int),
				position: a.current.position,
			}}
		case STRING:
			prefixExpression = &LiteralExpression{&LiteralString{
				value:    a.current.value.(string),
				position: a.current.position,
			}}
		default:
			prefixExpression = &LiteralExpression{&LiteralIdent{
				token:    a.current.token,
				position: a.current.position,
			}}
		}

		if a.peek.token.BinaryOperator() {

			operator := &Operator{
				token:    a.peek.token,
				position: a.peek.position,
			}

			if err = a.next(); err != nil {
				return nil, ErrInvalidSyntax
			}

			if err = a.next(); err != nil {
				return nil, ErrInvalidSyntax
			}

			suffixExpression, err := a.suffixExpression()
			if err != nil {
				return nil, err
			}

			return &BinaryExpression{&Binary{
				left:     prefixExpression,
				Operator: operator,
				right:    suffixExpression,
			}}, nil
		}

		return prefixExpression, nil
	}

	if a.current.token.UnaryOperator() {

		position := a.current.position

		if err = a.next(); err != nil {
			return nil, ErrInvalidSyntax
		}

		expression, err := a.expression()
		if err != nil {
			return nil, err
		}

		return &UnaryExpression{&UnaryNot{
			Expression: expression,
			position:   position,
		}}, nil
	}

	return nil, ErrInvalidSyntax
}
