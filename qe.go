package qe

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func Marshal(params interface{}) (string, error) {
	pairs := []string{}
	val := reflect.ValueOf(params)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		tag := typ.Field(i).Tag
		name, ok := tag.Lookup("q")
		value := val.Field(i).Interface()

		if ok {
			if tag.Get("required") == "true" {
				if value == nil || value == "" || value == 0 || val.Field(i).Kind() == reflect.Slice && fmt.Sprintf("%v", value) == "[]" {
					return "", newError("field '%s' is required, but the slice is empty", typ.Field(i).Name)
				} else if value == nil {
					return "", newError("field '%s' is required, but the value is nil", typ.Field(i).Name)
				}
			}

			if value != nil && value != "" && (val.Field(i).Kind() == reflect.String || fmt.Sprint(value) != "0") {
				encoded, err := encodeValue(typ.Field(i).Type, value)

				if err != nil {
					return "", err
				}

				if encoded != "" {
					escaped := url.QueryEscape(fmt.Sprintf("%v", encoded))
					pairs = append(pairs, name+"="+escaped)
				}
			}
		}
	}

	return strings.Join(pairs, "&"), nil
}

func encodeValue(typ reflect.Type, value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case bool:
		return fmt.Sprintf("%v", v), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", v), nil
	case float32, float64, complex64, complex128:
		return fmt.Sprintf("%v", v), nil
	case []string:
		return strings.Join(v, ","), nil
	case []bool, []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []float32, []float64, []complex64, []complex128:
		joined := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v)), ","), "[]")
		return joined, nil
	default:
		return "", newError("unable to encode type '%v'", typ)
	}
}

func isSupportedSlice(value interface{}) bool {
	switch value.(type) {
	case []bool, []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []float32, []float64, []complex64, []complex128:
		return true
	default:
		return false
	}
}

func newError(message string, args ...interface{}) error {
	return errors.New("query: " + fmt.Sprintf(message, args...))
}
