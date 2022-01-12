package value

import (
	"reflect"
	"strconv"
)

// Settable return a settable value indicated by path,
// path can be any of "struct field name/map index/slice index/array index/method name".
func Settable(value reflect.Value, path []string) reflect.Value {
	for _, name := range path {
		if v := trySettableField(value, name); v.IsValid() {
			value = v
		} else if v := tryMethod(value, name); v.IsValid() {
			value = v
		} else if value.Kind() != reflect.Ptr && value.CanAddr() {
			if v := tryMethod(value.Addr(), name); v.IsValid() {
				value = v
			} else {
				return reflect.Value{}
			}
		} else {
			return reflect.Value{}
		}
	}
	return value
}

func trySettableField(value reflect.Value, name string) reflect.Value {
	for {
		switch value.Kind() {
		case reflect.Interface:
			value = value.Elem()
		case reflect.Ptr:
			if value.IsNil() {
				value.Set(reflect.New(value.Type().Elem()))
			}
			value = value.Elem()
		case reflect.Struct:
			return SettableStructField(value, name)
		// map value is always not settable, so don't check it here.
		// maybe we should define and return a "Setter" interface for this case in the future.
		case reflect.Slice, reflect.Array:
			if i, err := strconv.Atoi(name); i >= 0 && err == nil {
				if i < value.Len() {
					return value.Index(i)
				}
				return reflect.Value{}
			}
			return reflect.Value{}
		default:
			return reflect.Value{}
		}
	}
}

func SettableStructField(value reflect.Value, name string) reflect.Value {
	field, ok := value.Type().FieldByName(name)
	if !ok {
		return reflect.Value{}
	}
	for i, index := range field.Index {
		if i > 0 {
			for value.Kind() == reflect.Ptr {
				if value.IsNil() {
					value.Set(reflect.New(value.Type().Elem()))
				}
				value = value.Elem()
			}
		}
		value = value.Field(index)
	}
	return value
}
