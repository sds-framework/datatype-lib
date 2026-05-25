// Package key_value defines the custom map and its additional functions.
//
// The package defines two different data types:
//   - [KeyValue] is the map where the kv is a string, and the value could be anything.
//     It defines additional functions that return the value converted to the desired type.
//   - [List] is the list of elements but based on the map.
//     For the user, the list acts as the array.
//     However, internally it uses a map for optimization.
package key_value

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/sds-framework/datatype"
)

// KeyValue is the golang's map with the functions.
// Important to know that, no value could be `nil`.
type KeyValue map[string]interface{}

// New empty KeyValue is created
func New() KeyValue {
	return map[string]interface{}{}
}

// NewFromString converts the s string with a json decoder into the kv value
func NewFromString(s string) (KeyValue, error) {
	var keyValue KeyValue

	decoder := json.NewDecoder(strings.NewReader(s))
	decoder.UseNumber()

	if err := decoder.Decode(&keyValue); err != nil {
		return nil, fmt.Errorf("json.decoder: '%w'", err)
	}

	err := keyValue.noNilValue()
	if err != nil {
		return nil, fmt.Errorf("value is nil: %w", err)
	}

	return keyValue, nil
}

// NewFromInterface converts the data structure "i" to KeyValue
// In order to do that, it serializes data structure using json
//
// The data structures should define the json variable names
func NewFromInterface(i interface{}) (KeyValue, error) {
	var k KeyValue
	bytes, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("json.marshal %T: '%w'", i, err)
	}
	err = json.Unmarshal(bytes, &k)
	if err != nil {
		return nil, fmt.Errorf("json:unmarshal %s: '%w'", bytes, err)
	}

	nilErr := k.noNilValue()
	if nilErr != nil {
		return nil, fmt.Errorf("value is nil: %w", nilErr)
	}

	return k, nil
}

// Checks that the values are not nil.
func (k KeyValue) noNilValue() error {
	for key, value := range k {
		if value == nil {
			return fmt.Errorf("kv %s is nil", key)
		}

		nestedKv, ok := value.(KeyValue)

		if ok {
			err := nestedKv.noNilValue()
			if err != nil {
				return fmt.Errorf("kv %s nested value nil: %w", key, err)
			}

			continue
		}

		nestedMap, ok := value.(map[string]interface{})

		if ok {
			nestedKv = nestedMap

			err := nestedKv.noNilValue()
			if err != nil {
				return fmt.Errorf("kv %s nested value nil: %w", key, err)
			}
		}
	}

	return nil
}

// It sets the numbers in a string format.
// The string format for the number means a json number
func (k KeyValue) setNumber() {
	for key, value := range k {
		if value == nil {
			continue
		}

		// even if it's a number wrapped as a string,
		// we won't convert it.
		_, ok := value.(string)
		if ok {
			continue
		}

		_, ok = value.(json.Number)
		if ok {
			continue
		}

		bigNum, err := k.BigIntValue(key)
		if err == nil {
			delete(k, key)

			jsonNumber := json.Number(bigNum.String())
			k.Set(key, jsonNumber)
			continue
		}

		floatNum, err := k.Float64Value(key)
		if err == nil {
			delete(k, key)

			jsonNumber := json.Number(strconv.FormatFloat(floatNum, 'G', -1, 64))
			k.Set(key, jsonNumber)
			continue
		}

		num, err := k.Uint64Value(key)
		if err == nil {
			delete(k, key)

			jsonNumber := json.Number(strconv.FormatUint(num, 10))
			k.Set(key, jsonNumber)
			continue
		}

		nestedKv, ok := value.(KeyValue)
		if ok {
			nestedKv.setNumber()

			delete(k, key)
			k.Set(key, nestedKv)
			continue
		}

		nestedMap, ok := value.(map[string]interface{})

		if ok {
			nestedKv = nestedMap
			// ToMap will call setNumber()
			nestedMap = nestedKv.Map()

			delete(k, key)
			k.Set(key, nestedMap)
			continue
		}
	}
}

