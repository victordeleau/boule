package boule

/*
Context-Free grammar

expression         -> grouping | suffixExpression
suffixExpression   -> literal | unary | binary
literal            -> INTEGER | STRING | IDENT
unary              -> NOT suffixExpression
binary             -> suffixExpression operator suffixExpression
grouping           -> OPEN expression CLOSE
operator           -> EQUAL | NOT_EQUAL | LESS | LESS_EQUAL | GREATER | GREATER_EQUAL | AND | OR
*/

type Node interface{}

type Expression interface {
	Node
	Interpret()
	Resolve()
	Analyze()
}

type LiteralExpression struct {
	LiteralInterface
}

func (l *LiteralExpression) Interpret() {}
func (l *LiteralExpression) Resolve()   {}
func (l *LiteralExpression) Analyze()   {}

type UnaryExpression struct {
	*UnaryNot
}

func (l *UnaryExpression) Interpret() {}
func (l *UnaryExpression) Resolve()   {}
func (l *UnaryExpression) Analyze()   {}

type BinaryExpression struct {
	*Binary
}

func (l *BinaryExpression) Interpret() {}
func (l *BinaryExpression) Resolve()   {}
func (l *BinaryExpression) Analyze()   {}

type GroupingExpression struct {
	*Grouping
}

func (l *GroupingExpression) Interpret() {}
func (l *GroupingExpression) Resolve()   {}
func (l *GroupingExpression) Analyze()   {}

// literal

type LiteralInterface interface {
	Node
	Interpret()
	Resolve()
	Analyze()
}

type LiteralInteger struct {
	value    int
	position int
}

func (l *LiteralInteger) Interpret() {}
func (l *LiteralInteger) Resolve()   {}
func (l *LiteralInteger) Analyze()   {}

type LiteralString struct {
	value    string
	position int
}

func (l *LiteralString) Interpret() {}
func (l *LiteralString) Resolve()   {}
func (l *LiteralString) Analyze()   {}

type LiteralIdent struct {
	token    Token
	position int
}

func (l *LiteralIdent) Interpret() {}
func (l *LiteralIdent) Resolve()   {}
func (l *LiteralIdent) Analyze()   {}

// unary

type UnaryNot struct {
	Expression
	position int
}

// binary

type Binary struct {
	left Expression
	*Operator
	right Expression
}

// grouping

type Grouping struct {
	openPosition int
	Expression
	closePosition int
}

// operator

type Operator struct {
	token    Token
	position int
}
