package gojison

import (
	"net/http"
	"strings"

	"encoding/json"

	"github.com/zenazn/goji/web"
)

func Response(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}
	return http.HandlerFunc(fn)
}

func Request(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		contentTypeSlice := strings.Split(r.Header.Get("Content-Type"), ";")
		if contentTypeSlice[0] == "application/json" {
			var params Params
			c.Env["GojisonDecodeError"] = json.NewDecoder(r.Body).Decode(&params)
			c.Env["Params"] = params
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
