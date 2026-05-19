package key_value

import (
	"fmt"
	"reflect"

	"github.com/sds-framework/datatype-lib/data_type"
)

type List struct {
	l         map[interface{}]interface{}
	length    uint
	cap       uint
	keyType   reflect.Type
	valueType reflect.Type
}

// DefaultCap max amount of data that this list could keep
const DefaultCap uint = 1_000_000

// NewList returns a new list of the elements that could contain
// maximum DefaultCap of elements.
//
// The queue has a function that returns the first element
// by taking it out from the list.
//
// The added elements attached after the last element.
func NewList() *List {
	return &List{
		keyType:   nil,
		valueType: nil,
		cap:       DefaultCap,
		length:    0,
		l:         map[interface{}]interface{}{},
	}
}

func (q *List) Len() uint {
	return q.length
}

func (q *List) IsEmpty() bool {
	return q.length == 0
}

func (q *List) IsFull() bool {
	return q.length == q.cap
}

func (q *List) List() map[interface{}]interface{} {
	return q.l
}

// SetCap updates the capacity of the list.
// If the list has more elements than newCap, it throws an error.
func (q *List) SetCap(newCap uint) error {
	if q.length > newCap {
		return fmt.Errorf("list has %d elements. can not set cap to %d", q.length, newCap)
	}

	q.cap = newCap

	return nil
}

// Cap returns the capacity of the list
func (q *List) Cap() uint {
	return q.cap
}

// Add a new element to the queue.
// If the element type is not the same as
// the expected type, then
// It will silently drop it.
// Silently drop if the queue is full
func (q *List) Add(key interface{}, value interface{}) error {
	if q.IsFull() {
		return fmt.Errorf("list is already full")
	}
	if data_type.IsNil(key) {
		return fmt.Errorf("the kv parameter is nil")
	}
	if data_type.IsPointer(key) {
		return fmt.Errorf("the kv was passed by the pointer")
	}
	if data_type.IsNil(value) {
		return fmt.Errorf("the value parameer is nil")
	}

	keyType := reflect.TypeOf(key)
	valueType := reflect.TypeOf(value)

	if q.keyType == nil {
		q.keyType = keyType
		q.valueType = valueType
	} else if _, ok := q.l[key]; ok {
		return fmt.Errorf("the element exists")
	}

	if keyType == q.keyType && valueType == q.valueType {
		q.l[key] = value
		q.length++
		return nil
	}

	return fmt.Errorf(
		"expected kv type %T against %T and expected value type %T against %T",
		q.keyType,
		keyType,
		q.valueType,
		valueType,
	)
}

func (q *List) Exist(key interface{}) bool {
	if data_type.IsNil(key) ||
		data_type.IsPointer(key) ||
		q.IsEmpty() {
		return false
	}

	keyType := reflect.TypeOf(key)
	if keyType != q.keyType {
		return false
	}

	_, ok := q.l[key]
	if !ok {
		return false
	}

	return true
}

// Get the element in the list to the value.
// Pointer should pass the value
func (q *List) Get(key interface{}) (interface{}, error) {
	if data_type.IsNil(key) {
		return nil, fmt.Errorf("the parameter is nil")
	}
	if data_type.IsPointer(key) {
		return nil, fmt.Errorf("the kv was passed by the pointer")
	}
	if q.IsEmpty() {
		return nil, fmt.Errorf("the list is empty")
	}

	keyType := reflect.TypeOf(key)
	if keyType != q.keyType {
		return nil, fmt.Errorf("the data mismatch: expected kv type %T against %T", q.keyType, keyType)
	}

	value, ok := q.l[key]
	if !ok {
		return nil, fmt.Errorf("the element not found")
	}
	return value, nil
}

// GetFirst returns the first added element.
// Returns the kv, value and error if it can not find it.
func (q *List) GetFirst() (interface{}, interface{}, error) {
	for key, value := range q.l {
		return key, value, nil
	}

	return nil, nil, fmt.Errorf("empty list")
}

// Take is a Get, but removes the returned element from the list
func (q *List) Take(key interface{}) (interface{}, error) {
	value, err := q.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get the element")
	}

	delete(q.l, key)
	q.length--

	return value, nil
}

// TakeFirst is a GetFirst, but removes the element from the list
func (q *List) TakeFirst() (interface{}, interface{}, error) {
	key, value, err := q.GetFirst()
	if err != nil {
		return nil, nil, fmt.Errorf("list.GetFirst: %w", err)
	}

	delete(q.l, key)
	q.length--

	return key, value, nil
}
