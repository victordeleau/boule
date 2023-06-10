package boule

import (
	"fmt"
	"github.com/victordeleau/boule/internal/prefixtree"
	"reflect"
)

type Data struct {
	prefixtree.Tree
}

func NewData() *Data {
	return new(Data)
}

/*
Context-Free grammar

expression         -> binary | suffixExpression
suffixExpression   -> grouping | literal | unary
literal            -> NUMBER | STRING | IDENT
unary              -> NOT suffixExpression
binary             -> expression operator suffixExpression
grouping           -> OPEN expression CLOSE
operator           -> EQUAL | NOT_EQUAL | LESS | LESS_EQUAL | GREATER | GREATER_EQUAL | AND | OR
*/

type node interface {
	Evaluate(data *Data) (interface{}, error)
}

type GroupingExpression struct {
	openPosition int
	node
	closePosition int
}

func (l *GroupingExpression) Evaluate(data *Data) (interface{}, error) {
	return l.node.Evaluate(data)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type UnaryExpression struct {
	node
	position int
}

func (l *UnaryExpression) Evaluate(data *Data) (interface{}, error) {

	value, err := l.node.Evaluate(data)
	if err != nil {
		return false, err
	}

	booleanValue, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("unary expression must be of type boolean (position=%d)", l.position)
	}

	return !booleanValue, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type BinaryExpression struct {
	left     node
	token    Token
	position int
	right    node
}

func (l *BinaryExpression) Evaluate(data *Data) (interface{}, error) {

	left, err := l.left.Evaluate(data)
	if err != nil {
		return nil, err
	}
	right, err := l.right.Evaluate(data)
	if err != nil {
		return nil, err
	}

	leftType, rightType := reflect.TypeOf(left), reflect.TypeOf(right)
	leftValue, rightValue := reflect.ValueOf(left), reflect.ValueOf(right)
	leftKind, rightKind := leftType.Kind(), rightType.Kind()

	if leftKind == reflect.Bool {

		if leftKind != rightType.Kind() {
			return false, fmt.Errorf("can't compare type '%s' with type '%s'", leftKind, rightKind)
		}

		leftBool, rightBool := left.(bool), right.(bool)

		if l.token == EQUAL {
			return leftBool == rightBool, nil
		}

		if l.token == NOT_EQUAL {
			return leftBool != rightBool, nil
		}

		if l.token == AND {
			return leftBool && rightBool, nil
		}

		if l.token == OR {
			return leftBool || rightBool, nil
		}

		fmt.Printf(">>> %+v\n\n", l.token)

		return false, fmt.Errorf("type '%s' only supports the EQUAL, NOT_EQUAL, AND and OR operators", leftKind.String())

	}
	if leftKind == reflect.String {

		if leftKind != rightType.Kind() {
			return false, fmt.Errorf("can't compare type '%s' with type '%s'", leftKind, rightKind)
		}

		leftString, rightString := left.(string), right.(string)

		if l.token == EQUAL {
			return leftString == rightString, nil
		}

		if l.token == NOT_EQUAL {
			return leftString != rightString, nil
		}

		return false, fmt.Errorf("type '%s' only supports the EQUAL and NOT_EQUAL operators", leftKind.String())

	}

	if isInt(leftKind) && isInt(rightKind) {

		leftInt, rightInt := leftValue.Int(), rightValue.Int()
		switch l.token {
		case EQUAL:
			return leftInt == rightInt, nil
		case NOT_EQUAL:
			return leftInt != rightInt, nil
		case LESS:
			return leftInt < rightInt, nil
		case LESS_OR_EQUAL:
			return leftInt <= rightInt, nil
		case GREATER:
			return leftInt > rightInt, nil
		case GREATER_OR_EQUAL:
			return leftInt >= rightInt, nil
		default:
			return false, fmt.Errorf("type '%s' only supports the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_GREATER operators", leftKind.String())
		}

	}

	if isUint(leftKind) && isUint(rightKind) {

		leftUint, rightUint := leftValue.Uint(), rightValue.Uint()
		switch l.token {
		case EQUAL:
			return leftUint == rightUint, nil
		case NOT_EQUAL:
			return leftUint != rightUint, nil
		case LESS:
			return leftUint < rightUint, nil
		case LESS_OR_EQUAL:
			return leftUint <= rightUint, nil
		case GREATER:
			return leftUint > rightUint, nil
		case GREATER_OR_EQUAL:
			return leftUint >= rightUint, nil
		default:
			return false, fmt.Errorf("type '%s' only supports the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_GREATER operators", leftKind.String())
		}
	}

	if isFloat(leftKind) && isFloat(rightKind) {

		leftFloat, rightFloat := leftValue.Float(), rightValue.Float()
		switch l.token {
		case EQUAL:
			return leftFloat == rightFloat, nil
		case NOT_EQUAL:
			return leftFloat != rightFloat, nil
		case LESS:
			return leftFloat < rightFloat, nil
		case LESS_OR_EQUAL:
			return leftFloat <= rightFloat, nil
		case GREATER:
			return leftFloat > rightFloat, nil
		case GREATER_OR_EQUAL:
			return leftFloat >= rightFloat, nil
		default:
			return false, fmt.Errorf("type '%s' only supports the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_GREATER operators", leftKind.String())
		}
	}

	if isInt(leftKind) && isFloat(rightKind) {

		leftInt, rightFloat := leftValue.Int(), rightValue.Float()
		switch l.token {
		case EQUAL:
			return float64(leftInt) == rightFloat, nil
		case NOT_EQUAL:
			return float64(leftInt) != rightFloat, nil
		case LESS:
			return float64(leftInt) < rightFloat, nil
		case LESS_OR_EQUAL:
			return float64(leftInt) <= rightFloat, nil
		case GREATER:
			return float64(leftInt) > rightFloat, nil
		case GREATER_OR_EQUAL:
			return float64(leftInt) >= rightFloat, nil
		default:
			return false, fmt.Errorf("type '%s' and '%s' only supports the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_GREATER operators", leftKind.String(), rightKind.String())
		}
	}

	if isFloat(leftKind) && isInt(rightKind) {

		leftFloat, rightInt := leftValue.Float(), rightValue.Int()
		switch l.token {
		case EQUAL:
			return leftFloat == float64(rightInt), nil
		case NOT_EQUAL:
			return leftFloat != float64(rightInt), nil
		case LESS:
			return leftFloat < float64(rightInt), nil
		case LESS_OR_EQUAL:
			return leftFloat <= float64(rightInt), nil
		case GREATER:
			return leftFloat > float64(rightInt), nil
		case GREATER_OR_EQUAL:
			return leftFloat >= float64(rightInt), nil
		default:
			return false, fmt.Errorf("type '%s' and '%s' only supports the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_GREATER operators", leftKind.String(), rightKind.String())
		}
	}

	// floats can be compared to ints, and conversely

	return false, fmt.Errorf("can't compare type '%s' with type '%s'", leftKind, rightKind)
}

func isInt(kind reflect.Kind) bool {
	return kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64
}

func isUint(kind reflect.Kind) bool {
	return kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64
}

func isFloat(kind reflect.Kind) bool {
	return kind == reflect.Float32 || kind == reflect.Float64
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// integer

type LiteralNumber struct {
	value    int
	position int
}

func (l *LiteralNumber) Evaluate(_ *Data) (interface{}, error) {
	return l.value, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// string

type LiteralString struct {
	value    string
	position int
}

func (l *LiteralString) Evaluate(_ *Data) (interface{}, error) {
	return l.value, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// identifier

type LiteralIdent struct {
	identifier string
	position   int
}

func (l *LiteralIdent) Evaluate(data *Data) (interface{}, error) {
	return data.Find(l.identifier)
}
