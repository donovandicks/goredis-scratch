package resp

import (
	"strconv"
)

type ValueType string

const (
	Array   ValueType = "array"
	Bulk              = "bulk"
	Error             = "error"
	Integer           = "integer"
	String            = "string"
	Null              = "null"
)

type Value struct {
	Typ   ValueType
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

func NewStringValue(s string) Value {
	return Value{Typ: String, Str: s}
}

func NewBulkValue(s string) Value {
	return Value{Typ: Bulk, Bulk: s}
}

func NewErrorValue(err string) Value {
	return Value{Typ: Error, Str: err}
}

func NewNullValue() Value { return Value{Typ: Null} }

func (v Value) Marshal() []byte {
	switch v.Typ {
	case Array:
		return v.marshalArray()
	case Bulk:
		return v.marshalBulk()
	case String:
		return v.marshalString()
	case Null:
		return v.marshalNull()
	case Error:
		return v.marshalErr()
	default:
		return []byte{}
	}
}

func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STR)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalArray() []byte {
	arrlen := len(v.Array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(arrlen)...)
	bytes = append(bytes, '\r', '\n')

	for _, item := range v.Array {
		bytes = append(bytes, item.Marshal()...)
	}

	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

func (v Value) marshalErr() []byte {
	var bytes []byte
	bytes = append(bytes, ERR)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}
