package pocket

import "net/http"

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
