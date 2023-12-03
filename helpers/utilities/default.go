package utilities

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func SetDefaults[T any](t *T) error {
	tType := reflect.TypeOf(*t)
	if tType.Kind() != reflect.Struct {
		fmt.Println(tType.Kind())
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
					if intValue, err := strconv.Atoi(value); err == nil {
						if !p.Field(i).OverflowInt(int64(intValue)) {
							p.Field(i).SetInt(int64(intValue))
						}

					}
				}
			case reflect.Float32, reflect.Float64:
				if p.Field(i).Float() == 0.0 {
					if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
						if !p.Field(i).OverflowFloat(floatValue) {
							p.Field(i).SetFloat(floatValue)
						}
					}
				}
			}
		}
	}

	return nil
}
