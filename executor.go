package pocket

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
)

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

// String implements fmt.Stringer
func (h PropsHandler) String() string {
	return fmt.Sprintf("%s -> %s", h.propsT.String(), h.returns.String())
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

// Execute is the function that will hydrate a handler's props for an incoming
// HTTP request. If you're just using the `Handler` helper, you won't need to
// call this function anywhere.
//
// It's similar to middleware except it's not meant to be used in the middle,
// only at the end of a request middleware chain where the actual business logic
// happens. So it's more like "endware" than middleware...
func (h PropsHandler) Execute(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing handler:", h)

	// This could just become a context.Context tbh...
	hctx := Ctx{
		Request: r,
	}

	// create a copy of the props value created during generation time. This is
	// faster than creating a new props object from scratch and is necessary
	// because we don't want to mutate the one stored in the handler as it's
	// shared across many goroutines.
	props := h.propsV
	if err := hydrate(&props, h.propsT, r); err != nil {
		// TODO: handle errors - HTTP bad request? pass to logger?
		log.Panic(err)
	}

	// If the return type of the handler is "Writer" that means it doesn't have
	// an explicit return type on the function signature and the function
	// instead needs access to the underlying http.ResponseWriter for response.
	if h.returns == returnTypeWriter {
		hctx.Writer = &w
	}

	// Do the handler call then react to the response.
	respond(h.function.Call([]reflect.Value{
		reflect.ValueOf(hctx),
		props,
	}), h.returns, w)
}

func respond(results []reflect.Value, rt returnType, w http.ResponseWriter) {
	if results == nil || len(results) == 0 {
		return
	}

	switch rt {
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
			} else if errstring := responder.Error(); errstring != "" {
				if _, err := io.Copy(w, bytes.NewBufferString(errstring)); err != nil {
					panic(err)
				}
			}
		}

	default:
		panic(fmt.Sprintf(
			"unknown handler return type %v something went wrong during handler generation!",
			rt,
		))
	}

	return
}

func hydrate(
	v *reflect.Value,
	t reflect.Type,
	r *http.Request,
) error {
	// reflect.Type's NumField appears to be slightly faster than reflect.Value!
	for i := 0; i < t.NumField(); i++ {
		if err := extractProp(r, t.Field(i), v.Field(i)); err != nil {
			return err
		}
	}
	return nil
}
