package main

import (
	"fmt"
	"reflect"
)

// i2s: written using 'happy path first' style
// iterate over out's fields, set each field from map `data`
func i2s(data any, out any) error {
	var writeBool = func(data any, target reflect.Value) error {
		x, ok := data.(bool)
		if ok {
			target.SetBool(x)
			return nil
		} else {
			return fmt.Errorf("failed cast to bool from %#v", data)
		}
	}
	var writeString = func(data any, target reflect.Value) error {
		x, ok := data.(string)
		if ok {
			target.SetString(x)
			return nil
		} else {
			return fmt.Errorf("failed cast to string from %#v", data)
		}
	}
	var writeInt = func(data any, target reflect.Value) error {
		x, ok := data.(float64) // json bug/feature
		if ok {
			target.SetInt(int64(x))
			return nil
		} else {
			return fmt.Errorf("failed cast to float64 from %#v", data)
		}
	}
	var writeSlice = func(data any, target reflect.Value) error {
		x, ok := data.([]interface{})
		if ok {
			for _, lstValue := range x {
				item := reflect.New(target.Type().Elem())
				err := i2s(lstValue, item.Interface())
				if err == nil {
					target.Set(reflect.Append(target, item.Elem()))
				} else {
					return err // i2s error
				}
			}
		} else {
			return fmt.Errorf("failed cast to []interface{} from %#v", data)
		}
		return nil
	}
	var writeStruct = func(data any, target reflect.Value) error {
		// struct decoded only from map
		x, ok := data.(map[string]interface{})
		if ok {
			// for each target field
			for i := 0; i < target.NumField(); i++ {
				fieldName := target.Type().Field(i).Name
				fieldValue, ok := x[fieldName]
				if ok {
					// recursion
					err := i2s(fieldValue, target.Field(i).Addr().Interface())
					if err != nil {
						return err // i2s error
					}
				} else {
					// probably should just skip this field
					return fmt.Errorf("field `%s` not found in given map %#v", fieldName, x)
				}
			}
		} else {
			return fmt.Errorf("i2s, data must be a map string-to-any, got %#v", data)
		}
		return nil
	}

	var writeDataToElem = func(data any, target reflect.Value) error {
		switch target.Kind() {
		case reflect.Bool:
			return writeBool(data, target)
		case reflect.String:
			return writeString(data, target)
		case reflect.Int:
			return writeInt(data, target)
		case reflect.Slice:
			return writeSlice(data, target)
		case reflect.Struct:
			return writeStruct(data, target)
		default:
			return fmt.Errorf("i2s, unsupportd type: %s", target.Kind())
		}
	}

	ptr := reflect.ValueOf(out)
	if ptr.Kind() == reflect.Ptr { // this check could be done only once for entire 'out' struct
		return writeDataToElem(data, ptr.Elem())
	} else {
		return fmt.Errorf("out must be a pointer to target container, got %#v", out)
	}
}
