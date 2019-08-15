package pocket

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
)

// Handler will take a function and, if it satisfies the rules of Pocket, will
// generate a HTTP handler with the universal signature:
// `func (w http.ResponseWriter, r *http.Request)`. The types of the parameters
// and returns will affect how the handler will preprocess the request before
// calling the handler as well as how it builds a response from the return
// value.
//
// If the supplied function does not match any of the criteria, Handler will
// panic with a description of why.
//
// This function is all you need on the mux side of things, the rest is done via
// responders (returning data and errors.)
func Handler(f interface{}) http.HandlerFunc {
	return GenerateHandler(f).Execute
}

// Ctx represents the underlying lower-level reader/writer interface typically
// associated with HTTP handlers in Go. This type _may_ be added to a handler's
// function signature and, if so, will be hydrated with the corresponding HTTP
// request and response writer at call-time. Note: If the handler specifies a
// return type, the writer will be empty since writing will be handled
// internally by the framework.
type Ctx struct {
	Writer  *http.ResponseWriter
	Request *http.Request
}

// Responder describes an error type that can resolve to a HTTP response. This
// means providing a response code and a body.
type Responder interface {
	error
	Headers() http.Header
	Body() io.ReadCloser
	Status() int
}

// PropsHandler holds information about a handler after it has been analysed for
// props and other features. Once this information has been generated, the
// associated handler can be executed with a HTTP request which will hydrate the
// handler's props with whatever data can be extracted from the HTTP request.
type PropsHandler struct {
	propsV   reflect.Value
	propsT   reflect.Type
	function reflect.Value
	returns  returnType
}

type returnType uint8

const (
	returnTypeWriter    returnType = iota
	returnTypeResponder returnType = iota
	returnTypeError     returnType = iota
)

func (r returnType) String() string {
	switch r {
	case returnTypeWriter:
		return "Writer"
	case returnTypeResponder:
		return "Responder"
	case returnTypeError:
		return "Error"
	}
	return "Unknown"
}

var reflectedResponseType = reflect.TypeOf((*Responder)(nil)).Elem()
var reflectedErrorType = reflect.TypeOf((*error)(nil)).Elem()

// String implements fmt.Stringer
func (h PropsHandler) String() string {
	return fmt.Sprintf("%s -> %s", h.propsT.String(), h.returns.String())
}

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

// Execute is the function that will hydrate a handler's props for an incoming
// HTTP request. If you're just using the `Handler` helper, you won't need to
// call this function anywhere.
func (h PropsHandler) Execute(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing handler:", h)

	for i := 0; i < h.propsV.NumField(); i++ {
		fv := h.propsV.Field(i)
		ft := h.propsT.Field(i)
		log.Printf("	'%v': '%v'\n", fv.Kind(), ft.Name)
		if strings.HasPrefix(ft.Name, "Param") {
			fieldName := ft.Name[5:]
			fieldValue := r.URL.Query().Get(fieldName)
			log.Println("		Field:", fieldName, fieldValue, ft.Type)
			fv.SetString(fieldValue)
		}
	}

	hctx := Ctx{
		Request: r,
	}

	if h.returns == returnTypeWriter {
		hctx.Writer = &w
	}

	results := h.function.Call([]reflect.Value{
		reflect.ValueOf(hctx),
		h.propsV,
	})

	if len(results) != 1 {
		return
	}

	switch h.returns {
	case returnTypeError:
		ev := results[0].Convert(reflectedErrorType)
		if !ev.IsNil() {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(ev.Elem().Interface().(error).Error())); err != nil {
				panic(err)
			}
		}

	case returnTypeResponder:
		ev := results[0].Convert(reflectedResponseType)
		if !ev.IsNil() {
			//nolint:errcheck - this was asserted during handler generation
			responder := ev.Elem().Interface().(Responder)
			w.WriteHeader(responder.Status())
			if body := responder.Body(); body != nil {
				if _, err := io.Copy(w, body); err != nil {
					panic(err)
				}
			}
		}

	case returnTypeWriter:
		ei := results[0].Interface()
		switch ev := ei.(type) {
		case fmt.Stringer:
			if _, err := w.Write([]byte(ev.String())); err != nil {
				panic(err)
			}

		case io.Reader:
			if _, err := io.Copy(w, ev); err != nil {
				panic(err)
			}

		default:
			panic(fmt.Sprintf("don't know how to respond with a %v", ev))
		}

	default:
		panic(fmt.Sprintf(
			"unknown handler return type %v something went wrong during handler generation!",
			h.returns,
		))
	}

	return
}
