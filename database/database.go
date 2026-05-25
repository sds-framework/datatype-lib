// Package database keeps the utility
// functions that converts
// database type to
// internal golang type
//
// It also keeps the interface for structs to implement database connection
package database

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/sds-framework/datatype"
)

// Returns the type of database type
// that matches to the golang type
//
// If the data type wasn't detected, then
// it returns an empty result.
func detectType(databaseType *sql.ColumnType) string {
	switch databaseType.DatabaseTypeName() {
	case "VARCHAR":
		return "string"
	case "JSON":
		return "json"
	case "SMALLINT":
		return "int64"
	case "BIGINT":
		return "int64"
	case "UNSIGNED SMALLINT":
		return "uint64"
	case "UNSIGNED BIGINT":
		return "uint64"
	}
	return ""
}

// SetValue sets the value into kv KeyValue.
// Before setting, the function converts the value into the desired
// golang parameter
func SetValue(kv datatype.KeyValue, databaseType *sql.ColumnType, raw any) error {
	golangType := detectType(databaseType)
	if golangType == "" {
		return fmt.Errorf("unsupported database type %s", databaseType.DatabaseTypeName())
	}

	switch golangType {
	case "string":
		if raw == nil {
			kv.Set(databaseType.Name(), "")
			return nil
		}
		value, ok := raw.(string)
		if !ok {
			bytes, ok := raw.([]byte)
			if !ok {
				return fmt.Errorf("couldn't convert %v of type %T into 'string'", raw, raw)
			}
			kv.Set(databaseType.Name(), string(bytes))
			return nil
		}
		kv.Set(databaseType.Name(), value)
		return nil
	case "json":
		if raw == nil {
			kv.Set(databaseType.Name(), []byte{})
			return nil
		}

		value, ok := raw.([]byte)
		if !ok {
			return fmt.Errorf("database value is expected to be '[]byte', but value %v of type %T", raw, raw)
		}
		kv.Set(databaseType.Name(), datatype.AddJsonPrefix(value))
		return nil
	case "int64":
		if raw == nil {
			kv.Set(databaseType.Name(), int64(0))
			return nil
		}
		value, ok := raw.(int64)
		if !ok {
			bytes, ok := raw.([]byte)
			if !ok {
				return fmt.Errorf("couldn't convert %v of type %T into 'int64'", raw, raw)
			}
			data, err := strconv.ParseInt(string(bytes), 10, 64)
			if err != nil {
				return fmt.Errorf("strconv.ParseInt: %w", err)
			}
			kv.Set(databaseType.Name(), data)
			return nil
		}
		kv.Set(databaseType.Name(), value)
		return nil
	case "uint64":
		if raw == nil {
			kv.Set(databaseType.Name(), uint64(0))
			return nil
		}
		value, ok := raw.(uint64)
		if !ok {
			bytes, ok := raw.([]byte)
			if !ok {
				newValue, ok := raw.(int64)
				if !ok {
					return fmt.Errorf("couldn't convert %v of type %T into 'uint64'", raw, raw)
				}
				kv.Set(databaseType.Name(), newValue)
				return nil
			}
			data, err := strconv.ParseUint(string(bytes), 10, 64)
			if err != nil {
				return fmt.Errorf("strconv.ParseUint: %w", err)
			}
			kv.Set(databaseType.Name(), data)
			return nil
		}
		kv.Set(databaseType.Name(), value)
		return nil
	}

	return fmt.Errorf("no switch/case for setting value into KeyValue for %s field of %s type", databaseType.Name(), golangType)
}
