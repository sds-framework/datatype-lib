package datatype

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func IsNil(v any) bool {
	return v == nil
}

func IsPointer(body any) bool {
	v := reflect.ValueOf(body)
	return v.Kind() == reflect.Pointer
}

// Serialize wraps the JSON abi interface to the internal data type.
// It's blockchain agnostic.
func Serialize(body any) ([]byte, error) {
	if IsPointer(body) {
		return nil, fmt.Errorf("body was passed by a pointer")
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("json: '%w'", err)
	}
	return bytes, nil
}

// Deserialize the given string to the map, slice
// or struct.
func Deserialize(bytes []byte, body any) error {
	if IsNil(body) {
		return fmt.Errorf("body parameter is a nil")
	}
	if !IsPointer(body) {
		return fmt.Errorf("body wasn't passed by pointer")
	}

	err := json.Unmarshal(bytes, body)

	if err != nil {
		return fmt.Errorf("json: '%w'", err)
	}

	return nil
}
