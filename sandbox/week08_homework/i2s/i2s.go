package main

import (
	"fmt"
	"reflect"
	"time"
)

func i2s(data any, out any) error {
	// show("i2s params: ", data, out)
	// iterate over out's fields, set each field from map `data`

	var target reflect.Value

	ptr := reflect.ValueOf(out)
	if ptr.Kind() != reflect.Ptr {
		return fmt.Errorf("i2s, out must be a pointer to target container, got %#v", out)
	} else {
		target = ptr.Elem()
	}

	switch target.Kind() {

	case reflect.Bool:
		b, ok := data.(bool)
		if ok {
			target.SetBool(b)
		} else {
			return fmt.Errorf("failed cast to bool from %#v", data)
		}

	case reflect.String:
		str, ok := data.(string)
		if ok {
			target.SetString(str)
		} else {
			return fmt.Errorf("failed cast to string from %#v", data)
		}

	case reflect.Int:
		num, ok := data.(float64) // json bug/feature
		if ok {
			target.SetInt(int64(num))
		} else {
			return fmt.Errorf("failed cast to float64 from %#v", data)
		}

	case reflect.Slice:
		lst, ok := data.([]interface{})
		if ok {
			for _, lstValue := range lst {
				item := reflect.New(target.Type().Elem())
				err := i2s(lstValue, item.Interface())
				if err == nil {
					target.Set(reflect.Append(target, item.Elem()))
				} else {
					return err // fmt.Errorf("failed to process slice element %d: %s", i, err)
				}
			}
		} else {
			return fmt.Errorf("failed cast to []interface{} from %#v", data)
		}

	case reflect.Struct:
		// struct decoded only from map
		dict, ok := data.(map[string]interface{})
		if ok {
			// for each target field
			for i := 0; i < target.NumField(); i++ {
				fieldName := target.Type().Field(i).Name
				fieldValue, ok := dict[fieldName]
				if ok {
					// recursion
					if err := i2s(fieldValue, target.Field(i).Addr().Interface()); err != nil {
						return err // fmt.Errorf("i2s, failed to decode field `%s`: %e", fieldName, err)
					}
				} else {
					// probably should just skip this field
					return fmt.Errorf("i2s, field `%s` not found in given map %#v", fieldName, dict)
				}
			}
		} else {
			return fmt.Errorf("i2s, data must be a map string-to-any, got %#v", data)
		}

	default:
		return fmt.Errorf("i2s, unsupportd type: %s", target.Kind())
	}

	return nil
}

// func userInput(msg string) (res string, err error) {
// 	show(msg)
// 	if n, e := fmt.Scanln(&res); n != 1 || e != nil {
// 		return "", e
// 	}
// 	return res, nil
// }

// func panicOnError(msg string, err error) {
// 	if err != nil {
// 		panic(msg + ": " + err.Error())
// 	}
// }

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	return time.Now().UTC().Format(RFC3339Milli)
}

// show writes message to standard output. Message combined from prefix msg and slice of arbitrary arguments
func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		// line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
