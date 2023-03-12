package boule

type Token int

const (
	EOF Token = iota
	ILLEGAL

	// literal
	INTEGER
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
	STRING:  "STRING",
	IDENT:   "IDENT",

	// binary operator
	EQUAL:            "==",
	NOT_EQUAL:        "!=",
	GREATER:          ">",
	GREATER_OR_EQUAL: ">=",
	LESS:             "<",
	LESS_OR_EQUAL:    "<=",
	AND:              "&&",
	OR:               "||",

	// unary operator
	NOT: "!",

	// group
	OPEN:  "(",
	CLOSE: ")",
}

func (t Token) String() string {
	return tokens[t]
}

func (t Token) Valid() bool {
	if t < 2 {
		return false
	}
	return true
}

func (t Token) Literal() bool {
	return t > 1 && t < 5
}

func (t Token) BinaryOperator() bool {
	return t > 4 && t < 13
}

func (t Token) UnaryOperator() bool {
	return t == 13
}

func (t Token) Group() bool {
	return t > 13
}
