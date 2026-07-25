# noPerfection/datatype

This module consists of the data types used in all go modules such as Queue, KeyValue, Lists.

Additionally it adds utility helper functions that are missing in standard go library.

> License? **Public Domain**

### Queue

First in, First out queue. Available in `data_type.Queue`.
The queue has a cap.

### KeyValue

Stored in the `data_type/key_value.KeyValue`.
The `KeyValue` is the wrapper around `map[string]interface{}`.
The structure adds methods to extract the data with the validation.

The structure also can not contain the null parameters.

All the keys are string.

The values can be:

- `NestedValue` &ndash; nested `KeyValue`
- `NestedListValue` &ndash; list of `KeyValue`.
- `Uint64` &ndash; any natural numbers and zero are converted into go's `uint64` type. **KeyValue negative numbers are represented as Float64**.
- `Float64` &ndash; the number is represented as a float of 64 bits.
- `String` &ndash; a string.
- `Strings` &ndash; a slice of string.
- `BigNumber` &ndash; the number of `big.Int` format.
- `Bool` &ndash; a boolean parameter.

### KeyValueList

Stored in the `data_type/key_value.List`.
The `List` is the `KeyValue` with two conditions:
Keys are interface, not string. 
The key types must be identical.
The value types must be identical.

**The `List` has a cap**.

If the first value that you set is the number key, and struct A value.  
Then all keys must be number, and all values must be struct A.