package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type DBJsonRawMessage json.RawMessage

func (rm *DBJsonRawMessage) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if bv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := bv.(string); ok {
			b := []byte(v)
			*rm = b
			return nil
		}
	} else {
		return err
	}
	return fmt.Errorf("Failed to scan DBJsonRawMessage")
}

func (rm DBJsonRawMessage) Value() (driver.Value, error) {
	return string(rm), nil
}

func (rm *DBJsonRawMessage) Type() JSONType {
	switch string(*rm) {
	default:
		switch (*rm)[0] {
		case '"':
			return JSON_STRING
		case '[':
			return JSON_ARRAY
		case '{':
			return JSON_OBJECT
		default:
			//if parsable to float64 then number
			//_,err:=strconv.ParseFloat(string(*rm),64)
			//if err!=nil{}
			// as this is already a valid json then this should be a number only
			return JSON_NUMBER
		}
	case "true", "false":
		return JSON_BOOL
	case "null":
		return JSON_NULL

	}
}

func (rm DBJsonRawMessage) Json() (*Json, error) {
	t := rm.Type()
	j := Json{}
	switch t {
	case JSON_NULL:
		return &j, nil
	case JSON_BOOL:
		b := false
		j.val = &b
	case JSON_NUMBER:
		var f float64
		j.val = &f
	case JSON_STRING:
		s := ""
		j.val = &s
	case JSON_ARRAY:
		a := []*Json{}
		j.val = &a
	case JSON_OBJECT:
		o := map[string]*Json{}
		j.val = &o
	}
	err := json.Unmarshal(rm, j.val)
	return &j, err
}

func (rm *DBJsonRawMessage) UnmarshalJSON(data []byte) error {
	r := json.RawMessage{}
	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}
	*rm = DBJsonRawMessage(r)
	return nil
}
