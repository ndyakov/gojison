package gojison

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Params map[string]interface{}

func (p Params) Add(key string, value interface{}) {
	p[key] = value
}

func (p Params) GetP(key string) Params {
	if val, ok := p[key]; ok {
		return Params(val.(map[string]interface{}))
	}

	return Params{}
}

func (p Params) GetI(key string) interface{} {
	if val, ok := p[key]; ok {
		return val
	}

	return nil
}

func (p Params) Get(key string) string {
	if val, ok := p[key]; ok {
		return stringify(val)
	}

	return ""
}

func (p Params) GetString(key string) string {
	if val, ok := p[key]; ok {
		if vs, ok := val.(string); ok {
			return vs
		}
	}
	return ""
}

func (p Params) GetInt(key string) int {
	result, err := strconv.ParseInt(p.Get(key), 0, 0)
	if err == nil {
		return int(result)
	}

	return 0
}

func (p Params) GetInt8(key string) int8 {
	result, err := strconv.ParseInt(p.Get(key), 0, 8)
	if err == nil {
		return int8(result)
	}

	return 0
}

func (p Params) GetInt64(key string) int64 {
	result, err := strconv.ParseInt(p.Get(key), 0, 64)
	if err == nil {
		return result
	}

	return 0
}

func (p Params) GetFloat32(key string) float32 {
	result, err := strconv.ParseFloat(p.Get(key), 32)
	if err == nil {
		return float32(result)
	}

	return 0
}

func (p Params) GetFloat64(key string) float64 {
	result, err := strconv.ParseFloat(p.Get(key), 64)
	if err == nil {
		return result
	}

	return 0
}

func (p Params) GetFloat(key string) float32 {
	return p.GetFloat32(key)
}

func (p Params) GetTime(key string) time.Time {
	if result, err := time.Parse(time.RFC3339, p.Get(key)); err == nil {
		return result
	}
	return time.Time{}
}

func (p Params) GetSlice(key string) []interface{} {
	if val, ok := p[key]; ok {
		if val, ok := val.([]interface{}); ok {
			return val
		}
	}

	return nil
}

func (p Params) GetSliceStrings(key string) []string {
	result := make([]string, 0)
	if val, ok := p[key]; ok {
		if val, ok := val.([]interface{}); ok {
			for _, v := range val {
				if vs, ok := v.(string); ok {
					result = append(result, vs)
				}
			}
		}
		return result
	}

	return nil
}

func (p Params) GetSliceInts(key string) []int {
	result := make([]int, 0)
	if val, ok := p[key]; ok {
		if slice, ok := val.([]interface{}); ok {
			for _, v := range slice {
				if intVal, err := strconv.ParseInt(stringify(v), 0, 0); err == nil {
					result = append(result, int(intVal))
				}
			}
		}
		return result
	}

	return nil
}

// Validators
func (p Params) Required(keys ...string) error {
	for _, key := range keys {
		if ok := p.exists(p, key); !ok {
			return fmt.Errorf("The parameter %s is required!", key)
		}
	}

	return nil
}

func (p Params) exists(input map[string]interface{}, key string) (ok bool) {
	if input == nil {
		return false
	}

	if index := strings.Index(key, "."); index != -1 {
		pair := strings.SplitN(key, ".", 2)
		params, asserted := input[pair[0]].(map[string]interface{})
		if !asserted {
			return false
		}
		ok = p.exists(params, pair[1])
	} else {
		var v interface{}
		v, ok = input[key]
		if s, isS := v.(string); isS && ok {
			ok = (s != "")
		}
	}

	return ok
}

func stringify(v interface{}) string {
	if vs, ok := v.(string); ok {
		return vs
	}

	return fmt.Sprintf("%v", v)
}
