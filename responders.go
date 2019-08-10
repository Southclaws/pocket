package pocket

import (
	"io"
	"net/http"
)

var _ = (Response)(&BasicResponder{})

// BasicResponder provides a basic implementation for the Response interface for
// writing simple responses.
type BasicResponder struct {
	H http.Header
	R io.ReadCloser
	S int
}

// Headers implements Response
func (r BasicResponder) Headers() http.Header {
	return r.H
}

// Body implements Response
func (r BasicResponder) Body() io.ReadCloser {
	return r.R
}

// Status implements Response
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
