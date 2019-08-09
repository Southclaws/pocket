package pocket

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_WithQueryParam(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Handler(func(c HandlerContext, props struct {
		MethodGet
		ParamUserID string
	}) (err error) {
		fmt.Println(
			"\ncontext:", c,
			"\nprops.ParamUserID:", props.ParamUserID,
		)
		return
	}))

	ts := httptest.NewServer(mux)
	defer ts.Close()

	_, err := http.Get(ts.URL + `/?UserID=user1`)
	if err != nil {
		t.Fatal(err)
	}
}
