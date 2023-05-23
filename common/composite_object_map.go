package common

import (
	"encoding/json"
	"reflect"
	"strings"
)

type CompositeObjectMap struct {
	m UnknownMap
}

func (com *CompositeObjectMap) Delete(key string) bool {
	return com.m.Delete(key)
}
func (com *CompositeObjectMap) GetStr(key string) (string, bool) {
	return com.m.GetStr(key)
}
func (com *CompositeObjectMap) GetKeys() []string {
	return com.m.GetKeys()
}
func (com *CompositeObjectMap) Set(key string, value interface{}) {
	com.m.Set(key, value)
}
func (com *CompositeObjectMap) Get(key string, store *interface{}) error {
	return com.m.Get(key, store)
}
func (com *CompositeObjectMap) GetRaw(key string) (interface{}, bool) {
	return com.m.GetRaw(key)
}
func (com *CompositeObjectMap) FullMarshal(c interface{}) ([]byte, error) {
	objTags := com.getObjJsonTags(c, func(key string, value interface{}) {
		com.m.Set(key, value)
	})
	b, err := json.Marshal(com.m)
	for _, i := range objTags {
		com.m.Delete(i)
	}
	return b, err
}

func (com *CompositeObjectMap) Unmarshal(b []byte, c interface{}) error {
	if err := com.m.UnmarshalJSON(b); err != nil {
		return err
	}
	objTags := com.getObjJsonTags(c, func(key string, value interface{}) {})
	for _, i := range objTags {
		com.m.Delete(i)
	}
	return nil
}

func (com *CompositeObjectMap) getObjJsonTags(c interface{}, fn func(key string, value interface{})) []string {
	a := []string{}
	v := reflect.ValueOf(c)
	t := v.Type()
	for t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			if f.Name == "CompositeObjectMap" {
				continue
			}
			a = append(a, com.getObjJsonTags(v.Field(i).Interface(), fn)...)
			continue
		}
		if f.Name[0:1] == strings.ToLower(f.Name[0:1]) {
			continue
		}
		j := f.Tag.Get("json")
		if j == "" {
			j = f.Name
		}
		a = append(a, j)
		fn(j, v.Field(i).Interface())
	}
	return a
}
