package boule

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T) {

	for _, test := range testCases {
		t.Run(fmt.Sprintf("testing string %v", test.string), func(t *testing.T) {

			data := NewData()
			assert.NoError(t, data.Add(test.data))

			evaluate, err := NewExpression(test.string)
			if test.valid {
				assert.NoError(t, err)

				result, err := evaluate(data)
				assert.NoError(t, err)

				assert.Equal(t, test.result, result)
				return
			}

			assert.Error(t, err)
		})
	}
}
