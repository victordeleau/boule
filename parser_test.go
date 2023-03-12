package boule

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/victordeleau/boule/prefixtree"
	"testing"
)

func TestParser(t *testing.T) {

	for _, test := range testCases {
		t.Run(fmt.Sprintf("testing string %v", test.string), func(t *testing.T) {

			_, err := newAST(newLexer(test.string, prefixtree.New()))
			if test.valid {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidSyntax)
		})
	}
}
