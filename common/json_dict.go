package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONSimpleStrDict map[string]string

func (jd *JSONSimpleStrDict) Scan(value interface{}) error {
	*jd = map[string]string{}
	if value == nil {
		return nil
	}
	if bv, err := driver.String.ConvertValue(value); err == nil {
		if s, ok := bv.(string); ok {
			b := []byte(s)
			return json.Unmarshal(b, &jd)
		}
	} else {
		return err
	}
	return fmt.Errorf("Failed to scan JSONSimpleStrDict")
}

func (jd JSONSimpleStrDict) Value() (driver.Value, error) {
	b, e := json.Marshal(jd)
	return string(b), e
}

//func (jd JSONSimpleStrDict) MarshalJSON() ([]byte, error) {
//	return json.Marshal(jd)
//}
//func (jd *JSONSimpleStrDict) UnmarshalJSON(data []byte) error {
//	err := json.Unmarshal([]byte(data), &(jd))
//	return err
//}
