package boule

import "math/big"

var testCases = []struct {
	string      string
	tokenStream []Token
	data        map[string]interface{}
	ast         node
	valid       bool
	result      bool
}{
	// valid tests
	{
		string:      `destination == "Saturn" && traveltime > 30000000`,
		tokenStream: []Token{IDENT, EQUAL, STRING, AND, IDENT, GREATER, INTEGER},
		data: map[string]interface{}{
			"destination": "Saturn",
			"traveltime":  50000000,
		},
		valid:  true,
		result: true,
	},
	{
		string:      `destination == 'Saturn'`,
		tokenStream: []Token{IDENT, EQUAL, STRING},
		data: map[string]interface{}{
			"destination": "Saturn",
		},
		valid:  true,
		result: true,
	},
	{
		string:      `arrived == true`,
		tokenStream: []Token{IDENT, EQUAL, IDENT},
		data: map[string]interface{}{
			"arrived": true,
		},
		valid:  true,
		result: true,
	},
	{
		string:      `core_frequency < 345738983260257983`,
		tokenStream: []Token{IDENT, LESS, INTEGER},
		data: map[string]interface{}{
			"core_frequency": big.NewInt(3895679862),
		},
		valid:  true,
		result: true,
	},
	{
		string:      `destination == "Saturn" && speed > 280.32`,
		tokenStream: []Token{IDENT, EQUAL, STRING, AND, IDENT, GREATER, FLOAT},
		data: map[string]interface{}{
			"destination": "Saturn",
			"speed":       300.89,
		},
		valid:  true,
		result: true,
	},
	{
		string:      `!(captain == "Henry Cavill") || !arrived`,
		tokenStream: []Token{NOT, OPEN, IDENT, EQUAL, STRING, CLOSE, OR, NOT, IDENT},
		data: map[string]interface{}{
			"captain": "Henry Cavill",
			"arrived": false,
		},
		valid:  true,
		result: true,
	},
	{
		string:      `(speed <= 1209843257) && (from == "Mars" || from != "Pluton")`,
		tokenStream: []Token{OPEN, IDENT, LESS_OR_EQUAL, INTEGER, CLOSE, AND, OPEN, IDENT, EQUAL, STRING, OR, IDENT, NOT_EQUAL, STRING, CLOSE},
		data: map[string]interface{}{
			"speed": 20000,
			"from":  "Mars",
		},
		valid:  true,
		result: true,
	},

	// invalid tests
	{
		string:      `destination == "Saturn" && speed > 280.32. && speed < 1000`,
		tokenStream: []Token{IDENT, EQUAL, STRING, AND, IDENT, GREATER, ILLEGAL, AND, IDENT, LESS, INTEGER},
		data: map[string]interface{}{
			"destination": "Saturn",
			"speed":       300.89,
		},
		valid:  false,
		result: false,
	},
	{
		string:      `== "Io"`,
		tokenStream: []Token{EQUAL, STRING},
		data:        map[string]interface{}{},
		valid:       false,
	},
	{
		string:      `!= speed)(`,
		tokenStream: []Token{NOT_EQUAL, IDENT, CLOSE, OPEN},
		data:        map[string]interface{}{},
		valid:       false,
	},
	{
		string:      `!= destination)(`,
		tokenStream: []Token{NOT_EQUAL, IDENT, CLOSE, OPEN},
		data:        map[string]interface{}{},
		valid:       false,
	},
	{
		string:      `239869235 >= speed && (> < || !))`,
		tokenStream: []Token{INTEGER, GREATER_OR_EQUAL, IDENT, AND, OPEN, GREATER, LESS, OR, NOT, CLOSE, CLOSE},
		data:        map[string]interface{}{},
		valid:       false,
	},
}
