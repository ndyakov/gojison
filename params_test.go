package gojison

import (
	"bytes"
	"encoding/json"
	"testing"
)

var body = []byte(`
{
	"int": -10,
	"int8": 1,
	"int64": 123456,
	"float32": 3.14,
	"float64": 3.14159265358979,
	"string": "test",
	"arrayString": ["one","two","three"],
	"arrayInt": [1,2,3,4],
	"params": {
		"one": 1,
		"two": 2,
		"params2": {
			"three": 3
		}
	}
}
`)

var bodyStrings = map[string]string{
	"int":         "-10",
	"int8":        "1",
	"int64":       "123456",
	"float64":     "3.14159265358979",
	"string":      "test",
	"arrayString": "[one two three]",
	"arrayInt":    "[1 2 3 4]",
}

func parse(contents []byte) Params {
	var params Params
	request := bytes.NewReader(contents)
	err := json.NewDecoder(request).Decode(&params)
	if err != nil {
		panic(err)
	}
	return params
}

func wrong(t *testing.T, method string, expected, got interface{}) {
	t.Errorf(
		"Params.%s was incorrect.\n Expected: %#v, Got: %#v",
		method,
		expected,
		got,
	)
}

func TestParamsUnmarshal(t *testing.T) {
	var params Params
	request := bytes.NewReader(body)
	err := json.NewDecoder(request).Decode(&params)
	if err != nil {
		t.Error(err)
	}
}

func TestParamsGet(t *testing.T) {
	keys := []string{"int", "int8", "int64", "float64", "string", "arrayString", "arrayInt"}
	params := parse(body)
	for _, key := range keys {
		if params.Get(key) != bodyStrings[key] {
			wrong(t, "Get", bodyStrings[key], params.Get(key))
		}
	}
}

func TestParamGetString(t *testing.T) {
	keys := []string{"float64", "string", "arrayString"}
	expected := map[string]string{"float64": "", "string": "test", "arrayString": ""}
	params := parse(body)
	for _, key := range keys {
		if params.GetString(key) != expected[key] {
			wrong(t, "GetString", expected[key], params.GetString(key))
		}
	}
}

func TestParamGetInt(t *testing.T) {
	params := parse(body)
	keys := []string{"string", "int", "int8", "int64"}
	expected := map[string]int{"string": 0, "int": -10, "int8": 1, "int64": 123456}
	for _, key := range keys {
		if params.GetInt(key) != expected[key] {
			wrong(t, "GetInt", expected[key], params.GetInt(key))
		}
	}
}

func TestParamGetInt8(t *testing.T) {
	params := parse(body)
	keys := []string{"string", "int", "int8", "int64"}
	expected := map[string]int8{"string": 0, "int": -10, "int8": 1, "int64": 0}
	for _, key := range keys {
		if params.GetInt8(key) != expected[key] {
			wrong(t, "GetInt8", expected[key], params.GetInt8(key))
		}
	}
}

func TestParamGetInt64(t *testing.T) {
	params := parse(body)
	keys := []string{"float64", "string", "int", "int8", "int64"}
	expected := map[string]int64{"float64": 0, "string": 0, "int": -10, "int8": 1, "int64": 123456}
	for _, key := range keys {
		if params.GetInt64(key) != expected[key] {
			wrong(t, "GetInt64", expected[key], params.GetInt64(key))
		}
	}
}

func TestParamGetFloat64(t *testing.T) {
	params := parse(body)
	keys := []string{"float64", "string", "int", "int8", "int64"}
	expected := map[string]float64{"float64": 3.14159265358979, "string": 0, "int": -10, "int8": 1, "int64": 123456}
	for _, key := range keys {
		if params.GetFloat64(key) != expected[key] {
			wrong(t, "GetFloat64", expected[key], params.GetFloat64(key))
		}
	}
}

func TestParamGetFloat32(t *testing.T) {
	params := parse(body)
	keys := []string{"float64", "string", "int", "int8", "int64"}
	expected := map[string]float32{"float64": 3.1415927, "string": 0, "int": -10, "int8": 1, "int64": 123456}
	for _, key := range keys {
		if params.GetFloat32(key) != expected[key] {
			wrong(t, "GetFloat32", expected[key], params.GetFloat32(key))
		}
	}
}
