package boule

import (
	"fmt"
	"math/big"

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

	switch lv := left.(type) {
	case bool:
		rv, ok := right.(bool)
		if !ok {
			return false, fmt.Errorf("can't compare type 'bool' with type '%T' (position=%d)", right, l.position)
		}
		switch l.token {
		case EQUAL:
			return lv == rv, nil
		case NOT_EQUAL:
			return lv != rv, nil
		case AND:
			return lv && rv, nil
		case OR:
			return lv || rv, nil
		default:
			return false, fmt.Errorf("type 'bool' only supports the EQUAL, NOT_EQUAL, AND and OR operators (position=%d)", l.position)
		}

	case string:
		rv, ok := right.(string)
		if !ok {
			return false, fmt.Errorf("can't compare type 'string' with type '%T' (position=%d)", right, l.position)
		}
		switch l.token {
		case EQUAL:
			return lv == rv, nil
		case NOT_EQUAL:
			return lv != rv, nil
		default:
			return false, fmt.Errorf("type 'string' only supports the EQUAL and NOT_EQUAL operators (position=%d)", l.position)
		}

	default:
		leftInt, leftBig, leftFloat, leftKind := toNumeric(left)
		rightInt, rightBig, rightFloat, rightKind := toNumeric(right)

		if leftKind == numNone || rightKind == numNone {
			return false, fmt.Errorf("can't compare type '%T' with type '%T' (position=%d)", left, right, l.position)
		}

		if leftKind == numInt64 && rightKind == numInt64 {
			return compareInt64(leftInt, rightInt, l.token, l.position)
		}

		if leftKind == numFloat64 && rightKind == numFloat64 {
			return compareFloat64(leftFloat, rightFloat, l.token, l.position)
		}

		leftBig = promoteToBI(leftInt, leftBig, leftKind)
		rightBig = promoteToBI(rightInt, rightBig, rightKind)

		if leftKind != numFloat64 && rightKind != numFloat64 {
			return compareBigInt(leftBig, rightBig, l.token, l.position)
		}

		if leftKind == numFloat64 {
			return compareFloatBigInt(leftFloat, rightBig, l.token, l.position)
		}
		return compareBigIntFloat(leftBig, rightFloat, l.token, l.position)
	}
}

type numKind int

const (
	numNone numKind = iota
	numInt64
	numBigInt
	numFloat64
)

func toNumeric(v interface{}) (int64, *big.Int, float64, numKind) {
	switch n := v.(type) {
	case int:
		return int64(n), nil, 0, numInt64
	case int8:
		return int64(n), nil, 0, numInt64
	case int16:
		return int64(n), nil, 0, numInt64
	case int32:
		return int64(n), nil, 0, numInt64
	case int64:
		return n, nil, 0, numInt64
	case uint:
		return 0, new(big.Int).SetUint64(uint64(n)), 0, numBigInt
	case uint8:
		return 0, new(big.Int).SetUint64(uint64(n)), 0, numBigInt
	case uint16:
		return 0, new(big.Int).SetUint64(uint64(n)), 0, numBigInt
	case uint32:
		return 0, new(big.Int).SetUint64(uint64(n)), 0, numBigInt
	case uint64:
		return 0, new(big.Int).SetUint64(uint64(n)), 0, numBigInt
	case float32:
		return 0, nil, float64(n), numFloat64
	case float64:
		return 0, nil, n, numFloat64
	case *big.Int:
		return 0, n, 0, numBigInt
	default:
		return 0, nil, 0, numNone
	}
}

func promoteToBI(i64 int64, bi *big.Int, kind numKind) *big.Int {
	if kind == numInt64 {
		return big.NewInt(i64)
	}
	return bi
}

func compareInt64(l, r int64, token Token, pos int) (interface{}, error) {
	switch token {
	case EQUAL:
		return l == r, nil
	case NOT_EQUAL:
		return l != r, nil
	case LESS:
		return l < r, nil
	case LESS_OR_EQUAL:
		return l <= r, nil
	case GREATER:
		return l > r, nil
	case GREATER_OR_EQUAL:
		return l >= r, nil
	default:
		return false, fmt.Errorf("numeric types only support the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_EQUAL operators (position=%d)", pos)
	}
}

