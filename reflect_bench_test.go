package pocket_test

import (
	"os"
	"reflect"
	"testing"
)

type P struct {
	A string
	B string
}

func handler(props P) error {
	return nil
}

var arg = P{A: "one", B: "two"}

var handlerV reflect.Value
var handlerT reflect.Type
var handlerI interface{}
var argV reflect.Value
var argT reflect.Type

func TestMain(m *testing.M) {
	handlerV = reflect.ValueOf(handler)
	handlerT = reflect.TypeOf(handler)
	handlerI = handlerV.Interface()
	argV = reflect.ValueOf(arg)
	argT = handlerT.In(0)

	os.Exit(m.Run())
}

// This is how calls are done by the library at the moment.
// It's quite slow unfortunately...
func BenchmarkReflectCall(b *testing.B) {
	for n := 0; n < b.N; n++ {
		handlerV.Call([]reflect.Value{argV})
	}
}

// This is a normal function call for reference. It's fast, obviously!
func BenchmarkFunctionCall(b *testing.B) {
	for n := 0; n < b.N; n++ {
		handler(arg)
	}
}

// This is an alternative way of doing a call, unfortunately it requires a type
// assertion so this cannot be automated by the library. Any interface that
// would require this, would require the user to explicitly specify this
// *somewhere*... (not sure where yet...)
func BenchmarkFunctionInterfaceCastCall(b *testing.B) {
	for n := 0; n < b.N; n++ {
		handlerI.(func(P) error)(arg)
	}
}
