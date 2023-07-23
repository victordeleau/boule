package boule

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer(t *testing.T) {

	for _, test := range testCases {
		t.Run(fmt.Sprintf("testing string %s", test.string), func(t *testing.T) {

			lexer := newLexer(test.string)

			var token *lexerTokenWithPosition
			output := make([]*lexerTokenWithPosition, 0, 10)
			for {
				if token = lexer.Yield(); token.token == EOF {
					break
				}
				output = append(output, token)
			}

			for _, t := range output {
				fmt.Printf("%v\n", t.value)
			}
			if assert.Equal(t, len(test.tokenStream), len(output)) {
				for i, token := range output {
					assert.Equal(t, test.tokenStream[i], token.token)
				}
			}
		})
	}
}
