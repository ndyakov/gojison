// Package gojison provides two simple middlewares
// for working with JSON data in goji and two simple
// helper function to write out Error and Success responses
// to the http.RequestWriter.
// This package includes the following middlewares:
//
// Response:
//
// This middleware will set the content type of the response to:
//     "application/json"
//
// Request:
//
// This middlware will unmarshal the request parameters to an
// whatever.Params structure and the attach it to the goji context map.
// You can access that structure from the context and work with it as either:
//     map[string]iterface{} or whatever.Params
//
// For more information about how to work with the whatever.Params type, please refer to:
// http://godoc.org/github.com/ndyakov/whatever
//
//
// Example:
//
//		package main
//
//		import (
//				"fmt"
//				"net/http"
//
//				"github.com/zenazn/goji"
//				"github.com/zenazn/goji/web"
//				"github.com/ndyakov/whatever"
//				"github.com/ndyakov/gojison"
//		)
//
//		func main() {
//				goji.Use(gojison.Request)
//				goji.Use(gojison.Response)
//				goji.Get("/save", handleSave)
//				goji.Serve()
//		}
//
//		func handleSave(c web.C, w http.ResponseWriter, r *http.Request) {
//				params := c.Env["Params"].(whatever.Params)
//				if err := params.Required("user.name", "user.email"); err != nil {
//						gojison.Error(err, 0)
//						return
//				}
//				// do something...
//				gojison.Success("saved", 0)
//		}
package gojison
