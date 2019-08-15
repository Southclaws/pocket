package gorilla

import (
	"fmt"
	"testing"

	"github.com/gorilla/mux"
)

func TestPocket(t *testing.T) {
	r := mux.NewRouter()

	Pocket(
		r.NewRoute().Path("/users/{UserID}"),
		func(props struct {
			UserID string
		}) {
			fmt.Println(props.UserID)
		},
	)
}
