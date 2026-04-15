package prefixtree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTree_Add(t *testing.T) {

	t.Run("test key/value pair", func(t *testing.T) {

		t.Run("can add key/value pair", func(t *testing.T) {

			tree := new(Tree)
			assert.NoError(t, tree.Add("some_key", 380))
			value, err := tree.Find("some_key")
			if assert.NoError(t, err) {
				if assert.IsType(t, 10, value) {
					assert.Equal(t, 380, value)
				}
			}
		})

		t.Run("key must be a string", func(t *testing.T) {

			tree := new(Tree)
			assert.Error(t, tree.Add(3, 380))
		})

		t.Run("can't pass more than 2 arguments", func(t *testing.T) {

			tree := new(Tree)
			assert.Error(t, tree.Add("some_key", 380, "oups"))
		})
	})

	t.Run("test key validation", func(t *testing.T) {

		t.Run("rejects reserved keyword 'true'", func(t *testing.T) {
			assert.Error(t, new(Tree).Add("true", 1))
		})

		t.Run("rejects reserved keyword 'false'", func(t *testing.T) {
			assert.Error(t, new(Tree).Add("false", 1))
		})

		t.Run("rejects empty key", func(t *testing.T) {
			assert.Error(t, new(Tree).Add("", 1))
		})

		t.Run("rejects key starting with digit", func(t *testing.T) {
			assert.Error(t, new(Tree).Add("1abc", 1))
		})

		t.Run("rejects key containing operator characters", func(t *testing.T) {
			for _, key := range []string{"a==b", "a>b", "a<b", "a!", "a&b", "a|b"} {
				assert.Error(t, new(Tree).Add(key, 1), "expected error for key %q", key)
			}
		})

		t.Run("rejects key containing parentheses", func(t *testing.T) {
			assert.Error(t, new(Tree).Add("foo(bar)", 1))
		})

		t.Run("rejects key containing spaces", func(t *testing.T) {
			assert.Error(t, new(Tree).Add("foo bar", 1))
		})

		t.Run("accepts valid identifier with dots and underscores", func(t *testing.T) {
			tree := new(Tree)
			assert.NoError(t, tree.Add("ship.max_speed", 100))
		})

		t.Run("rejects reserved keyword in map mode", func(t *testing.T) {
			assert.Error(t, new(Tree).Add(map[string]interface{}{
				"true": 1,
			}))
		})
	})

	t.Run("test map", func(t *testing.T) {

		t.Run("can add map as data", func(t *testing.T) {

			tree := new(Tree)

			data := map[string]interface{}{
				"road":   "Wellington",
				"number": 20,
			}

			assert.NoError(t, tree.Add(data))

			value, err := tree.Find("road")
			if assert.NoError(t, err) {
				if assert.IsType(t, "string", value) {
					assert.Equal(t, "Wellington", value)
				}
			}

			value, err = tree.Find("number")
			if assert.NoError(t, err) {
				if assert.IsType(t, 0, value) {
					assert.Equal(t, 20, value)
				}
			}
		})

		t.Run("map can't index slice", func(t *testing.T) {

			assert.Error(t, new(Tree).Add(map[string]interface{}{
				"index": []int{0, 1, 2},
			}))
		})

		t.Run("map can't index map", func(t *testing.T) {

			assert.Error(t, new(Tree).Add(map[string]interface{}{
				"index": map[string]int{"un": 1, "deux": 2},
			}))
		})
	})

	t.Run("test struct", func(t *testing.T) {

		t.Run("can add struct as data", func(t *testing.T) {

			tree := new(Tree)

			data := struct {
				Road   string `json:"road"`
				Number int    `json:"number"`
			}{
				Road:   "Wellington",
				Number: 20,
			}

			assert.NoError(t, tree.Add(data))

			value, err := tree.Find("road")
			if assert.NoError(t, err) {
				if assert.IsType(t, "string", value) {
					assert.Equal(t, "Wellington", value)
				}
			}

			value, err = tree.Find("number")
			if assert.NoError(t, err) {
				if assert.IsType(t, 0, value) {
					assert.Equal(t, 20, value)
				}
			}
		})

		t.Run("fields without json tag are ignored", func(t *testing.T) {

			tree := new(Tree)

			data := struct {
				Road   string `json:"road"`
				Number int
			}{
				Road:   "Wellington",
				Number: 20,
			}

			assert.NoError(t, tree.Add(data))

			value, err := tree.Find("road")
			if assert.NoError(t, err) {
				if assert.IsType(t, "string", value) {
					assert.Equal(t, "Wellington", value)
				}
			}

			value, err = tree.Find("number")
			assert.Error(t, err)
		})

		t.Run("embedded structs are supported", func(t *testing.T) {

			type Owner struct {
				Name string `json:"name,omitempty"`
			}

			tree := new(Tree)

			data := struct {
				Road  string `json:"road"`
				Owner Owner  `json:"owner"`
			}{
				Road: "Wellington",
				Owner: Owner{
					Name: "Rodolph",
				},
			}

			assert.NoError(t, tree.Add(data))

			value, err := tree.Find("road")
			if assert.NoError(t, err) {
				if assert.IsType(t, "string", value) {
					assert.Equal(t, "Wellington", value)
				}
			}

			value, err = tree.Find("owner.name")
			if assert.NoError(t, err) {
				if assert.IsType(t, "string", value) {
					assert.Equal(t, "Rodolph", value)
				}
			}
		})

		t.Run("embedded maps are not supported", func(t *testing.T) {

			tree := new(Tree)

			assert.NoError(t, tree.Add(struct {
				Index map[string]int `json:"index"`
			}{
				Index: map[string]int{"un": 1, "deux": 2, "trois": 3},
			}))

			_, err := tree.Find("index")
			assert.Error(t, err)
		})

		t.Run("embedded slice are not supported", func(t *testing.T) {

			tree := new(Tree)

			assert.NoError(t, tree.Add(struct {
				Index []int `json:"index"`
			}{
				Index: []int{1, 2, 3},
			}))

			_, err := tree.Find("index")
			assert.Error(t, err)
		})
	})
}
