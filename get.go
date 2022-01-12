package value

import (
	"reflect"
	"strconv"
)

// Get return a value indicated by path,
// path can be any of "struct field name/map index/slice index/array index/method name".
func Get(value reflect.Value, path []string) reflect.Value {
	for _, name := range path {
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
		case reflect.Interface:
			value = value.Elem()
		case reflect.Ptr:
			if value.IsNil() {
				value = reflect.Zero(value.Type().Elem())
			} else {
				value = value.Elem()
			}
		case reflect.Struct:
			return GetStructField(value, name)
		case reflect.Map:
			var v = value.MapIndex(reflect.ValueOf(name))
			if v.IsValid() {
				return v
			}
			return reflect.Zero(value.Type().Elem())
		case reflect.Slice, reflect.Array:
			if i, err := strconv.Atoi(name); i >= 0 && err == nil {
				if i < value.Len() {
					return value.Index(i)
				}
				return reflect.Zero(value.Type().Elem())
			}
			return reflect.Value{}
		default:
			return reflect.Value{}
		}
	}
}

func GetStructField(value reflect.Value, name string) reflect.Value {
	field, ok := value.Type().FieldByName(name)
	if !ok {
		return reflect.Value{}
	}
	for i, index := range field.Index {
		if i > 0 {
			for value.Kind() == reflect.Ptr {
				if value.IsNil() {
					value = reflect.Zero(value.Type().Elem())
				} else {
					value = value.Elem()
				}
			}
		}
		value = value.Field(index)
	}
	return value
}

func tryMethod(value reflect.Value, name string) reflect.Value {
	if v := value.MethodByName(name); v.IsValid() && v.Type().NumOut() == 1 {
		return v.Call(nil)[0]
	}
	return reflect.Value{}
}
