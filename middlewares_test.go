package gojison

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
}

func TestResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := &http.Request{}
	Response(TestHandler{}).ServeHTTP(w, r)
	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected content-type to be set to application/json")
	}
}

type TestHandler struct{}

func (th TestHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
}
