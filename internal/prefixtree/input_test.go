package prefixtree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTree_AddKeyValue(t *testing.T) {

	t.Run("can add key/value pair", func(t *testing.T) {
		tree := new(Tree)
		assert.NoError(t, tree.AddKeyValue("some_key", 380))
		value, err := tree.Find("some_key")
		if assert.NoError(t, err) {
			if assert.IsType(t, 10, value) {
				assert.Equal(t, 380, value)
			}
		}
	})

	t.Run("rejects unsupported value type", func(t *testing.T) {
		assert.Error(t, new(Tree).AddKeyValue("key", []int{1, 2}))
	})
}

func TestTree_AddMap(t *testing.T) {

	t.Run("can add map as data", func(t *testing.T) {
		tree := new(Tree)
		assert.NoError(t, tree.AddMap(map[string]interface{}{
			"road":   "Wellington",
			"number": 20,
		}))

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
		assert.Error(t, new(Tree).AddMap(map[string]interface{}{
			"index": []int{0, 1, 2},
		}))
	})

	t.Run("map can't index map", func(t *testing.T) {
		assert.Error(t, new(Tree).AddMap(map[string]interface{}{
			"index": map[string]int{"un": 1, "deux": 2},
		}))
	})

	t.Run("rejects reserved keyword in map", func(t *testing.T) {
		assert.Error(t, new(Tree).AddMap(map[string]interface{}{
			"true": 1,
		}))
	})
}

func TestTree_AddStruct(t *testing.T) {

	t.Run("can add struct as data", func(t *testing.T) {
		tree := new(Tree)
		data := struct {
			Road   string `json:"road"`
			Number int    `json:"number"`
		}{
			Road:   "Wellington",
			Number: 20,
		}

		assert.NoError(t, tree.AddStruct(data))

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

		assert.NoError(t, tree.AddStruct(data))

		value, err := tree.Find("road")
		if assert.NoError(t, err) {
			if assert.IsType(t, "string", value) {
				assert.Equal(t, "Wellington", value)
			}
		}

		_, err = tree.Find("number")
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

		assert.NoError(t, tree.AddStruct(data))

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
		assert.NoError(t, tree.AddStruct(struct {
			Index map[string]int `json:"index"`
		}{
			Index: map[string]int{"un": 1, "deux": 2, "trois": 3},
		}))

		_, err := tree.Find("index")
		assert.Error(t, err)
	})

	t.Run("embedded slice are not supported", func(t *testing.T) {
		tree := new(Tree)
		assert.NoError(t, tree.AddStruct(struct {
			Index []int `json:"index"`
		}{
			Index: []int{1, 2, 3},
		}))

		_, err := tree.Find("index")
		assert.Error(t, err)
	})

	t.Run("rejects non-struct input", func(t *testing.T) {
		assert.Error(t, new(Tree).AddStruct("not a struct"))
	})
}

func TestTree_KeyValidation(t *testing.T) {

	t.Run("rejects reserved keyword 'true'", func(t *testing.T) {
		assert.Error(t, new(Tree).AddKeyValue("true", 1))
	})

	t.Run("rejects reserved keyword 'false'", func(t *testing.T) {
		assert.Error(t, new(Tree).AddKeyValue("false", 1))
	})

	t.Run("rejects empty key", func(t *testing.T) {
		assert.Error(t, new(Tree).AddKeyValue("", 1))
	})

	t.Run("rejects key starting with digit", func(t *testing.T) {
		assert.Error(t, new(Tree).AddKeyValue("1abc", 1))
	})

	t.Run("rejects key containing operator characters", func(t *testing.T) {
		for _, key := range []string{"a==b", "a>b", "a<b", "a!", "a&b", "a|b"} {
			assert.Error(t, new(Tree).AddKeyValue(key, 1), "expected error for key %q", key)
		}
	})

	t.Run("rejects key containing parentheses", func(t *testing.T) {
		assert.Error(t, new(Tree).AddKeyValue("foo(bar)", 1))
	})

	t.Run("rejects key containing spaces", func(t *testing.T) {
		assert.Error(t, new(Tree).AddKeyValue("foo bar", 1))
	})

	t.Run("accepts valid identifier with dots and underscores", func(t *testing.T) {
		assert.NoError(t, new(Tree).AddKeyValue("ship.max_speed", 100))
	})
}

func TestTree_AddBackwardCompat(t *testing.T) {

	t.Run("key/value via Add", func(t *testing.T) {
		tree := new(Tree)
		assert.NoError(t, tree.Add("some_key", 380))
		value, err := tree.Find("some_key")
		if assert.NoError(t, err) {
			assert.Equal(t, 380, value)
		}
	})

	t.Run("map via Add", func(t *testing.T) {
		tree := new(Tree)
		assert.NoError(t, tree.Add(map[string]interface{}{"road": "Wellington"}))
		value, err := tree.Find("road")
		if assert.NoError(t, err) {
			assert.Equal(t, "Wellington", value)
		}
	})

	t.Run("struct via Add", func(t *testing.T) {
		tree := new(Tree)
		assert.NoError(t, tree.Add(struct {
			Road string `json:"road"`
		}{Road: "Wellington"}))
		value, err := tree.Find("road")
		if assert.NoError(t, err) {
			assert.Equal(t, "Wellington", value)
		}
	})

	t.Run("wrong arg count via Add", func(t *testing.T) {
		assert.Error(t, new(Tree).Add("a", 1, "extra"))
	})

	t.Run("non-string key via Add", func(t *testing.T) {
		assert.Error(t, new(Tree).Add(3, 380))
	})
}
