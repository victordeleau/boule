package boule

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/victordeleau/boule/internal/prefixtree"
)

// Data holds the variables that expressions are evaluated against.
type Data struct {
	prefixtree.Tree
}

// NewData returns an empty Data store ready for variable insertion via Add.
func NewData() *Data {
	return new(Data)
}

/*
Context-Free grammar

expression         -> binary | suffixExpression
suffixExpression   -> grouping | literal | unary
literal            -> INTEGER | STRING | IDENT
unary              -> NOT suffixExpression
binary             -> expression operator suffixExpression
grouping           -> OPEN expression CLOSE
operator           -> EQUAL | NOT_EQUAL | LESS | LESS_EQUAL | GREATER | GREATER_EQUAL | AND | OR
*/

// Node represents an evaluable node in the expression AST.
type Node interface {
	Evaluate(data *Data) (interface{}, error)
}

// GroupingExpression represents a parenthesized expression.
type GroupingExpression struct {
	openPosition int
	Node
	closePosition int
}

// Evaluate delegates to the inner expression.
func (l *GroupingExpression) Evaluate(data *Data) (interface{}, error) {
	return l.Node.Evaluate(data)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// UnaryExpression represents a NOT (!) expression.
type UnaryExpression struct {
	Node
	position int
}

// Evaluate returns the logical negation of the inner node.
func (l *UnaryExpression) Evaluate(data *Data) (interface{}, error) {

	value, err := l.Node.Evaluate(data)
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

// BinaryExpression represents a comparison or logical operation between two nodes.
type BinaryExpression struct {
	left     Node
	token    Token
	position int
	right    Node
}

// Evaluate computes the result of the binary operation on the left and right operands.
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

	if leftKind == reflect.Bool { // bool <-> bool comparison

		if leftKind != rightType.Kind() {
			return false, fmt.Errorf("can't compare type '%s' with type '%s' (position=%d)", leftKind, rightKind, l.position)
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

		return false, fmt.Errorf("type '%s' only supports the EQUAL, NOT_EQUAL, AND and OR operators (position=%d)", leftKind.String(), l.position)

	}

	if leftKind == reflect.String { // string <-> string comparison

		if leftKind != rightType.Kind() {
			return false, fmt.Errorf("can't compare type '%s' with type '%s' (position=%d)", leftKind, rightKind, l.position)
		}

		leftString, rightString := left.(string), right.(string)

		if l.token == EQUAL {
			return leftString == rightString, nil
		}

		if l.token == NOT_EQUAL {
			return leftString != rightString, nil
		}

		return false, fmt.Errorf("type '%s' only supports the EQUAL and NOT_EQUAL operators (position=%d)", leftKind.String(), l.position)

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
			return false, fmt.Errorf("type '%s' only supports the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_EQUAL operators (position=%d)", leftKind.String(), l.position)
		}
	}

	var ok bool
	var leftIsBigInt, rightIsBigInt bool
	var leftBigInt, rightBigInt *big.Int

	if isUint(leftKind) {
		leftBigInt, ok = (&big.Int{}).SetString(strconv.FormatUint(leftValue.Uint(), 10), 10)
		if !ok {
			return false, fmt.Errorf("uint to big int conversion failed (position=%d)", l.position)
		}
		leftIsBigInt = true

	} else if isInt(leftKind) {
		leftBigInt = big.NewInt(leftValue.Int())
		leftIsBigInt = true

	} else if leftBigInt, ok = leftValue.Interface().(*big.Int); ok {
		leftIsBigInt = true
	}

	if isUint(rightKind) {
		rightBigInt, ok = (&big.Int{}).SetString(strconv.FormatUint(rightValue.Uint(), 10), 10)
		if !ok {
			return false, fmt.Errorf("uint to big int conversion failed (position=%d)", l.position)
		}
		rightIsBigInt = true

	} else if isInt(rightKind) {
		rightBigInt = big.NewInt(rightValue.Int())
		rightIsBigInt = true

	} else if rightBigInt, ok = rightValue.Interface().(*big.Int); ok {
		rightIsBigInt = true
	}

	var leftIsFloat, rightIsFloat bool
	var leftFloat64, rightFloat64 float64

	if isFloat(leftKind) {
		leftFloat64 = leftValue.Float()
		leftIsFloat = true
	}

	if isFloat(rightKind) {
		rightFloat64 = rightValue.Float()
		rightIsFloat = true
	}

	if leftIsFloat && rightIsFloat {
		switch l.token {
		case EQUAL:
			return leftFloat64 == rightFloat64, nil
		case NOT_EQUAL:
			return leftFloat64 != rightFloat64, nil
		case LESS:
			return leftFloat64 < rightFloat64, nil
		case LESS_OR_EQUAL:
			return leftFloat64 <= rightFloat64, nil
		case GREATER:
			return leftFloat64 > rightFloat64, nil
		case GREATER_OR_EQUAL:
			return leftFloat64 >= rightFloat64, nil
		default:
			return false, fmt.Errorf("type 'float64' only supports the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_EQUAL operators (position=%d)", l.position)
		}
	}

	if leftIsBigInt && rightIsBigInt {
		switch l.token {
		case EQUAL:
			return leftBigInt.Cmp(rightBigInt) == 0, nil
		case NOT_EQUAL:
			return leftBigInt.Cmp(rightBigInt) != 0, nil
		case LESS:
			return leftBigInt.Cmp(rightBigInt) == -1, nil
		case LESS_OR_EQUAL:
			return leftBigInt.Cmp(rightBigInt) <= 0, nil
		case GREATER:
			return leftBigInt.Cmp(rightBigInt) == 1, nil
		case GREATER_OR_EQUAL:
			return leftBigInt.Cmp(rightBigInt) >= 0, nil
		default:
			return false, fmt.Errorf("type 'big.Int' only supports the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_EQUAL operators (position=%d)", l.position)
		}
	}

	if leftIsBigInt && rightIsFloat {
		rightBigFloat := big.NewFloat(rightFloat64)
		rightRoundedFloat, accuracy := rightBigFloat.Int(nil)
		switch l.token {
		case EQUAL:
			return leftBigInt.Cmp(rightRoundedFloat) == 0 && accuracy == big.Exact, nil
		case NOT_EQUAL:
			return leftBigInt.Cmp(rightRoundedFloat) != 0 || accuracy != big.Exact, nil
		case LESS:
			return (leftBigInt.Cmp(rightRoundedFloat) == 0 && accuracy == big.Below) || leftBigInt.Cmp(rightRoundedFloat) == -1, nil
		case LESS_OR_EQUAL:
			return (leftBigInt.Cmp(rightRoundedFloat) == 0 && (accuracy == big.Exact || accuracy == big.Below)) || leftBigInt.Cmp(rightRoundedFloat) == -1, nil
		case GREATER:
			return (leftBigInt.Cmp(rightRoundedFloat) == 0 && accuracy == big.Above) || leftBigInt.Cmp(rightRoundedFloat) == 1, nil
		case GREATER_OR_EQUAL:
			return (leftBigInt.Cmp(rightRoundedFloat) == 0 && (accuracy == big.Exact || accuracy == big.Above)) || leftBigInt.Cmp(rightRoundedFloat) == 1, nil
		default:
			return false, fmt.Errorf("numeric types only support the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_EQUAL operators (position=%d)", l.position)
		}
	}

	if leftIsFloat && rightIsBigInt {
		leftBigFloat := big.NewFloat(leftFloat64)
		leftRoundedFloat, accuracy := leftBigFloat.Int(nil)
		switch l.token {
		case EQUAL:
			return leftRoundedFloat.Cmp(rightBigInt) == 0 && accuracy == big.Exact, nil
		case NOT_EQUAL:
			return leftRoundedFloat.Cmp(rightBigInt) != 0 || accuracy != big.Exact, nil
		case LESS:
			return (leftRoundedFloat.Cmp(rightBigInt) == 0 && accuracy == big.Above) || leftRoundedFloat.Cmp(rightBigInt) == -1, nil
		case LESS_OR_EQUAL:
			return (leftRoundedFloat.Cmp(rightBigInt) == 0 && (accuracy == big.Exact || accuracy == big.Above)) || leftRoundedFloat.Cmp(rightBigInt) == -1, nil
		case GREATER:
			return (leftRoundedFloat.Cmp(rightBigInt) == 0 && accuracy == big.Below) || leftRoundedFloat.Cmp(rightBigInt) == 1, nil
		case GREATER_OR_EQUAL:
			return (leftRoundedFloat.Cmp(rightBigInt) == 0 && (accuracy == big.Exact || accuracy == big.Below)) || leftRoundedFloat.Cmp(rightBigInt) == 1, nil
		default:
			return false, fmt.Errorf("numeric types only support the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_EQUAL operators (position=%d)", l.position)
		}
	}

	return false, fmt.Errorf("can't compare type '%s' with type '%s' (position=%d)", leftKind, rightKind, l.position)
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

// LiteralInteger represents an arbitrary-precision integer literal.
type LiteralInteger struct {
	value    *big.Int
	position int
}

// Evaluate returns the integer value.
func (l *LiteralInteger) Evaluate(_ *Data) (interface{}, error) {
	return l.value, nil
}

// LiteralFloat represents a 64-bit floating-point literal.
type LiteralFloat struct {
	value    float64
	position int
}

// Evaluate returns the float value.
func (l *LiteralFloat) Evaluate(_ *Data) (interface{}, error) {
	return l.value, nil
}

// LiteralString represents a quoted string literal.
type LiteralString struct {
	value    string
	position int
}

// Evaluate returns the string value.
func (l *LiteralString) Evaluate(_ *Data) (interface{}, error) {
	return l.value, nil
}

// LiteralIdent represents a variable reference or boolean keyword (true/false).
type LiteralIdent struct {
	identifier string
	position   int
}

// Evaluate resolves the identifier: "true" and "false" return booleans,
// anything else is looked up in the data store.
func (l *LiteralIdent) Evaluate(data *Data) (interface{}, error) {

	if l.identifier == "true" {
		return true, nil
	}

	if l.identifier == "false" {
		return false, nil
	}

	return data.Find(l.identifier)
}