func compareFloat64(l, r float64, token Token, pos int) (interface{}, error) {
	switch token {
	case EQUAL:
		return l == r, nil
	case NOT_EQUAL:
		return l != r, nil
	case LESS:
		return l < r, nil
	case LESS_OR_EQUAL:
		return l <= r, nil
	case GREATER:
		return l > r, nil
	case GREATER_OR_EQUAL:
		return l >= r, nil
	default:
		return false, fmt.Errorf("numeric types only support the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_EQUAL operators (position=%d)", pos)
	}
}

func compareBigInt(l, r *big.Int, token Token, pos int) (interface{}, error) {
	switch token {
	case EQUAL:
		return l.Cmp(r) == 0, nil
	case NOT_EQUAL:
		return l.Cmp(r) != 0, nil
	case LESS:
		return l.Cmp(r) == -1, nil
	case LESS_OR_EQUAL:
		return l.Cmp(r) <= 0, nil
	case GREATER:
		return l.Cmp(r) == 1, nil
	case GREATER_OR_EQUAL:
		return l.Cmp(r) >= 0, nil
	default:
		return false, fmt.Errorf("numeric types only support the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_EQUAL operators (position=%d)", pos)
	}
}

func compareBigIntFloat(l *big.Int, r float64, token Token, pos int) (interface{}, error) {
	rBigFloat := big.NewFloat(r)
	rRounded, accuracy := rBigFloat.Int(nil)
	switch token {
	case EQUAL:
		return l.Cmp(rRounded) == 0 && accuracy == big.Exact, nil
	case NOT_EQUAL:
		return l.Cmp(rRounded) != 0 || accuracy != big.Exact, nil
	case LESS:
		return (l.Cmp(rRounded) == 0 && accuracy == big.Below) || l.Cmp(rRounded) == -1, nil
	case LESS_OR_EQUAL:
		return (l.Cmp(rRounded) == 0 && (accuracy == big.Exact || accuracy == big.Below)) || l.Cmp(rRounded) == -1, nil
	case GREATER:
		return (l.Cmp(rRounded) == 0 && accuracy == big.Above) || l.Cmp(rRounded) == 1, nil
	case GREATER_OR_EQUAL:
		return (l.Cmp(rRounded) == 0 && (accuracy == big.Exact || accuracy == big.Above)) || l.Cmp(rRounded) == 1, nil
	default:
		return false, fmt.Errorf("numeric types only support the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_EQUAL operators (position=%d)", pos)
	}
}

func compareFloatBigInt(l float64, r *big.Int, token Token, pos int) (interface{}, error) {
	lBigFloat := big.NewFloat(l)
	lRounded, accuracy := lBigFloat.Int(nil)
	switch token {
	case EQUAL:
		return lRounded.Cmp(r) == 0 && accuracy == big.Exact, nil
	case NOT_EQUAL:
		return lRounded.Cmp(r) != 0 || accuracy != big.Exact, nil
	case LESS:
		return (lRounded.Cmp(r) == 0 && accuracy == big.Above) || lRounded.Cmp(r) == -1, nil
	case LESS_OR_EQUAL:
		return (lRounded.Cmp(r) == 0 && (accuracy == big.Exact || accuracy == big.Above)) || lRounded.Cmp(r) == -1, nil
	case GREATER:
		return (lRounded.Cmp(r) == 0 && accuracy == big.Below) || lRounded.Cmp(r) == 1, nil
	case GREATER_OR_EQUAL:
		return (lRounded.Cmp(r) == 0 && (accuracy == big.Exact || accuracy == big.Below)) || lRounded.Cmp(r) == 1, nil
	default:
		return false, fmt.Errorf("numeric types only support the EQUAL, NOT_EQUAL, LESS, LESS_OR_EQUAL, GREATER and GREATER_OR_EQUAL operators (position=%d)", pos)
	}
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
