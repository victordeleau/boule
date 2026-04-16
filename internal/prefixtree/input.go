package prefixtree

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
)

// AddKeyValue adds a single identifier to the prefix tree.
//
// Keys must start with an ASCII letter (a-z, A-Z) and may only contain ASCII letters,
// digits (0-9), underscores, and dots. Reserved keywords "true" and "false" are rejected.
// Supported value types: bool, string, int, int8, int16, int32, int64, uint, uint8,
// uint16, uint32, uint64, float32, float64, and *big.Int.
func (p *Tree) AddKeyValue(key string, value interface{}) error {
	return p.addKeyValue(key, value)
}

// AddMap adds all entries from a map[string]interface{} to the prefix tree.
// Keys and values follow the same rules as AddKeyValue.
func (p *Tree) AddMap(m map[string]interface{}) error {
	for k, v := range m {
		if err := p.addKeyValue(k, v); err != nil {
			return err
		}
	}
	return nil
}

// AddStruct adds fields from a struct to the prefix tree. Field names are derived
// from json struct tags (fields without a json tag are ignored). Nested structs are
// supported via dot notation (e.g. "owner.name"). Slice and map fields are skipped.
func (p *Tree) AddStruct(s interface{}) error {
	fieldMap, err := p.structToJsonFieldMap(s)
	if err != nil {
		return err
	}
	for k, v := range fieldMap {
		if err := p.addKeyValue(k, v); err != nil {
			return err
		}
	}
	return nil
}

// Deprecated: Add is kept for backward compatibility. Use AddKeyValue, AddMap, or AddStruct instead.
func (p *Tree) Add(input ...interface{}) error {
	if len(input) == 1 {
		if m, ok := input[0].(map[string]interface{}); ok {
			return p.AddMap(m)
		}
		return p.AddStruct(input[0])
	}
	if len(input) == 2 {
		key, ok := input[0].(string)
		if !ok {
			return fmt.Errorf("key must be of type string")
		}
		return p.AddKeyValue(key, input[1])
	}
	return fmt.Errorf("expected 1 or 2 arguments, got %d", len(input))
}

var reservedKeywords = map[string]struct{}{
	"true":  {},
	"false": {},
}

// validateKey checks that key is a valid ASCII identifier: non-empty, starts with an
// ASCII letter (a-z, A-Z), contains only ASCII letters/digits/underscores/dots, and
// is not a reserved keyword.
func validateKey(key string) error {
	if len(key) == 0 {
		return fmt.Errorf("key must not be empty")
	}
	if _, reserved := reservedKeywords[key]; reserved {
		return fmt.Errorf("key %q is a reserved keyword", key)
	}
	if !isASCIILetter(key[0]) {
		return fmt.Errorf("key %q must start with an ASCII letter", key)
	}
	for i := 0; i < len(key); i++ {
		c := key[i]
		if !isASCIILetter(c) && !isASCIIDigit(c) && c != '_' && c != '.' {
			return fmt.Errorf("key %q contains invalid character %q", key, c)
		}
	}
	return nil
}

func isASCIILetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isASCIIDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (p *Tree) addKeyValue(key string, value interface{}) error {
	if err := validateKey(key); err != nil {
		return err
	}
	switch value.(type) {
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *big.Int:
		p.add(key, value)
		return nil
	default:
		return fmt.Errorf("'value' type %T is not supported", value)
	}
}

func (p *Tree) structToJsonFieldMap(input interface{}) (map[string]interface{}, error) {

	fieldMap := make(map[string]interface{})

	inputType, inputValue := reflect.TypeOf(input), reflect.ValueOf(input)

	if inputType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	var structField reflect.StructField
	var fieldType reflect.Type
	var fieldValue reflect.Value
	var jsonName string
	var ok bool

	for i := 0; i < inputType.NumField(); i++ {

		structField = inputType.Field(i)

		jsonName, ok = structField.Tag.Lookup("json")
		if !ok {
			continue // ignore field without 'json' tag name
		}
		if split := strings.Split(jsonName, ","); len(split) > 1 {
			jsonName = split[0]
		}

		fieldType, fieldValue = structField.Type, inputValue.Field(i)

		if fieldType.Kind() == reflect.Ptr {
			fieldType, fieldValue = fieldType.Elem(), fieldValue.Elem()
		}

		if fieldType.Kind() == reflect.Struct { // recurse on struct field
			subFieldMap, err := p.structToJsonFieldMap(fieldValue.Interface())
			if err != nil {
				return nil, err
			}
			for k, v := range subFieldMap {
				fieldMap[jsonName+"."+k] = v
			}
			continue
		}

		if fieldType.Kind() == reflect.Slice {
			continue // slice fields are not supported for now
		}

		if fieldType.Kind() == reflect.Map {
			continue // map fields are not supported for now
		}

		fieldMap[jsonName] = fieldValue.Interface()
	}

	return fieldMap, nil
}