// Map converts the k to golang's map
func (k KeyValue) Map() map[string]interface{} {
	k.setNumber()
	return k
}

// MapString returns the map with the string values only
func (k KeyValue) MapString() map[string]string {
	converted := k.Map()

	data := map[string]string{}

	for key := range converted {
		value, err := k.StringValue(key)
		if err != nil {
			continue
		}
		data[key] = value
	}

	return data
}

// Bytes serialize k into the series of bytes
func (k KeyValue) Bytes() ([]byte, error) {
	err := k.noNilValue()
	if err != nil {
		return []byte{}, fmt.Errorf("nil value: %w", err)
	}
	k.setNumber()

	bytes, err := json.Marshal(k)
	if err != nil {
		return []byte{}, fmt.Errorf("json.serialize: '%w'", err)
	}

	return bytes, nil
}

// Returns the serialized kv-value as a string
func (k KeyValue) String() string {
	bytes, err := k.Bytes()
	if err != nil {
		return ""
	}

	return string(bytes)
}

// Interface representation of this KeyValue
func (k KeyValue) Interface(i interface{}) error {
	if !datatype.IsPointer(i) {
		return fmt.Errorf("interface wasn't passed by pointer")
	}
	bytes, err := k.Bytes()
	if err != nil {
		return fmt.Errorf("k.ToBytes of %v: '%w'", k, err)
	}
	err = json.Unmarshal(bytes, i)
	if err != nil {
		return fmt.Errorf("json.deserialize(%s to %T): '%w'", bytes, i, err)
	}

	return nil
}

// Set the parameter in KeyValue
func (k KeyValue) Set(key string, value interface{}) KeyValue {
	k[key] = value

	return k
}

func (k KeyValue) Exist(key string) (exists bool) {
	_, exists = k[key]
	return
}

// Uint64Value returns the parameter as an uint64
func (k KeyValue) Uint64Value(key string) (uint64, error) {
	if !k.Exist(key) {
		return 0, fmt.Errorf("not exist")
	}
	raw := k[key]
	if raw == nil {
		return 0, fmt.Errorf("kv %s is nil", key)
	}

	pureValue, ok := raw.(uint64)
	if ok {
		return pureValue, nil
	}
	value, ok := raw.(float64)
	if ok {
		return uint64(value), nil
	}

	jsonValue, ok := raw.(json.Number)
	if ok {
		number, err := strconv.ParseUint(string(jsonValue), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("strconv.ParseUint(%v (type %T) as json number %v): '%w'", raw, raw, jsonValue, err)
		}
		return number, nil
	}

	stringValue, ok := raw.(string)
	if !ok {
		return 0, fmt.Errorf("'%s' parameter type %T, can not convert to number", key, raw)
	}
	number, err := strconv.ParseUint(stringValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseUint string %v (original: %v): '%w'", stringValue, raw, err)
	}

	return number, nil
}

// Float64Value extracts the float number
func (k KeyValue) Float64Value(key string) (float64, error) {
	if !k.Exist(key) {
		return 0, fmt.Errorf("not exist")

	}
	raw := k[key]
	if raw == nil {
		return 0, fmt.Errorf("kv %s is nil", key)
	}

	pureValue, ok := raw.(float64)
	if ok {
		return pureValue, nil
	}
	value, ok := raw.(json.Number)
	if ok {
		v, err := value.Float64()
		if err != nil {
			return 0, fmt.Errorf("json.Number.Float64() of %v (original: %v): '%w'", value, raw, err)
		}
		return v, nil
	}
	stringValue, ok := raw.(string)
	if !ok {
		return 0, fmt.Errorf("'%s' parameter type %T, can not convert to number", key, raw)
	}
	number, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseUint string %v (original: %v): '%w'", stringValue, raw, err)
	}

	return number, nil
}

