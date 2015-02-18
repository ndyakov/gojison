package gojison

import (
	"fmt"
	"strings"
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

func (p Params) Get(key string) string {
	if val, ok := p[key]; ok {
		return fmt.Sprintf("%v", val)
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
	if val, ok := p[key]; ok {
		if vi, ok := val.(int); ok {
			return vi
		}
	}

	return 0
}

func (p Params) GetFloat(key string) float64 {
	if val, ok := p[key]; ok {
		if vf, ok := val.(float64); ok {
			return vf
		}
	}

	return 0
}

func (p Params) GetI(key string) interface{} {
	if val, ok := p[key]; ok {
		return val
	}

	return nil
}

func (p Params) GetSlice(key string) []interface{} {
	if val, ok := p[key]; ok {
		return val.([]interface{})
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
