package boule

import (
	"fmt"
	"github.com/victordeleau/boule/prefixtree"
	"io"
)

type AST struct {
	program node
	lexer   *lexer
	current *lexerToken
	peek    *lexerToken
}

func NewBouleExpression(input string) (func(data *prefixtree.Tree) (bool, error), error) {

	ast := &AST{
		lexer: newLexer(input),
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
	if err != nil {
		return nil, err
	}

	return func(data *prefixtree.Tree) (bool, error) {
		result, err := ast.program.Evaluate(data)
		if err != nil {
			return false, err
		}

		resultBoolean, ok := result.(bool)
		if !ok {
			return false, fmt.Errorf("can't evaluate non-boolean identifier")
		}
		return resultBoolean, nil
	}, nil
}

func (a *AST) next() error {

	a.current = a.peek
	a.peek = a.lexer.Yield()

	if a.current.token == EOF {
		return io.EOF
	}

	return nil
}

func (a *AST) expression() (node, error) {

	prefixExpression, err := a.suffixExpression()
	if err != nil {
		return nil, err
	}

	if a.peek.token.BinaryOperator() {

		token := a.peek.token
		position := a.peek.position

		if err = a.next(); err != nil {
			return nil, err
		}

		if err = a.next(); err != nil {
			return nil, err
		}

		suffixExpression, err := a.suffixExpression()
		if err != nil {
			return nil, err
		}

		binaryExpression := &BinaryExpression{
			left:     prefixExpression,
			token:    token,
			position: position,
			right:    suffixExpression,
		}

		if a.peek.token.BooleanOperator() {

			token = a.peek.token
			position = a.peek.position

			if err = a.next(); err != nil {
				return nil, err
			}

			if err = a.next(); err != nil {
				return nil, err
			}

			expression, err := a.expression()
			if err != nil {
				return nil, err
			}

			return &BinaryExpression{
				left:     binaryExpression,
				token:    token,
				position: position,
				right:    expression,
			}, nil
		}

		return binaryExpression, nil
	}

	return prefixExpression, nil
}

func (a *AST) suffixExpression() (node, error) {

	var expression node
	var err error

	if a.current.token.Literal() {

		switch a.current.token {
		case NUMBER:
			return &LiteralNumber{
				value:    a.current.value.(int),
				position: a.current.position,
			}, nil
		case STRING:
			return &LiteralString{
				value:    a.current.value.(string),
				position: a.current.position,
			}, nil
		default:

			valueString, ok := a.current.value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid syntax: raw identifier is not of type 'string' (position=%d)", a.current.position)
			}

			return &LiteralIdent{
				identifier: valueString,
				position:   a.current.position,
			}, nil
		}
	}

	if a.current.token.UnaryOperator() {

		position := a.current.position

		if err = a.next(); err != nil {
			return nil, err
		}

		expression, err = a.suffixExpression()
		if err != nil {
			return nil, err
		}

		return &UnaryExpression{
			node:     expression,
			position: position,
		}, nil
	}

	if a.current.token == OPEN {

		position := a.current.position

		if err = a.next(); err != nil {
			return nil, err
		}

		expression, err = a.expression()
		if err != nil {
			return nil, err
		}

		if err = a.next(); err != nil {
			return nil, err
		}

		if a.current.token != CLOSE {
			return nil, fmt.Errorf("invalid syntax: group expression not closed (position=%d)", a.current.position)
		}

		return &GroupingExpression{
			openPosition:  position,
			node:          expression,
			closePosition: a.current.position,
		}, nil
	}

	return nil, fmt.Errorf("invalid suffix expression (position=%d)", a.current.position)
}
