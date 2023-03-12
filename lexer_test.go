package boule

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/victordeleau/boule/prefixtree"
	"testing"
)

func TestLexer(t *testing.T) {

	for _, test := range testCases {
		t.Run(fmt.Sprintf("testing string %s", test.string), func(t *testing.T) {

			lexer := newLexer(test.string, prefixtree.New())

			var token *lexerToken
			output := make([]*lexerToken, 0, 10)
			for {
				if token = lexer.Yield(); token.token == EOF {
					break
				}
				output = append(output, token)
			}

			assert.Equal(t, len(test.tokenStream), len(output))
			for i, token := range output {
				assert.Equal(t, test.tokenStream[i], token.token)
			}
		})
	}
}
