# gojison

Package __gojison__ provides two simple middlewares
for working with JSON data in goji and two simple
helper function to write out Error and Success responses
to the http.RequestWriter.

[![BuildStatus](https://travis-ci.org/ndyakov/gojison.png)](https://travis-ci.org/ndyakov/gojison)
[![Coverage Status](https://coveralls.io/repos/ndyakov/gojison/badge.svg?branch=master)](https://coveralls.io/r/ndyakov/gojison?branch=master)
[![GoDoc](https://godoc.org/github.com/ndyakov/gojison?status.png)](https://godoc.org/github.com/ndyakov/gojison)
[![status](https://sourcegraph.com/api/repos/github.com/ndyakov/gojison/.badges/status.svg)](https://sourcegraph.com/github.com/ndyakov/gojison)

## Example

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/zenazn/goji"
    "github.com/zenazn/goji/web"
    "github.com/ndyakov/whatever"
    "github.com/ndyakov/gojison"
)

func main() {
    goji.Use(gojison.Request)
    goji.Use(gojison.Response)
    goji.Get("/save", handleSave)
    goji.Serve()
}

func HandleSave(c web.C, w http.ResponseWriter, r *http.Request) {
    params := c.Env["Params"].(whatever.Params)
    if err := params.Required("user.name", "user.email"); err != nil {
        gojison.Error(err, 0)
        return
    }
    // do something...
    gojison.Success("saved", 0)
}
```

## Install

### Get the package

`go get github.com/ndyakov/gojison`

### Import in your source

`import "github.com/ndyakov/gojison"`

## Middlewares

This package includes the following middlewares:

### Response:

This middleware will set the content type of the response to
`"application/json"`

### Request:

This middlware will unmarshal the request parameters to an
[whatever.Params](http://godoc.org/github.com/ndyakov/whatever) structure
and the attach it to the goji context map.
You can access that structure from the context and work with it as either
`map[string]iterface{}` or `whatever.Params`

For more information about how to work with the whatever.Params type, please refer to:
http://godoc.org/github.com/ndyakov/whatever

## Helpers

__gojison__ exports two helper functions:
* Error(w http.ResponseWriter, err error, code int)
* Success(w http.ResponseWriter, message string, code int)

Those helpers can be used to render an error or
success message as json and set the response code.

If the response code argument is 0, the following will be used:
* 400(http.StatusBadRequest) for Error
* 200(http.StatusOK) for Success

## Contributions

Before contributing please execute:
* gofmt
* golint
* govet
