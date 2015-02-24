package gojison

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zenazn/goji/web"
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

func TestRequest(t *testing.T) {
	request := []byte(`
		{
			"one": 1,
			"nested":{
				"two": 2
			}
		}
	`)
	body := bytes.NewBuffer(request)
	w := httptest.NewRecorder()
	c := &web.C{Env: make(map[interface{}]interface{})}
	r, err := http.NewRequest("POST", "/", body)

	if err != nil {
		t.Error(err)
	}

	r.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}

	Request(c, TestHandler{}).ServeHTTP(w, r)
	paramsAsInterface, ok := c.Env["Params"]
	if !ok {
		t.Error("Expected params to be set into the context.")
	}

	params, ok := paramsAsInterface.(Params)
	if params == nil || !ok {
		t.Error("Expected params to be unmarshalled into the context.")
	}

	if c.Env["GojisonDecodeError"] != nil {
		t.Error("Expected params to be decoded without an error.")
	}

	if params.Get("one") != "1" {
		wrong(t, "Get on unmarshaled params", "1", params.Get("one"))
	}

	if params.GetP("nested").GetInt("two") != 2 {
		wrong(t, "GetP#GetInt on unmarshaled params", 2, params.GetP("nested").GetInt("two"))
	}
}

func TestRequestWithoutJSON(t *testing.T) {
	request := []byte{}
	body := bytes.NewBuffer(request)
	w := httptest.NewRecorder()
	c := &web.C{Env: make(map[interface{}]interface{})}
	r, err := http.NewRequest("GET", "/", body)

	if err != nil {
		t.Error(err)
	}

	r.Header = map[string][]string{}

	Request(c, TestHandler{}).ServeHTTP(w, r)
	_, ok := c.Env["Params"]
	if ok {
		t.Error("Expected params to be empty when there is no application/json request.")
	}
}

func TestRequestWithInvalidBody(t *testing.T) {
	request := []byte(`
		{"one"}
	`)
	body := bytes.NewBuffer(request)
	w := httptest.NewRecorder()
	c := &web.C{Env: make(map[interface{}]interface{})}
	r, err := http.NewRequest("POST", "/", body)

	if err != nil {
		t.Error(err)
	}

	r.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}

	Request(c, TestHandler{}).ServeHTTP(w, r)
	_, ok := c.Env["Params"]
	if !ok {
		t.Error("Expected params to be set into the context.")
	}

	if c.Env["GojisonDecodeError"] == nil {
		t.Error("Expected to be an Error when decoding the params.")
	}
}

type TestHandler struct{}

func (th TestHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
}
