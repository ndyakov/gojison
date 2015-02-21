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
	"nestedParams": {
		"one": 1,
		"two": 2,
		"params2": {
			"three": 3
		}
	}
}
`)

func TestParamsUnmarshal(t *testing.T) {
	var params Params
	request := bytes.NewReader(body)
	err := json.NewDecoder(request).Decode(&params)
	if err != nil {
		t.Error(err)
	}
}

func TestParamsAdd(t *testing.T) {
	p := Params{}
	if _, ok := p["one"]; ok {
		wrong(t, "[\"one\"]", false, true)
	}

	overwrite := p.Add("one", 2)
	if overwrite {
		wrong(t, "Add", false, true)
	}

	overwrite = p.Add("one", 1)
	if !overwrite {
		wrong(t, "Add", true, false)
	}

	if val, ok := p["one"]; !ok {
		wrong(t, "[\"one\"]", true, false)
		if val != 1 {
			wrong(t, "[\"one\"]", 1, val)
		}
	}
}

func TestParamsAddRecursive(t *testing.T) {
	p := Params{"one": 1}
	// This should detect the recursive adding and
	// should NOT add p with the params key.
	p.Add("params", p)
	if _, ok := p["params"]; ok {
		wrong(t, "[\"params\"]", false, true)
	}
}

func TestParamsRequired(t *testing.T) {
	p := Params{"one": 1, "emtpyString": ""}
	if err := p.Required("one"); err != nil {
		wrong(t, "Required", nil, err)
	}

	if err := p.Required("emptyString"); err == nil {
		wrong(t, "Required", "the parameter emptyString is required", nil)
	}

	if err := p.Required(); err != nil {
		wrong(t, "Required", nil, err)
	}

	if err := p.Required("two"); err == nil {
		wrong(t, "Required", "the parameter two is required", nil)
	}

	p1 := Params{"two": 2}
	p.Add("nestedParams", p1)

	if err := p.Required("nestedParams.one"); err == nil {
		wrong(t, "Required", "the parameter nestedParams.one is required", nil)
	}

	if err := p.Required("nestedParams.two"); err != nil {
		wrong(t, "Required", nil, err)
	}

	m := map[string]interface{}{"three": 3, "emptyString": ""}
	p.Add("nestedMap", m)

	if err := p.Required("nestedMap.one"); err == nil {
		wrong(t, "Required", "the parameter nestedMap.one is required", nil)
	}

	if err := p.Required("nestedMap.three"); err != nil {
		wrong(t, "Required", nil, err)
	}

	if err := p.Required("nestedMap.emptyString"); err == nil {
		wrong(t, "Required", "the parameter nestedMap.emptyString is required", err)
	}

	if err := p.Required("missing.missing"); err == nil {
		wrong(t, "Required", "the parameter missing.missing is required", nil)
	}

	if p.exists(nil, "nil") {
		wrong(t, "exists", false, true)
	}
}

func TestParamsDelete(t *testing.T) {
	p := Params{"one": 1}
	p.Remove("one")

	if _, ok := p["one"]; ok {
		wrong(t, "[\"one\"]", false, true)
	}
}

func TestParamsEmpty(t *testing.T) {
	p := Params{}
	if !p.Empty() {
		wrong(t, "Empty", true, false)
	}

	p.Add("one", 1)

	if p.Empty() {
		wrong(t, "Empty", false, true)
	}
}

func TestParamsGet(t *testing.T) {
	keys := []string{
		"int",
		"int8",
		"int64",
		"float64",
		"string",
		"arrayString",
		"arrayInt",
	}

	expected := map[string]string{
		"int":          "-10",
		"int8":         "1",
		"int64":        "123456",
		"float64":      "3.14159265358979",
		"string":       "test",
		"arrayStrings": "[one two three]",
		"arrayInts":    "[1 2 3 4]",
	}

	params := parse(body)
	for _, key := range keys {
		got := params.Get(key)
		if got != expected[key] {
			wrong(t, "Get", expected[key], got)
		}
	}
}

func TestParamsGetP(t *testing.T) {
	keys := []string{"one", "two"}
	expected := map[string]string{
		"one": "1",
		"two": "2",
	}
	params := parse(body)
	nestedParams := params.GetP("nestedParams")
	if nestedParams.Empty() {
		wrong(t, "GetP", "existing nested params", "empty params set")
	}

	for _, key := range keys {
		got := nestedParams.Get(key)
		if got != expected[key] {
			wrong(t, "Get", expected[key], got)
		}
	}

	nestedParams.Add("whoah", Params{"one": 1})
	whoah := nestedParams.GetP("whoah")
	if whoah.Empty() {
		wrong(t, "GetP", "existing nested params", "empty params set")
	}

	missing := whoah.GetP("missing")
	if !missing.Empty() {
		wrong(t, "GetP", "to return empty Params when missing", "something that was not empty")
	}
}

// Well... this will return the same as
// p[key].
func TestParamsGetI(t *testing.T) {
	keys := []string{
		"int",
		"string",
		"time",
		"incorectTime",
		"arrayStrings",
		"missing",
	}

	params := parse(body)
	for _, key := range keys {
		got := params.GetI(key)
		v := params[key]
		if stringify(got) != stringify(v) {
			wrong(t, "GetI", v, got)
		}
	}

	now := time.Now()
	params.Add("now", now)
	t1 := params.GetI("now").(time.Time)
	if now != t1 {
		wrong(t, "GetI", now, t1)
	}

	// or even pointers
	params.Add("testing", t)
	ti := params.GetI("testing").(*testing.T)
	if t != ti {
		wrong(t, "GetI", t, ti)
	}
}

func TestParamsURLValues(t *testing.T) {
	keys := []string{
		"int",
		"int8",
		"int64",
		"float64",
		"string",
		"arrayString",
		"arrayInt",
		"nestedParams.one",
		"nestedParams.params2.three",
		"inParams.one",
	}

	expected := map[string]string{
		"int":                        "-10",
		"int8":                       "1",
		"int64":                      "123456",
		"float64":                    "3.14159265358979",
		"string":                     "test",
		"arrayStrings":               "one",
		"arrayInts":                  "1",
		"nestedParams.one":           "1",
		"nestedParams.params2.three": "3",
		"inParams.one":               "1",
	}
	params := parse(body)
	params.Add("inParams", Params{"one": 1})
	urlValues := params.URLValues("", "")
	for _, key := range keys {
		got := urlValues.Get(key)
		if got != expected[key] {
			wrong(t, fmt.Sprintf("URLValues key: %s", key), expected[key], got)
		}
	}
}

func TestParamsURLValuesWithPrefixAndSufix(t *testing.T) {
	keys := []string{
		"int",
		"int8",
		"int64",
		"float64",
		"string",
		"arrayString",
		"arrayInt",
		"nestedParams[one]",
		"nestedParams[params2][three]",
		"inParams[one]",
	}

	expected := map[string]string{
		"int":                          "-10",
		"int8":                         "1",
		"int64":                        "123456",
		"float64":                      "3.14159265358979",
		"string":                       "test",
		"arrayStrings":                 "one",
		"arrayInts":                    "1",
		"nestedParams[one]":            "1",
		"nestedParams[params2][three]": "3",
		"inParams[one]":                "1",
	}
	params := parse(body)
	params.Add("inParams", Params{"one": 1})
	urlValues := params.URLValues("[", "]")
	for _, key := range keys {
		got := urlValues.Get(key)
		if got != expected[key] {
			wrong(t, fmt.Sprintf("URLValues key: %s", key), expected[key], got)
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
