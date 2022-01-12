package value_test

import (
	"fmt"
	"reflect"
	"time"

	"github.com/lovego/value"
)

func ExampleSettable() {
	var ts = TestStruct{
		Interface: &TestStruct2{},
	}
	var v = reflect.ValueOf(&ts)
	value.Settable(v, []string{"String"}).Set(reflect.ValueOf("ok"))
	fmt.Println(ts.String)

	value.Settable(v, []string{"Bool"}).Set(reflect.ValueOf(true))
	fmt.Println(ts.Bool)

	var time = reflect.ValueOf(time.Date(2022, 1, 12, 19, 30, 30, 0, time.UTC))
	value.Settable(v, []string{"Time"}).Set(time)
	fmt.Println(ts.Time)

	value.Settable(v, []string{"Interface", "Time"}).Set(time)
	fmt.Println(ts.Interface.(*TestStruct2).Time)

	value.Settable(v, []string{"Interface"}).Set(reflect.ValueOf(9))
	fmt.Println(ts.Interface)

	var s = "string"
	value.Settable(v, []string{"Pointer"}).Set(reflect.ValueOf(&s))
	fmt.Println(*ts.Pointer)

	fmt.Println(value.Settable(v, []string{"NonExists"}))
	// Output:
	// ok
	// true
	// 2022-01-12 19:30:30 +0000 UTC
	// 2022-01-12 19:30:30 +0000 UTC
	// 9
	// string
	// <invalid reflect.Value>
}

func ExampleSettable_layer() {
	var ts = TestStruct{
		Layer: TestStruct2{Slice: []int{0, 1, 2}},
	}
	var v = reflect.ValueOf(&ts).Elem()

	value.Settable(v, []string{"Layer", "Slice", "2"}).Set(reflect.ValueOf(22))
	fmt.Println(ts.Layer.Slice[2])
	fmt.Println(value.Settable(v, []string{"Layer", "Slice", "3"}))
	fmt.Println(value.Settable(v, []string{"Layer", "Slice", "none"}))

	value.Settable(v, []string{"Layer", "Map"}).Set(reflect.ValueOf(map[string][]int{"k": {7}}))
	fmt.Println(ts.Layer.Map)

	value.Settable(v, []string{"PointerLayer", "Map2"}).Set(reflect.ValueOf(map[string]int{"k": 8}))
	fmt.Println(ts.PointerLayer.Map2)

	fmt.Println(value.Settable(v, []string{"Layer", "Map", "key"}))
	fmt.Println(value.Settable(v, []string{"NonExists"}))
	// Output:
	// 22
	// <invalid reflect.Value>
	// <invalid reflect.Value>
	// map[k:[7]]
	// map[k:8]
	// <invalid reflect.Value>
	// <invalid reflect.Value>
}

func ExampleSettable_method() {
	var ts = TestStruct{}
	var v = reflect.ValueOf(&ts)

	value.Settable(v, []string{"SettableMethod"}).Elem().Set(reflect.ValueOf("ok"))
	fmt.Println(ts.String)

	value.Settable(v.Elem(), []string{"SettableMethod"}).Elem().Set(reflect.ValueOf("ok2"))
	fmt.Println(ts.String)

	// Output:
	// ok
	// ok2
}

func ExampleValue_FieldByName() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	fmt.Println(reflect.ValueOf(TestStruct{}).FieldByName("Time"))

	// Output:
	// reflect: indirection through nil pointer to embedded struct
}
