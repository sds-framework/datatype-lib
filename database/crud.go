package database

import (
	"github.com/noPerfection/datatype"
)

// Crud interface adds the database CRUD operations to the data struct.
//
// The interface that it accepts is the *remote.ClientSocket from the
// "github.com/ahmetson/service-lib/remote" package.
type Crud interface {
	// Update the parameters by int flag. It calls UPDATE command
	Update(any, uint8) error
	// Exist in the database or not. It calls EXIST command
	Exist(any) bool

	// Insert into the database. It calls INSERT command
	Insert(any) error
	// Select selects the single row from the database. It calls SELECT_ROW command
	Select(any) error

	// SelectAll selects the multiple rows from the database. It calls SELECT_ALL without WHERE clause of query.
	//
	// Result is then put to the second argument
	SelectAll(any, any) error

	// SelectAllByCondition returns structs from database to the second argument.
	// The database query should match to the condition.
	//
	// It calls SELECT_ALL with WHERE clause
	SelectAllByCondition(any, datatype.KeyValue, any) error // uses SELECT_ROW
}
