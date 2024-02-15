package utilities

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

func SetDefaults[T any](t *T) error {
	tType := reflect.TypeOf(*t)
	if tType.Kind() != reflect.Struct {
		return errors.New("default allowed only allowed on struct")
	}

	p := reflect.ValueOf(t).Elem()

	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		if value, ok := field.Tag.Lookup("default"); ok {
			switch field.Type.Kind() {
			case reflect.String:
				if p.Field(i).String() == "" {
					p.Field(i).SetString(value)
				}
			case reflect.Int:
				if p.Field(i).Int() == 0 {
					intValue, err := strconv.Atoi(value)
					if err != nil {
						return err
					}
					if !p.Field(i).OverflowInt(int64(intValue)) {
						p.Field(i).SetInt(int64(intValue))
					}
				}
			case reflect.Float32, reflect.Float64:
				if p.Field(i).Float() == 0 {
					floatValue, err := strconv.ParseFloat(value, 64)
					if err != nil {
						return err
					}
					if !p.Field(i).OverflowFloat(floatValue) {
						p.Field(i).SetFloat(floatValue)
					}
				}
			case reflect.Bool:
				if !p.Field(i).Bool() {
					b, err := strconv.ParseBool(value)
					if err != nil {
						return err
					}
					if !b {
						return errors.New("cannot default bool value to be false")
					}
					p.Field(i).SetBool(b)
				}

			case reflect.Struct:
				//TODO: time datatype is considered for now, in future other struct will be deeply reflected
				if p.Field(i).IsZero() {
					t, err := time.Parse(time.RFC3339, value)
					if err != nil {
						return err
					}
					p.Field(i).Set(reflect.ValueOf(t))
				}
			}
		}
	}

	return nil
}
