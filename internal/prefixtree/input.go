package prefixtree

import (
	"fmt"
	"reflect"
	"strings"
)

var errBadInput = fmt.Errorf("input is neither a struct{}, a map[string]interface{}, nor a key/value pair")

// Add identifiers to the prefix tree. Two modes are supported: (key, value) mode and map mode.
// (key, value) mode requires a key of type string, and a value of type bool, string, int, int8, int16, int32, int64,
// uint, uint8, uint16, uint32, uint64, float32, or float64. The key and the value are passed as separate arguments.
// map mode requires a map of type map[string]interface{}, and the key/value pairs must follow the same type restrictions
// as the (key, value) mode.
func (p *Tree) Add(input ...interface{}) error {

	var err error

	if len(input) == 1 { // either a map or a struct

		fieldMap, isMap := input[0].(map[string]interface{})

		if !isMap { // if not map, then must be a struct
			if fieldMap, err = p.structToJsonFieldMap(input[0]); err != nil {
				return errBadInput
			}

			for k, v := range fieldMap {
				if err = p.addKeyValue(k, v); err != nil {
					return err
				}
			}
		}

		for k, v := range fieldMap {
			if err = p.addKeyValue(k, v); err != nil {
				return err
			}
		}
		return nil

	} else if len(input) == 2 { // then must be a key/value pair

		key, ok := input[0].(string)
		if !ok {
			return fmt.Errorf("(key, value) mode: 'key' must be of type 'string'")
		}

		if err = p.addKeyValue(key, input[1]); err != nil {
			return err
		}
		return nil
	}

	return errBadInput
}

func (p *Tree) addKeyValue(key string, value interface{}) error {
	switch value.(type) {
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
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
