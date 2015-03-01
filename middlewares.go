package gojison

import (
	"net/http"
	"strings"

	"encoding/json"

	"github.com/ndyakov/whatever"
	"github.com/zenazn/goji/web"
)

// Response will set the Content-Type of the http response to:
//     "application/json"
func Response(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}
	return http.HandlerFunc(fn)
}

// Request will parse the request body to a whatever.Params structure and
// then add this structure to the goji context map with the key "Params".
// The error (or nil) of the decoding will be available in the context with
// the key "GojisonDecodeError".
//
// The parsing of the body will happen only if the Content-Type of the request
// is application/json.
//
// For more information about how to work with the whatever.Params type, please refer to:
// http://godoc.org/github.com/ndyakov/whatever
func Request(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		contentTypeSlice := strings.Split(r.Header.Get("Content-Type"), ";")
		if contentTypeSlice[0] == "application/json" {
			var params whatever.Params
			c.Env["GojisonDecodeError"] = json.NewDecoder(r.Body).Decode(&params)
			c.Env["Params"] = params
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
