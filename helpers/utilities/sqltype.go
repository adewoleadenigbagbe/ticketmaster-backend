package utilities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
)

var (
	ErrUnsupportedConversion = errors.New("unsupported type conversion")
)

type Nullable[T any] struct {
	Val   T
	Valid bool
}

func (n *Nullable[T]) Scan(value interface{}) error {
	if value == nil {
		n.Val = zeroTypeValue[T]()
		n.Valid = false
		return nil
	}

	var err error
	n.Val, err = convertToGenericType[T](value)
	if err == nil {
		n.Valid = true
	}
	return err
}

func (n Nullable[T]) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Val, nil
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	n.Val = value
	n.Valid = true
	return nil
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Val)
}

func zeroTypeValue[T any]() T {
	var zero T
	return zero
}

func convertToGenericType[T any](value interface{}) (T, error) {
	switch v := value.(type) {
	case T:
		return v, nil
	case int64:
		// Convert to correspondint int type
		switch t := reflect.Zero(reflect.TypeOf((*T)(nil)).Elem()).Interface().(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			if reflect.TypeOf(t).ConvertibleTo(reflect.TypeOf((*T)(nil)).Elem()) {
				return reflect.ValueOf(value).Convert(reflect.TypeOf((*T)(nil)).Elem()).Interface().(T), nil
			}
		}
	}
	var zero T
	return zero, ErrUnsupportedConversion
}
