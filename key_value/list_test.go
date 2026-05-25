package key_value

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// We won't test the requests.
// The requests are tested in the controllers
// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestListQueue struct {
	suite.Suite
	list *List
}

// SetupTest
// Setup checks the New() functions
// Setup checks ToMap() functions
func (suite *TestListQueue) SetupTest() {

	list := NewList()
	suite.list = list

	suite.Require().True(list.IsEmpty())
	suite.Require().False(list.IsFull())
	suite.Require().Zero(list.Len())
}

func (suite *TestListQueue) TestAddGet() {
	type Item struct {
		param1 string
		param2 uint64
	}
	// This type of data can not be added if the first
	// element was added
	type InvalidItem struct {
		param1 string
		param2 uint64
	}
	sample := Item{param1: "hello", param2: uint64(0)}
	err := suite.list.Add(uint64(1), sample)
	suite.Require().NoError(err)
	suite.Require().EqualValues(suite.list.Len(), 1)
	suite.Require().False(suite.list.IsFull())
	suite.Require().False(suite.list.IsEmpty())

	// the value type are not matching
	// therefore it should fail
	invalidSample := InvalidItem{param1: "hello", param2: uint64(0)}
	err = suite.list.Add(uint64(2), invalidSample)
	suite.Require().Error(err)
	suite.Require().EqualValues(suite.list.Len(), 1)

	// invalid type
	// already added by value, now pointer type is not valid
	err = suite.list.Add(uint64(3), &sample)
	suite.Require().Error(err)

	// invalid kv type
	err = suite.list.Add(5, sample)
	suite.Require().Error(err)

	// kv value already exists
	err = suite.list.Add(uint64(1), sample)
	suite.Require().Error(err)

	// kv can not be a pointer
	key := uint64(6)
	err = suite.list.Add(&key, sample)
	suite.Require().Error(err)

	// get the data
	list := suite.list.List()
	listSample := list[uint64(1)].(Item)
	suite.Require().EqualValues(sample, listSample)

	// should be successful
	returnedSample, err := suite.list.Get(uint64(1))
	suite.Require().NoError(err)
	suite.Require().EqualValues(sample, returnedSample)

	// should fail since kv doesn't exist
	_, err = suite.list.Get(uint64(10))
	suite.Require().Error(err)

	// should fail since kv type is invalid
	_, err = suite.list.Get(1)
	suite.Require().Error(err)

}

func (suite *TestListQueue) TestListLimit() {
	newList := NewList()

	// index till QUEUE_LENGTH - 2
	for i := uint(0); i < DefaultCap; i++ {
		err := newList.Add(i, i*2)
		suite.Require().NoError(err)
	}

	suite.Require().True(newList.IsFull())
	suite.Require().Equal(newList.Len(), DefaultCap)

	// can not add when the new list is full
	err := newList.Add(DefaultCap, DefaultCap*2)
	suite.Require().Error(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestList(t *testing.T) {
	suite.Run(t, new(TestListQueue))
}