// BoolValue extracts the value as boolean
func (k KeyValue) BoolValue(key string) (bool, error) {
	if !k.Exist(key) {
		return false, fmt.Errorf("not exist")
	}
	raw := k[key]
	if raw == nil {
		return false, fmt.Errorf("kv %s is nil", key)
	}

	pureValue, ok := raw.(bool)
	if ok {
		return pureValue, nil
	}

	return false, fmt.Errorf("'%s' parameter type %T, can not convert to boolean", key, raw)
}

// BigIntValue extracts the value as the parsed large number. Use this if the number size is more than 64 bits.
func (k KeyValue) BigIntValue(key string) (*big.Int, error) {
	if !k.Exist(key) {
		return nil, fmt.Errorf("not exist")
	}
	raw := k[key]
	if raw == nil {
		return nil, fmt.Errorf("kv %s is nil", key)
	}

	value, ok := raw.(json.Number)
	if !ok {
		return nil, fmt.Errorf("json.Number: '%s' parameter type %T", key, raw)
	}

	number, ok := big.NewInt(0).SetString(string(value), 10)
	if !ok {
		return nil, fmt.Errorf("math.ParseBig256 failed to parse %s from '%s'", key, value)
	}

	return number, nil
}

// StringValue returns the parameter as a string
func (k KeyValue) StringValue(key string) (string, error) {
	if !k.Exist(key) {
		return "", fmt.Errorf("not exist")
	}
	raw := k[key]
	if raw == nil {
		return "", fmt.Errorf("kv %s is nil", key)
	}

	value, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("%s parameter type %T, can not convert to string", key, raw)
	}

	return value, nil
}

// StringsValue returns the list of strings
func (k KeyValue) StringsValue(key string) ([]string, error) {
	if !k.Exist(key) {
		return nil, fmt.Errorf("not exist")
	}
	raw := k[key]

	values, ok := raw.([]interface{})
	if !ok {
		readyList, ok := raw.([]string)
		if !ok {
			return nil, fmt.Errorf("'%s' parameter type %T, can not convert to string list", key, raw)
		} else {
			return readyList, nil
		}
	}

	list := make([]string, len(values))
	for i, rawValue := range values {
		v, ok := rawValue.(string)
		if !ok {
			return nil, fmt.Errorf("parameter %s[%d] type is %T, can not convert to string %v", key, i, rawValue, rawValue)
		}

		list[i] = v
	}

	return list, nil
}

// NestedListValue returns the parameter as a slice of map:
//
// []key_value.KeyValue
func (k KeyValue) NestedListValue(key string) ([]KeyValue, error) {
	if !k.Exist(key) {
		return nil, fmt.Errorf("not exist")
	}
	raw := k[key]

	values, ok := raw.([]interface{})
	if !ok {
		readyList, ok := raw.([]KeyValue)
		if !ok {
			return nil, fmt.Errorf("'%s' parameter type %T, can not convert to kv-value list", key, raw)
		} else {
			return readyList, nil
		}
	}

	list := make([]KeyValue, len(values))
	for i, rawValue := range values {
		v, ok := rawValue.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("parameter %s[%d] type is %T, can not convert to kv-value %v", key, i, rawValue, rawValue)
		}

		list[i] = v
	}

	return list, nil
}

// NestedValue returns the parameter as a KeyValue
func (k KeyValue) NestedValue(key string) (KeyValue, error) {
	if !k.Exist(key) {
		return nil, fmt.Errorf("not exist")
	}
	raw := k[key]
	if raw == nil {
		return nil, fmt.Errorf("kv %s is nil", key)
	}

	value, ok := raw.(KeyValue)
	if ok {
		err := value.noNilValue()
		if err != nil {
			return nil, fmt.Errorf("kv %s is nil", key)
		}

		return value, nil
	}

	rawMap, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'%s' parameter type %T, can not convert to kv-value", key, raw)
	}

	var nestedKv KeyValue = rawMap
	err := nestedKv.noNilValue()
	if err != nil {
		return nil, fmt.Errorf("kv %s is nil", key)
	}

	return nestedKv, nil
}
