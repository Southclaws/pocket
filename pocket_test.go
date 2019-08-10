package pocket_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Southclaws/pocket"
	"gotest.tools/assert"
)

func withHandler(h http.HandlerFunc, pathQuery string) *http.Response {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + pathQuery)
	if err != nil {
		panic(err)
	}

	return resp
}

func TestHandler_WithQueryParam(t *testing.T) {
	withHandler(pocket.Handler(func(c pocket.Ctx, props struct {
		pocket.MethodGet
		ParamUserID string
	}) {

		assert.Equal(t, "user1", props.ParamUserID)
		assert.Assert(t, c.Writer != nil,
			"the writer should not be nil as there is no return value")

		return
	}), `/?UserID=user1`)
}

func TestHandler_WithReturn(t *testing.T) {
	withHandler(pocket.Handler(func(c pocket.Ctx, props struct {
		pocket.MethodGet
	}) error {

		assert.Assert(t, c.Writer == nil,
			"the writer should be nil when there is a return value present")

		return nil
	}), `/?UserID=user1`)
}
