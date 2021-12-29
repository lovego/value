package value_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/lovego/value"
)

type TestStruct struct {
	String  string
	Bool    bool
	Layer   TestStruct2
	Pointer *string
}

type TestStruct2 struct {
	Time time.Time
	Map  map[string][]int
}

func (t TestStruct) Method() string {
	return "方法"
}

func (t *TestStruct) PtrMethod() string {
	return "指针方法"
}

var ts = TestStruct{
	String: "中国",
	Layer: TestStruct2{
		Time: time.Date(2020, 8, 07, 9, 10, 11, 0, time.UTC),
		Map:  map[string][]int{"x": {1, 2, 3}},
	},
}

func ExampleGet_NonPtrType() {
	v := reflect.ValueOf(ts)
	fmt.Println(value.Get(v, []string{"String"}))
	fmt.Println(value.Get(v, []string{"Layer", "Time"}))
	fmt.Println(value.Get(v, []string{"Layer", "Map", "x", "2"}))
	fmt.Println(value.Get(v, []string{"Method"}))
	fmt.Println(value.Get(v, []string{"PtrMethod"}))

	v = reflect.ValueOf(&ts).Elem()
	fmt.Println(value.Get(v, []string{"PtrMethod"}))

	// Output:
	// 中国
	// 2020-08-07 09:10:11 +0000 UTC
	// 3
	// 方法
	// <invalid reflect.Value>
	// 指针方法
}

func ExampleGet_PtrType() {
	v := reflect.ValueOf(&ts)
	fmt.Println(value.Get(v, []string{"String"}))
	fmt.Println(value.Get(v, []string{"Layer", "Time"}))
	fmt.Println(value.Get(v, []string{"Method"}))
	fmt.Println(value.Get(v, []string{"PtrMethod"}))

	// none exists
	fmt.Println(value.Get(v, []string{"NoneExists"}))
	fmt.Println(value.Get(v, []string{"Layer", "Map", "x", "two"}))
	fmt.Println(value.Get(v, []string{"Layer", "Map", "x", "2", "two"}))

	// Output:
	// 中国
	// 2020-08-07 09:10:11 +0000 UTC
	// 方法
	// 指针方法
	// <invalid reflect.Value>
	// <invalid reflect.Value>
	// <invalid reflect.Value>
}

// https://golang.org/ref/spec#Method_sets
// The method set of the pointer type *T is the set of all methods declared with receiver *T or T
// (that is, it also contains the method set of T).
func ExampleNonPtrType_MethodByName() {
	t := reflect.TypeOf(TestStruct{})
	method, _ := t.MethodByName("Method")
	fmt.Println(method.Type)
	method, _ = t.MethodByName("PtrMethod")
	fmt.Println(method.Type)

	// Output:
	// func(value_test.TestStruct) string
	// <nil>
}

func ExamplePtrType_MethodByName() {
	t := reflect.TypeOf(&TestStruct{})
	method, _ := t.MethodByName("Method")
	fmt.Println(method.Type)
	method, _ = t.MethodByName("PtrMethod")
	fmt.Println(method.Type)

	// Output:
	// func(*value_test.TestStruct) string
	// func(*value_test.TestStruct) string
}

func BenchmarkGet(b *testing.B) {
	v := reflect.ValueOf(ts)
	for i := 0; i < b.N; i++ {
		value.Get(v, []string{"Layer", "Map", "x", "2"}) // 400~500 ns
	}
}
