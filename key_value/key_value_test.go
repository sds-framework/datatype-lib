package key_value

import (
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/suite"
)

// We won't test the requests.
// The requests are tested in the controllers
// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestKeyValueSuite struct {
	suite.Suite
	kv KeyValue
}

// SetupTest
// Setup checks the New() functions
// Setup checks ToMap() functions
func (suite *TestKeyValueSuite) SetupTest() {
	empty := map[string]interface{}{}
	var kv KeyValue = empty
	suite.Require().EqualValues(empty, kv)
	suite.Require().Equal(empty, kv.Map())

	// no null value could be used
	invalidStr := `{"param_1":null,"param_2":"string_value","param_3":{"nested_1":5,"nested_2":"hello"}}`
	_, err := NewFromString(invalidStr)
	suite.Require().Error(err)

	// no null value could be used in the nested values
	invalidStr = `{"param_1":1,"param_2":"string_value","param_3":{"nested_1":5,"nested_2":null}}`
	_, err = NewFromString(invalidStr)
	suite.Require().Error(err)

	// validate the parameters
	str := `{"param_1":2,"param_2":"string_value","param_3":{"nested_1":5,"nested_2":"hello"}}`
	strKv, err := NewFromString(str)
	suite.Require().NoError(err)
	mapKey := strKv.Map()

	var num2 = json.Number("2")
	var num5 = json.Number("5")

	strMap := map[string]interface{}{
		"param_1": num2,
		"param_2": "string_value",
		"param_3": map[string]interface{}{
			"nested_1": num5,
			"nested_2": "hello",
		},
	}
	invalidMap := map[string]interface{}{
		"param_1": 2,
		"param_2": "string_value",
		"param_3": map[string]interface{}{
			"nested_1": uint64(5),
			"nested_2": "hello",
		},
	}

	// one of the parameters is not uint64
	suite.Require().NotEqual(invalidMap, mapKey)
	suite.Require().Equal(strMap, mapKey)

	type Nested struct {
		Nested1 uint64 `json:"nested_1"`
		Nested2 string `json:"nested_2"`
	}
	type Temp struct {
		Param1 uint64 `json:"param_1"`
		Param2 string `json:"param_2"`
		Param3 Nested `json:"param_3"`
	}
	newTemp := Temp{
		Param1: uint64(2),
		Param2: "string_value",
		Param3: Nested{
			Nested1: uint64(5),
			Nested2: "hello",
		},
	}
	interfaceKv, err := NewFromInterface(newTemp)
	suite.Require().NoError(err)
	// The number type in the kv is json.Number
	// But in the temp it's not
	suite.Require().NotEqual(strKv, interfaceKv)
	suite.Require().EqualValues(mapKey, interfaceKv.Map())

	// invalid, the parameters are as is in the struct
	// it misses `json:<param>`
	type InvalidTemp struct {
		Param1 uint64
		Param2 string `json:"param_2"`
		Param3 Nested `json:"param_3"`
	}
	invalidTemp := InvalidTemp{
		Param1: uint64(2),
		Param2: "string_value",
		Param3: Nested{
			Nested1: uint64(5),
			Nested2: "hello",
		},
	}
	interfaceKv, err = NewFromInterface(invalidTemp)
	suite.Require().NoError(err)
	suite.Require().NotEqual(strKv, interfaceKv)

	// Any number is returned as uint64
	type TempUint struct {
		Param1 uint   `json:"param_1"`
		Param2 string `json:"param_2"`
		Param3 Nested `json:"param_3"`
	}
	uintTemp := TempUint{
		Param1: uint(2),
		Param2: "string_value",
		Param3: Nested{
			Nested1: uint64(5),
			Nested2: "hello",
		},
	}
	interfaceKv, err = NewFromInterface(uintTemp)
	suite.Require().NoError(err)
	param1, err := interfaceKv.Uint64Value("param_1")
	suite.Require().NoError(err)
	suite.Require().Equal(uint64(2), param1)

	suite.kv = strKv
}

func (suite *TestKeyValueSuite) TestToString() {
	str := `{"param_1":2,"param_2":"string_value","param_3":{"nested_1":5,"nested_2":"hello"}}`
	kvStr := suite.kv.String()
	suite.Require().Equal(str, kvStr)

	// nil parameter is not allowed
	nilKv := KeyValue(map[string]interface{}{"nil_param": nil})
	kvStr = nilKv.String()
	suite.Require().Empty(kvStr)

	// Empty parameter is okay
	emptyParam := KeyValue(map[string]interface{}{"empty_param": ""})
	kvStr = emptyParam.String()
	suite.Require().NotEmpty(kvStr)
}

func (suite *TestKeyValueSuite) TestToInterface() {
	type Nested struct {
		Nested1 uint64 `json:"nested_1"`
		Nested2 string `json:"nested_2"`
	}
	type Temp struct {
		Param1 uint64 `json:"param_1"`
		Param2 string `json:"param_2"`
		Param3 Nested `json:"param_3"`
	}
	var newTemp Temp
	err := suite.kv.Interface(&newTemp)
	suite.Require().NoError(err)

	// Can not convert to the scalar format,
	// But it will be empty
	// since it's not passed by a pointer
	var invalidTemp string
	err = suite.kv.Interface(invalidTemp)
	suite.Require().Error(err)

	// Can convert with the wrong type
	// But check it in the struct
	type InvalidTemp struct {
		Param1 uint64
		Param2 string `json:"param_2"`
		Param3 Nested `json:"param_3"`
	}
	var noJsonTemp InvalidTemp
	err = suite.kv.Interface(&noJsonTemp)
	suite.Require().NoError(err)

	// Can convert to another type
	// with the invalid parameter type.
	// The map's param_2 is a string.
	type InvalidType struct {
		Param1 uint64 `json:"param_1"`
		Param2 uint64 `json:"param_2"`
		Param3 Nested `json:"param_3"`
	}
	var hasMoreTemp InvalidType
	err = suite.kv.Interface(&hasMoreTemp)
	suite.Require().Error(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestKeyValue(t *testing.T) {
	suite.Run(t, new(TestKeyValueSuite))
}
