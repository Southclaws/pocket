/*
Package pocket provides a declarative HTTP handler generator which uses
reflection to facilitate declaring exactly what inputs (query parameters, body
content) a HTTP handler expects. Then, at request time, the function parameters
are hydrated by using the appropriate decoding methods on the incoming HTTP
request. It also allows easier responses by utilising the request's return value
instead of a writer interface.

Note: Documentation is still a little dry while the library is being developed!
This section will eventually contain a full set of examples for usage as well as
details for integrations and best practices.
*/

package pocket
