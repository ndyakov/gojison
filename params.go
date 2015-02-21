package gojison

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Params map[string]interface{}

func (p Params) Add(key string, value interface{}) bool {
	if fmt.Sprintf("%p", p) == fmt.Sprintf("%p", value) {
		return false
	}

	_, ok := p[key]
	p[key] = value
	return ok
}

func (p Params) Empty() bool {
	for _ = range p {
		return false
	}
	return true
}

func (p Params) Remove(key string) {
	delete(p, key)
}

func (p Params) GetP(key string) Params {
	if val, ok := p[key]; ok {
		if vmap, ok := val.(map[string]interface{}); ok {
			return Params(vmap)
		} else if vp, ok := val.(Params); ok {
			return vp
		}
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
	var result []string
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
	var result []int
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

// Get map[string][]string
func (p Params) GetURLValues() url.Values {
	return p.getURLValues(p)
}

func (p Params) getURLValues(set Params) url.Values {
	result := url.Values{}
	var subset url.Values
	for key, value := range set {
		var foundSubset bool
		if v, ok := value.(Params); ok {
			subset = p.getURLValues(v)
			foundSubset = true
		} else if v, ok := value.(map[string]interface{}); ok {
			subset = p.getURLValues(Params(v))
			foundSubset = true
		} else if v, ok := value.([]interface{}); ok {
			for _, el := range v {
				result[key] = append(result[key], stringify(el))
			}
			continue
		}
		if foundSubset {
			for k, v := range subset {
				nestedKey := fmt.Sprintf("%s.%s", key, k)
				result[nestedKey] = append(result[nestedKey], v...)
			}
		} else {
			result[key] = append(result[key], stringify(value))
		}
	}
	return result
}

// Validators
func (p Params) Required(keys ...string) error {
	for _, key := range keys {
		if ok := p.exists(p, key); !ok {
			return fmt.Errorf("the parameter %s is required", key)
		}
	}

	return nil
}

func (p Params) exists(input map[string]interface{}, key string) (ok bool) {
	if input == nil {
		return false
	}

	if index := strings.Index(key, "."); index != -1 {
		var casted bool
		var params map[string]interface{}
		pair := strings.SplitN(key, ".", 2)

		if params, casted = input[pair[0]].(Params); !casted {
			params, casted = input[pair[0]].(map[string]interface{})
		}

		if !casted {
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
