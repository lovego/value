package value

import (
	"reflect"
	"strconv"
)

// Get a value indicated by paths, path can be any of "struct field name/map index/slice index/array index/method name".
func Get(value reflect.Value, paths []string) reflect.Value {
	for _, name := range paths {
		if v := tryField(value, name); v.IsValid() {
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

func tryField(value reflect.Value, name string) reflect.Value {
	for {
		switch value.Kind() {
		case reflect.Ptr, reflect.Interface:
			value = value.Elem()
		case reflect.Struct:
			return value.FieldByName(name)
		case reflect.Map:
			var v = value.MapIndex(reflect.ValueOf(name))
			if v.IsValid() {
				return v
			}
			return reflect.Zero(value.Type().Elem())
		case reflect.Slice, reflect.Array, reflect.String:
			if i, err := strconv.Atoi(name); err == nil {
				return value.Index(i)
			}
			return reflect.Value{}
		default:
			return reflect.Value{}
		}
	}
}

func tryMethod(value reflect.Value, name string) reflect.Value {
	if v := value.MethodByName(name); v.IsValid() && v.Type().NumOut() == 1 {
		return v.Call(nil)[0]
	}
	return reflect.Value{}
}
