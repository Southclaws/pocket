package gorilla

import (
	"fmt"
	"reflect"

	"github.com/Southclaws/pocket"
	"github.com/gorilla/mux"
)

func Pocket(r *mux.Route, f interface{}) *mux.Route {
	handler := pocket.GenerateHandler(f)
	rv := reflect.ValueOf(*r)

	if path := rv.FieldByName("regexp").FieldByName("path"); !path.IsNil() {
		path := path.Elem()
		for i := 0; i < path.NumField(); i++ {
			fv := path.Field(i)
			fmt.Println(fv)

			// TODO: build some kind of plugin API for registering the things
			// that match to prop names as well as functions for extracting the
			// values.
			// also come up with a name for said things that match to props...
			// handler.TryRegister(fv, func())
		}
	}
	return r.HandlerFunc(handler.Execute)
}
