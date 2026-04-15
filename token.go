package boule

// Token represents the type of a lexical token in the expression language.
type Token int

const (
	EOF Token = iota
	ILLEGAL

	// literal
	INTEGER
	FLOAT
	STRING
	IDENT

	// binary operator
	EQUAL
	NOT_EQUAL
	GREATER
	GREATER_OR_EQUAL
	LESS
	LESS_OR_EQUAL
	AND
	OR

	// unary operator
	NOT

	// group
	OPEN
	CLOSE
)

var tokens = map[Token]string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",

	// literal
	INTEGER: "INTEGER",
	FLOAT:   "FLOAT",
	STRING:  "STRING",
	IDENT:   "IDENT",

	// binary operator
	EQUAL:            "==",
	NOT_EQUAL:        "!=",
	GREATER:          ">",
	GREATER_OR_EQUAL: ">=",
	LESS:             "<",
	LESS_OR_EQUAL:    "<=",

	// boolean operator
	AND: "&&",
	OR:  "||",

	// unary operator
	NOT: "!",

	// group
	OPEN:  "(",
	CLOSE: ")",
}

// String returns the human-readable representation of the token.
func (t Token) String() string {
	return tokens[t]
}

// Valid reports whether the token is a usable language token (not EOF or ILLEGAL).
func (t Token) Valid() bool {
	if t < 2 {
		return false
	}
	return true
}

// Literal reports whether the token is a literal type (INTEGER, FLOAT, STRING, or IDENT).
func (t Token) Literal() bool {
	return t > 1 && t < 6
}

// BinaryOperator reports whether the token is a binary operator (comparison or logical).
func (t Token) BinaryOperator() bool {
	return t > 5 && t < 14
}

// BooleanOperator reports whether the token is a boolean connective (AND or OR).
func (t Token) BooleanOperator() bool {
	return t > 11 && t < 14
}

// UnaryOperator reports whether the token is a unary operator (NOT).
func (t Token) UnaryOperator() bool {
	return t == 14
}

// Group reports whether the token is a grouping delimiter (OPEN or CLOSE).
func (t Token) Group() bool {
	return t > 14
}
