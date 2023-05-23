package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

/*
*
This will help in marshalling and unmarshalling maps of interface properly. This makes s
*/
type ErrKeyNotFound struct {
	Key string
}

func (e *ErrKeyNotFound) Error() string {
	return e.Key + " not found"
}

type UnknownMap struct {
	data    map[string]interface{}
	rawData map[string]json.RawMessage
}

func (um *UnknownMap) Get(key string, store interface{}) error {
	rv := reflect.ValueOf(store)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(store)}
	}
	if v, ok := um.data[key]; !ok {
		if i, ok := um.rawData[key]; !ok {
			return &ErrKeyNotFound{Key: key}
		} else {
			if err := json.Unmarshal(i, store); err == nil {
				um.data[key] = reflect.Indirect(rv).Interface()
				delete(um.rawData, key)
				return nil
			} else {
				return err
			}
		}

	} else {
		//TODO: need to fix this with proper assignment
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, store)
		//rv.Set(reflect.ValueOf(v).Elem())
		return err
	}
}

func (um *UnknownMap) GetRef(key string) (interface{}, bool) {
	if v, ok := um.data[key]; ok {
		return v, true
	}
	return nil, false
}

func (um *UnknownMap) Set(key string, value interface{}) {
	if v, ok := value.(json.RawMessage); ok {
		if um.rawData == nil {
			um.rawData = map[string]json.RawMessage{}
		}
		delete(um.data, key)
		um.rawData[key] = v
		return
	}
	if um.data == nil {
		um.data = map[string]interface{}{}
	}
	um.data[key] = value
}

func (um UnknownMap) MarshalJSON() ([]byte, error) {
	//preparing for marshal. Not all the data would be there in data. we might have to combine them
	d := map[string]interface{}{}
	for k, v := range um.rawData {
		d[k] = v
	}
	for k, v := range um.data {
		d[k] = v
	}
	return json.Marshal(d)
}

func (um *UnknownMap) UnmarshalJSON(data []byte) error {
	//s := string(data)
	//fmt.Println(s)
	um.rawData = map[string]json.RawMessage{}
	um.data = map[string]interface{}{}
	err := json.Unmarshal(data, &(um.rawData))
	return err
}

func (um *UnknownMap) GetKeys() []string {
	keys := []string{}
	for k := range um.data {
		keys = append(keys, k)
	}
	for k := range um.rawData {
		keys = append(keys, k)
	}
	return keys
}

func (um *UnknownMap) Delete(key string) bool {
	if _, ok := um.data[key]; ok {
		delete(um.data, key)
		return true
	}
	if _, ok := um.rawData[key]; ok {
		delete(um.rawData, key)
		return true
	}
	return false
}

func (um *UnknownMap) Copy(shallowCopy bool) *UnknownMap {
	data := um.data
	rawData := um.rawData
	if shallowCopy {
		data = map[string]interface{}{}
		rawData = map[string]json.RawMessage{}
		for k, v := range um.data {
			data[k] = v
		}
		for k, v := range um.rawData {
			rawData[k] = v
		}
	}
	return &UnknownMap{
		data:    data,
		rawData: rawData,
	}
}

func (um *UnknownMap) GetRaw(key string) (interface{}, bool) {
	if v, ok := um.data[key]; ok {
		return v, true
	}
	if v, ok := um.rawData[key]; ok {
		return v, true
	}
	return nil, false
}

func (um *UnknownMap) GetStr(key string) (string, bool) {
	if v, ok := um.data[key]; !ok {
		if i, ok := um.rawData[key]; !ok {
			return "", false
		} else {
			var s string
			if err := json.Unmarshal(i, &s); err != nil {
				return "", false
			}
			delete(um.rawData, key)
			um.data[key] = s
			return s, true
		}
	} else {
		if s, ok := v.(string); ok {
			return s, true
		} else {
			if j, ok := v.(json.RawMessage); ok {
				_ = json.Unmarshal(j, &s)
				return s, true
			}
			return "", false
		}
	}
}

func (um *UnknownMap) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if bv, err := driver.String.ConvertValue(value); err == nil {
		if s, ok := bv.(string); ok {
			v := []byte(s)
			return um.UnmarshalJSON(v)
		}
	} else {
		return err
	}
	return fmt.Errorf("Failed to scan UnknownMap")
}

func (um UnknownMap) Value() (driver.Value, error) {
	b, err := um.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (um UnknownMap) String() string {
	b, _ := um.MarshalJSON()
	fmt.Println(string(b))
	return string(b)
}
