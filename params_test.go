package gojison

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

var body = []byte(`
{
	"int": -10,
	"int8": 1,
	"int64": 123456,
	"float32": 3.14,
	"float64": 3.14159265358979,
	"string": "test",
	"time": "2015-02-20T21:22:23.24Z",
	"incorectTime": "2015-02-20",
	"arrayStrings": ["one","two","three"],
	"arrayInts": [1,2,3,4],
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
	"int":          "-10",
	"int8":         "1",
	"int64":        "123456",
	"float64":      "3.14159265358979",
	"string":       "test",
	"arrayStrings": "[one two three]",
	"arrayInts":    "[1 2 3 4]",
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
		got := params.Get(key)
		if got != bodyStrings[key] {
			wrong(t, "Get", bodyStrings[key], got)
		}
	}
}

func TestParamGetString(t *testing.T) {
	keys := []string{"float64", "string", "arrayString"}
	expected := map[string]string{
		"float64":     "",
		"string":      "test",
		"arrayString": "",
	}
	params := parse(body)
	for _, key := range keys {
		got := params.GetString(key)
		if got != expected[key] {
			wrong(t, "GetString", expected[key], got)
		}
	}
}

func TestParamGetInt(t *testing.T) {
	params := parse(body)
	keys := []string{"string", "int", "int8", "int64"}
	expected := map[string]int{"string": 0, "int": -10, "int8": 1, "int64": 123456}
	for _, key := range keys {
		got := params.GetInt(key)
		if got != expected[key] {
			wrong(t, "GetInt", expected[key], got)
		}
	}
}

func TestParamGetInt8(t *testing.T) {
	params := parse(body)
	keys := []string{"string", "int", "int8", "int64"}
	expected := map[string]int8{
		"string": 0,
		"int":    -10,
		"int8":   1,
		"int64":  0,
	}
	for _, key := range keys {
		got := params.GetInt8(key)
		if got != expected[key] {
			wrong(t, "GetInt8", expected[key], got)
		}
	}
}

func TestParamGetInt64(t *testing.T) {
	params := parse(body)
	keys := []string{"float64", "string", "int", "int8", "int64"}
	expected := map[string]int64{
		"float64": 0,
		"string":  0,
		"int":     -10,
		"int8":    1,
		"int64":   123456,
	}
	for _, key := range keys {
		got := params.GetInt64(key)
		if got != expected[key] {
			wrong(t, "GetInt64", expected[key], got)
		}
	}
}

func TestParamGetFloat(t *testing.T) {
	params := parse(body)
	keys := []string{"float64", "string", "int", "int8", "int64"}
	expected := map[string]float32{
		"float64": 3.1415927,
		"string":  0,
		"int":     -10,
		"int8":    1,
		"int64":   123456,
	}
	for _, key := range keys {
		got := params.GetFloat(key)
		if got != expected[key] {
			wrong(t, "GetFloat", expected[key], got)
		}
	}
}
func TestParamGetFloat64(t *testing.T) {
	params := parse(body)
	keys := []string{"float64", "string", "int", "int8", "int64"}
	expected := map[string]float64{
		"float64": 3.14159265358979,
		"string":  0,
		"int":     -10,
		"int8":    1,
		"int64":   123456,
	}
	for _, key := range keys {
		got := params.GetFloat64(key)
		if got != expected[key] {
			wrong(t, "GetFloat64", expected[key], got)
		}
	}
}

func TestParamGetFloat32(t *testing.T) {
	params := parse(body)
	keys := []string{
		"float64",
		"string",
		"int",
		"int8",
		"int64",
	}

	expected := map[string]float32{
		"float64": 3.1415927,
		"string":  0,
		"int":     -10,
		"int8":    1,
		"int64":   123456,
	}

	for _, key := range keys {
		got := params.GetFloat32(key)
		if got != expected[key] {
			wrong(t, "GetFloat32", expected[key], got)
		}
	}
}

func TestParamGetTime(t *testing.T) {
	params := parse(body)
	keys := []string{"float64", "string", "time", "incorectTime"}
	expected := map[string]time.Time{
		"float64": time.Time{},
		"string":  time.Time{},
		"time": time.Date(
			2015,          //Year
			time.February, //Month
			20,            //Day
			21,            //Hours
			22,            //Minutes
			23,            //Sec
			240000000,     //Nanosec
			time.UTC),     //Location (UTC)
		"incorectTime": time.Time{},
	}
	for _, key := range keys {
		got := params.GetTime(key)
		if got != expected[key] {
			wrong(t, "GetTime", expected[key], got)
		}
	}
}

func TestParamsGetSliceStrings(t *testing.T) {
	params := parse(body)
	keys := []string{"missing", "string", "arrayStrings"}
	expected := map[string][]string{
		"missing":      []string{},
		"string":       []string{},
		"arrayStrings": []string{"one", "two", "three"},
	}
	for _, key := range keys {
		got := params.GetSliceStrings(key)
		if !equalSlicesStrings(expected[key], got) {
			wrong(t, "GetSliceStrings", expected[key], got)
		}
	}
}

func TestParamsGetSliceInts(t *testing.T) {
	params := parse(body)
	keys := []string{"missing", "int", "arrayInts"}
	expected := map[string][]int{
		"missing":   []int{},
		"int":       []int{},
		"arrayInts": []int{1, 2, 3, 4},
	}
	for _, key := range keys {
		got := params.GetSliceInts(key)
		if !equalSlicesInts(expected[key], got) {
			wrong(t, "GetSliceInts", expected[key], got)
		}
	}
}

func TestParamsGetSlice(t *testing.T) {
	params := parse(body)
	keys := []string{"missing", "int", "arrayInts", "arrayStrings"}
	expected := map[string][]interface{}{
		"missing":      []interface{}{},
		"int":          []interface{}{},
		"arrayInts":    []interface{}{1, 2, 3, 4},
		"arrayStrings": []interface{}{"one", "two", "three"},
	}
	for _, key := range keys {
		got := params.GetSlice(key)
		if !equalSlicesInterfaces(expected[key], got) {
			wrong(t, "GetSlice", expected[key], got)
		}
	}
}

// Comparing slices
func equalSlicesStrings(expected, got []string) bool {
	if len(expected) != len(got) {
		return false
	}

	for ie, e := range expected {
		for ig, g := range got {
			if g == e {
				got = append(got[:ig], got[ig+1:]...)
			}
		}

		if len(got) != len(expected)-(1+ie) {
			return false
		}
	}

	if len(got) != 0 {
		return false
	}

	return true
}

func equalSlicesInts(expected, got []int) bool {
	if len(expected) != len(got) {
		return false
	}

	for ie, e := range expected {
		for ig, g := range got {
			if g == e {
				got = append(got[:ig], got[ig+1:]...)
			}
		}

		if len(got) != len(expected)-(1+ie) {
			return false
		}
	}

	if len(got) != 0 {
		return false
	}

	return true
}

func equalSlicesInterfaces(expected, got []interface{}) bool {
	if len(expected) != len(got) {
		return false
	}

	for ie, e := range expected {
		for ig, g := range got {
			if fmt.Sprintf("%v", g) == fmt.Sprintf("%v", e) {
				got = append(got[:ig], got[ig+1:]...)
			}
		}

		if len(got) != len(expected)-(1+ie) {
			return false
		}
	}

	if len(got) != 0 {
		return false
	}

	return true
}
