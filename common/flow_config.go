package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type FlowConfig map[string]ApiTemplates

type ApiTemplates struct {
	Status       string `json:"status" yaml:"status"`               // after Template execution this has to give a string parsable to in for parseInt
	Header       string `json:"header" yaml:"header"`               // After template execution this has to give a map[string][]string in Json
	ResponseBody string `json:"response_body" yaml:"response_body"` // After template execution this has to give a string but would depend on the content type
}

func (fc *FlowConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if bv, err := driver.String.ConvertValue(value); err == nil {
		if s, ok := bv.(string); ok {
			v := []byte(s)
			return json.Unmarshal(v, fc)
		}
	} else {
		return err
	}
	return fmt.Errorf("Failed to scan DBJsonRawMessage")
}

func (fc FlowConfig) Value() (driver.Value, error) {
	b, err := json.Marshal(fc)
	return string(b), err
}
