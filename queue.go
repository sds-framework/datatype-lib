// Package data_type defines the generic data types used in SDS.
//
// Supported data types are:
//   - Queue is the list where the new element is added to the end,
//     but when an element is taken its taken from the top.
//     Queue doesn't allow addition of any kind of element. All elements should have the same type.
//   - Key_value different kinds of maps
//   - serialize functions to serialize any structure to the bytes and vice versa.
package datatype

import (
	"container/list"
	"fmt"
	"reflect"
)

type Queue struct {
	l           *list.List
	cap         uint
	elementType reflect.Type
}

const QueueCap uint = 10

// NewQueue returns the queue of the elements that could contain
// maximum QUEUE_LENGTH number of elements.
//
// The queue has a function that returns the first element
// by taking it out from the list.
//
// The added elements attached after the last element.
func NewQueue() *Queue {
	return &Queue{
		elementType: nil,
		cap:         QueueCap,
		l:           list.New(),
	}
}

func (q *Queue) Len() uint {
	return uint(q.l.Len())
}

func (q *Queue) IsEmpty() bool {
	return q.l.Len() == 0
}

func (q *Queue) IsFull() bool {
	return uint(q.l.Len()) == q.cap
}

// SetCap updates the queue size if it's possible.
// If the new cap is less that current queue size throws an error.
// Otherwise, queue will lose access to the items.
func (q *Queue) SetCap(newCap uint) error {
	if q.Len() > newCap {
		return fmt.Errorf("trying to set %d as cap, however queue has %d elements", newCap, q.Len())
	}

	q.cap = newCap

	return nil
}

// Cap returns the capacity of the queue
func (q *Queue) Cap() uint {
	return q.cap
}

// Push the element into the queue.
// If the element type is not the same as
// the expected type, then
// It will silently drop it.
// Silently drop if the queue is full
func (q *Queue) Push(item interface{}) {
	if q.IsFull() {
		return
	}
	if q.elementType == nil {
		q.elementType = reflect.TypeOf(item)
		q.l.PushBack(item)
	} else if reflect.TypeOf(item) == q.elementType {
		q.l.PushBack(item)
	}
}

// First returns the first element without removing it from the queue
// If there is no element, then returns nil
func (q *Queue) First() interface{} {
	if q.IsEmpty() {
		return nil
	}
	return q.l.Front().Value
}

// Pop takes the first element from the list and returns it.
// If there is no element in the list, then return nil
func (q *Queue) Pop() interface{} {
	if q.IsEmpty() {
		return nil
	}
	return q.l.Remove(q.l.Front())
}
