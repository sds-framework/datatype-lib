package datatype

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// We won't test the requests.
// The requests are tested in the controllers
// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestSerializerSuite struct {
	suite.Suite
}

// SetupTest prepares the test.
// Setup checks the New() functions
// Setup checks ToMap() functions
func (suite *TestSerializerSuite) SetupTest() {
}

func (suite *TestSerializerSuite) TestSerialization() {
	type Item struct {
		Param1 string `json:"param_1"`
		Param2 uint64 `json:"param_2"`
	}
	sample := Item{Param1: "hello", Param2: uint64(5)}

	body, err := Serialize(sample)
	suite.Require().NoError(err)

	expected := `{"param_1":"hello","param_2":5}`
	suite.Require().EqualValues(expected, string(body))

	var newSample Item
	err = Deserialize(body, &newSample)
	suite.Require().NoError(err)
	suite.Require().EqualValues(newSample, sample)

	// try to serialize without passing the reference
	var noRef Item
	err = Deserialize(body, noRef)
	suite.Require().Error(err)
	suite.Require().Empty(noRef)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSerializer(t *testing.T) {
	suite.Run(t, new(TestSerializerSuite))
}
