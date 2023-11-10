package fields

import (
	"reflect"
	"strings"
)

func Extract[T any](m any, fields ...any) T {
	result := ExtractdAny(m, fields...)
	if ret, ok := result.(T); ok {
		return ret
	}
	var zero T
	return zero
}

func ExtractdAny(m any, fields ...any) any {
	if len(fields) == 0 {
		return m
	}

	v := reflect.ValueOf(m)
	for idx, f := range fields {
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		} else if v.Kind() == reflect.Interface {
			v = reflect.ValueOf(v.Interface())
		}

		switch f := f.(type) {
		case string:
			switch v.Kind() {
			case reflect.Map:
				key := reflect.ValueOf(f)
				val := v.MapIndex(key)
				if val.IsValid() {
					v = val
					continue
				}
			case reflect.Struct:
				fieldName := strings.Title(f)
				fieldValue := v.FieldByName(fieldName)
				if fieldValue.IsValid() {
					v = fieldValue
					continue
				}
			case reflect.Slice:
				if f != "*" || idx >= len(fields) {
					return nil
				}
				for i := 0; i < v.Len(); i++ {
					tmp := ExtractdAny(v.Index(i).Interface(), fields[idx+1:]...)
					if tmp != nil {
						return tmp
					}
				}
			}
			return nil
		case int:
			if v.Kind() == reflect.Slice {
				if v.Len() > f {
					v = v.Index(f)
					continue
				}
			}
			return nil
		default:
			return nil
		}
	}

	return v.Interface()
}
