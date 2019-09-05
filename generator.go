package pocket

import (
	"encoding"
	"fmt"
	"log"
	"reflect"
	"strings"
)

var reflectedResponseType = reflect.TypeOf((*Responder)(nil)).Elem()
var reflectedErrorType = reflect.TypeOf((*error)(nil)).Elem()
var reflectedTextUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

// GenerateHandler generates a props-based HTTP handler from the given function.
// Because this function should only ever be run during program initialisation,
// it will panic if it encounters any errors. To prevent running into errors by
// modifying handler signatures, it's advised you write some simple tests into
// your test suite that ensure handlers will be generated successfully,
// therefore capturing any issues before deployment.
//
// A PropsHandler is a HTTP request handler that makes use of reflected values
// to hydrate a "props" structure argument in the actual handler function.
// This function aims to cache as much as possible so fewer calls to reflection
// functions are necessary during request handling.
func GenerateHandler(f interface{}) (h PropsHandler) {
	h.function = reflect.ValueOf(f)
	if h.function.Kind() != reflect.Func {
		panic(fmt.Sprintf(
			"attempt to generate a props-based handler from a non-function argument %v",
			h.function.Kind(),
		))
	}

	t := reflect.TypeOf(f)

	log.Println("Generating handler for", t)

	for i := 0; i < t.NumIn(); i++ {
		pt := t.In(i)
		pk := pt.Kind()
		log.Println("	Param", pt, "Kind:", pk, pt.PkgPath())

		if strings.HasPrefix(pt.PkgPath(), "github.com/Southclaws/pocket") {
			// internal structure
			continue
		}

		switch pk {
		case reflect.Struct:
			h.propsT = pt

			structFields := []reflect.StructField{}
			for i := 0; i < pt.NumField(); i++ {
				ft := pt.Field(i)
				// TODO: process tags, load plugins
				structFields = append(structFields, ft)
			}

			// generate struct instance now instead of for each request
			h.propsV = reflect.New(reflect.StructOf(structFields)).Elem()

		default:
			panic(fmt.Sprintf("unsupported type in handler parameters: %v", pk))
		}
	}

	numOut := t.NumOut()
	if numOut > 1 {
		panic("attempt to generate a props-based handler with multiple return values")
	}

	if numOut == 1 {
		rt := t.Out(0)
		if rt.Implements(reflectedResponseType) {
			h.returns = returnTypeResponder
		} else if rt.Implements(reflectedErrorType) {
			h.returns = returnTypeError
		} else {
			panic(fmt.Sprintf("unsupported handler return type %v", rt.Name()))
		}
	} else {
		h.returns = returnTypeWriter
	}

	return
}
