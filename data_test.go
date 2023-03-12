package boule

var testCases = []struct {
	string      string
	tokenStream []Token
	ast         Expression
	valid       bool
}{
	// valid tests
	{
		string:      `destination == "Saturn" && traveltime > 30000000`,
		tokenStream: []Token{IDENT, EQUAL, STRING, AND, IDENT, GREATER, INTEGER},
		valid:       true,
	},
	{
		string:      `!(captain == "Henry Cavill") || !arrived`,
		tokenStream: []Token{NOT, OPEN, IDENT, EQUAL, STRING, CLOSE, OR, NOT, IDENT},
		valid:       true,
	},
	{
		string:      `(speed <= 1209843257) && (from == "Mars" || from != "Pluton")`,
		tokenStream: []Token{OPEN, IDENT, LESS_OR_EQUAL, INTEGER, CLOSE, AND, OPEN, IDENT, EQUAL, STRING, OR, IDENT, NOT_EQUAL, STRING, CLOSE},
		valid:       true,
	},

	// invalid tests
	{
		string:      `== "Io"`,
		tokenStream: []Token{EQUAL, STRING},
		valid:       false,
	},
	{
		string:      `!= speed)(`,
		tokenStream: []Token{NOT_EQUAL, IDENT, CLOSE, OPEN},
		valid:       false,
	},
	{
		string:      `!= destination)(`,
		tokenStream: []Token{NOT_EQUAL, IDENT, CLOSE, OPEN},
		valid:       false,
	},
	{
		string:      `239869235 >= speed && (> < || !))`,
		tokenStream: []Token{INTEGER, GREATER_OR_EQUAL, IDENT, AND, OPEN, GREATER, LESS, OR, NOT, CLOSE, CLOSE},
		valid:       false,
	},
}
