package common

import (
	"bytes"
	"database/sql/driver"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
)

type CSVArray struct {
	A []string
}

func (ca *CSVArray) Scan(value interface{}) error {
	if value == nil {
		ca.A = []string{}
		return nil
	}
	if bv, err := driver.String.ConvertValue(value); err == nil {
		if s, ok := bv.(string); ok {
			var err error
			s = strings.TrimPrefix(s, ",")
			s = strings.TrimSuffix(s, ",")
			ca.A, err = csv.NewReader(bytes.NewBuffer([]byte(s))).Read()
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		return err
	}
	return fmt.Errorf("Failed to scan CSVArray")
}

func (ca CSVArray) Value() (driver.Value, error) {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	if err := w.Write(ca.A); err != nil {
		return nil, err
	}
	w.Flush()
	s := buf.String()
	s = strings.TrimSuffix(s, "\n")
	if s != "" {
		s = "," + s + ","
	}
	return s, nil
}

func (ca *CSVArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(ca.A)
}
func (ca *CSVArray) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal([]byte(data), &(ca.A))
	return err
}
