package boule

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/victordeleau/boule/internal/prefixtree"
	"testing"
)

func TestParser(t *testing.T) {

	for _, test := range testCases {
		t.Run(fmt.Sprintf("testing string %v", test.string), func(t *testing.T) {

			data := prefixtree.New()
			assert.NoError(t, data.Add(test.data))

			evaluate, err := NewAST(test.string, data)
			if test.valid {
				assert.NoError(t, err)

				result, err := evaluate()
				assert.NoError(t, err)

				assert.Equal(t, test.result, result)
				return
			}

			assert.Error(t, err)
		})
	}
}
