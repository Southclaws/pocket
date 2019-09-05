package pocket

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

var _ = (Responder)(&BasicResponder{})

// Responder describes an error type that can resolve to a HTTP response. This
// means providing a response code and a body.
type Responder interface {
	error
	Headers() http.Header
	Body() io.ReadCloser
	Status() int
}

// BasicResponder provides a basic implementation for the Responder interface
// for writing simple responses.
type BasicResponder struct {
	H http.Header
	R io.ReadCloser
	S int
}

// Error implements the error interface
func (r BasicResponder) Error() string {
	if r.R != nil {
		b, err := ioutil.ReadAll(r.R)
		if err != nil {
			panic(err)
		}
		return string(b)
	}
	return "(none)"
}

// Headers implements Responder
func (r BasicResponder) Headers() http.Header {
	return r.H
}

// Body implements Responder
func (r BasicResponder) Body() io.ReadCloser {
	return r.R
}

// Status implements Responder
func (r BasicResponder) Status() int {
	return r.S
}

// -
// Basic responder helpers
// -

// OK returns a responder with HTTP 200 OK default
func OK() BasicResponder {
	return BasicResponder{S: http.StatusOK}
}

// ErrInternalServerError returns a responder for generic internal server errors
func ErrInternalServerError(e error) BasicResponder {
	return BasicResponder{
		S: http.StatusInternalServerError,
		R: ioutil.NopCloser(bytes.NewBufferString(e.Error())),
	}
}

// ErrUnauthorized returns a responder for generic authorization errors
func ErrUnauthorized(e error) BasicResponder {
	return BasicResponder{
		S: http.StatusUnauthorized,
		R: ioutil.NopCloser(bytes.NewBufferString(e.Error())),
	}
}

// ErrForbidden returns a responder for generic forbidden errors
func ErrForbidden(e error) BasicResponder {
	return BasicResponder{
		S: http.StatusForbidden,
		R: ioutil.NopCloser(bytes.NewBufferString(e.Error())),
	}
}

// ErrNotFound returns a responder for generic not found errors
func ErrNotFound(e error) BasicResponder {
	return BasicResponder{
		S: http.StatusNotFound,
		R: ioutil.NopCloser(bytes.NewBufferString(e.Error())),
	}
}

// ErrBadRequest returns a responder for generic request errors
func ErrBadRequest(e error) BasicResponder {
	return BasicResponder{
		S: http.StatusBadRequest,
		R: ioutil.NopCloser(bytes.NewBufferString(e.Error())),
	}
}
