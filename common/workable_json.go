package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type JSONType int

const (
	JSON_NULL = iota
	JSON_BOOL
	JSON_NUMBER
	JSON_STRING
	JSON_ARRAY
	JSON_OBJECT
)

var JsonTypeToStr = map[JSONType]string{
	JSON_NULL:   "NULL",
	JSON_BOOL:   "BOOL",
	JSON_NUMBER: "NUMBER",
	JSON_STRING: "STRING",
	JSON_ARRAY:  "ARRAY",
	JSON_OBJECT: "OBJECT",
}

type ErrJsonNotIndexable JSONType

func (e *ErrJsonNotIndexable) Error() string {
	return JsonTypeToStr[JSONType(*e)]
}

func NewErrJsonNotIndexable(t JSONType) error {
	e := ErrJsonNotIndexable(t)
	return &e
}

type Json struct {
	val interface{}
}

func (j *Json) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if bv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := bv.(string); ok {
			b := []byte(v)
			return json.Unmarshal(b, j)
		}
	} else {
		return err
	}
	return fmt.Errorf("Failed to scan Json")
}

func (j Json) Value() (driver.Value, error) {
	b, err := json.Marshal(j.val)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (j *Json) Get(keys ...interface{}) (*Json, error) {
	if j == nil {
		return nil, nil
	}
	cur := j
	//keys could be string>dict or numbers>array
	for _, key := range keys {
		if cur == nil {
			return nil, nil
		}
		t := cur.Type()
		switch t {
		case JSON_OBJECT:
			cur = (*cur.val.(*map[string]*Json))[key.(string)]
		case JSON_ARRAY:
			index := key.(int)
			arr := (*cur.val.(*[]*Json))
			if index < 1 || index > len(arr) {
				return nil, nil
			}
			cur = arr[index-1]
		default:
			return nil, NewErrJsonNotIndexable(t)
		}
	}
	return cur, nil
}

func (j *Json) Set(keysAndVal ...interface{}) error {
	if j == nil {
		return nil
	}
	lk := len(keysAndVal)
	if lk < 2 {
		return nil
	}
	val, err := NewJson(keysAndVal[lk-1])
	if err != nil {
		return err
	}
	keys := keysAndVal[:lk-2]
	cur, err := j.Get(keys...)
	if err != nil {
		return err
	}
	t := cur.Type()
	switch t {
	case JSON_OBJECT:
		(*cur.val.(*map[string]*Json))[keysAndVal[lk-2].(string)] = val
	case JSON_ARRAY:
		arr := (cur.val.(*[]*Json))
		index := keysAndVal[lk-2].(int)
		l := len(*arr)
		switch {
		case index == 0 || index == -1: //insert to start
			*arr = append([]*Json{val}, *(arr)...)
		case l+1 <= index || l+1 <= -index: // index here is not yet corrected for 0 based indexing from 1 base indexing
			*arr = append(*arr, val)
		case index < -1:
			index = -index - 1
			t := append([]*Json{val}, (*arr)[index:]...)
			*arr = append((*arr)[:index], t...)
		default:
			(*arr)[index-1] = val //correcting index from 1 based to 0 based
		}

	default:
		return NewErrJsonNotIndexable(t)
	}
	return nil
}

func (j *Json) Delete(keysAndVal ...interface{}) error {
	if j == nil {
		return nil
	}
	lk := len(keysAndVal)
	if lk < 1 {
		return nil
	}
	keys := keysAndVal[:lk-1]
	cur, err := j.Get(keys...)
	if err != nil {
		return err
	}
	t := cur.Type()
	switch t {
	case JSON_OBJECT:
		delete(*cur.val.(*map[string]*Json), keysAndVal[lk-1].(string))
	case JSON_ARRAY:
		arr := (cur.val.(*[]*Json))
		index := keysAndVal[lk-1].(int)
		l := len(*arr)
		if index < 1 || index > l {
			return nil
		}
		index--
		*arr = append((*arr)[:index], (*arr)[index+1:]...)
	default:
		return NewErrJsonNotIndexable(t)
	}
	return nil
}

func (j *Json) GetI(keys ...interface{}) interface{} {
	//the  value can be true, false, null, or start with {,[,"  else number parsable to float64
	o, err := j.Get(keys...)
	if err != nil {
		return nil
	}
	switch o.Type() {
	case JSON_BOOL:
		b := o.val.(*bool)
		return *b
	case JSON_NUMBER:
		n := o.val.(*float64)
		return *n
	case JSON_STRING:
		s := o.val.(*string)
		return *s
	case JSON_NULL:
		return nil
	default:
		return o.val
	}
}

func (j *Json) Type() JSONType {
	if j == nil || j.val == nil {
		return JSON_NULL
	}
	switch j.val.(type) {
	case *bool:
		return JSON_BOOL
	case *float64:
		return JSON_NUMBER
	case *string:
		return JSON_STRING
	case *[]*Json:
		return JSON_ARRAY
	default: //case map[string]*Json:
		return JSON_OBJECT
	}
}

func (j Json) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.val)
}
func (j *Json) UnmarshalJSON(data []byte) error {
	rm := DBJsonRawMessage{}
	err := json.Unmarshal(data, &rm)
	if err != nil {
		return err
	}
	jp, err := rm.Json()
	if err != nil {
		return err
	}
	*j = *jp
	return err
}

func (j Json) String() string {
	sb := strings.Builder{}
	switch v := j.val.(type) {
	case *bool:
		sb.WriteString(fmt.Sprint(*v))
	case *float64:
		sb.WriteString(fmt.Sprint(*v))
	case *string:
		//sb.WriteString("\"")
		sb.WriteString(fmt.Sprint(*v))
		//sb.WriteString("\"")
	case *[]*Json:
		sb.WriteString("[")
		flag := false
		for _, o := range *v {
			if flag {
				sb.WriteString(",")
			}
			flag = true
			sb.WriteString(fmt.Sprint(o))
		}
		sb.WriteString("]")
	case *map[string]*Json:
		sb.WriteString("{")
		flag := false
		for k, o := range *v {
			if flag {
				sb.WriteString(",")
			}
			flag = true
			sb.WriteString("\"")
			sb.WriteString(k)
			sb.WriteString("\":")
			sb.WriteString(fmt.Sprint(o))
		}
		sb.WriteString("]")
	}
	return sb.String()
}

func (j Json) MarshalYAML() (interface{}, error) {
	return j.val, nil
}

func NewJson(v interface{}) (*Json, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	rm := DBJsonRawMessage{}
	err = json.Unmarshal(b, &rm)
	if err != nil {
		return nil, err
	}
	return rm.Json()
}

func NewJsonArr() *Json {
	return &Json{
		val: &[]*Json{},
	}
}
func NewJsonObj() *Json {
	return &Json{
		val: &map[string]*Json{},
	}
}
