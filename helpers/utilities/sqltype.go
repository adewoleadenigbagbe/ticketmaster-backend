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

func zeroTypeValue[T any]() T {
	var zero T
	return zero
}

func NewNullable[T any](value T, isValid bool) Nullable[T] {
	return Nullable[T]{Val: value, Valid: isValid}
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

func convertToGenericType[T any](value interface{}) (T, error) {
	switch v := value.(type) {
	case []uint8, []uint16, []uint32, []uint64:
		if reflect.TypeOf(v).ConvertibleTo(reflect.TypeOf((*T)(nil)).Elem()) {
			return reflect.ValueOf(value).Convert(reflect.TypeOf((*T)(nil)).Elem()).Interface().(T), nil
		}
	}

	var zero T
	return zero, ErrUnsupportedConversion
}
