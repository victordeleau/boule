package prefixtree

import "fmt"

// Add identifiers to the prefix tree. Two modes are supported: (key, value) mode and map mode.
// (key, value) mode requires a key of type string, and a value of type bool, string, int, int8, int16, int32, int64,
// uint, uint8, uint16, uint32, uint64, float32, or float64. The key and the value are passed as separate arguments.
// map mode requires a map of type map[string]interface{}, and the key/value pairs must follow the same type restrictions
// as the (key, value) mode.
func (p *Tree) Add(input ...interface{}) error {

	if len(input) == 1 {

		inputMap, ok := input[0].(map[string]interface{})
		if !ok {
			return fmt.Errorf("map mode: input is not of type map[string]interface{}")
		}

		for k, v := range inputMap {
			switch v.(type) {
			case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
				p.add(k, v)
				continue
			default:
				return fmt.Errorf("map mode: 'value' type %T is not supported", v)
			}
		}

		return nil

	} else if len(input) == 2 {

		key, ok := input[0].(string)
		if !ok {
			return fmt.Errorf("(key, value) mode: 'key' must be of type 'string'")
		}

		switch input[1].(type) {
		case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			p.add(key, input[1])
			return nil
		default:
			return fmt.Errorf("(key, value) mode: 'value' type %T is not supported", input[1])
		}

	}

	return fmt.Errorf("invalid input mode: either a (key, value) or a map must be passed as argument")
}
