package pocket

import (
	"io"
	"net/http"
)

var _ = (Responder)(&BasicResponder{})

// BasicResponder provides a basic implementation for the Responder interface
// for writing simple responses.
type BasicResponder struct {
	H http.Header
	R io.ReadCloser
	S int
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
