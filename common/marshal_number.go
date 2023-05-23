package common

import (
	"math"
	"strconv"
)

func getFloatFromJsonByte(data []byte) (float64, error) {
	s := string(data)
	l := len(s)
	if l < 0 {
		return 0, nil
	}
	if s[0] == '"' {
		s = s[1:]
		l = l - 1
	}

	if s[l-1] == '"' {
		l = l - 1
		s = s[:l]
	}
	if l <= 0 {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

type JsonInt int

func (i *JsonInt) UnmarshalJSON(data []byte) error {

	if f, err := getFloatFromJsonByte(data); err != nil {
		return err
	} else {
		*i = JsonInt(math.Round(f))
	}
	return nil
}

type JsonInt64 int64

func (i *JsonInt64) UnmarshalJSON(data []byte) error {
	if f, err := getFloatFromJsonByte(data); err != nil {
		return err
	} else {
		*i = JsonInt64(math.Round(f))
	}
	return nil
}

type JsonInt32 int32

func (i *JsonInt32) UnmarshalJSON(data []byte) error {
	if f, err := getFloatFromJsonByte(data); err != nil {
		return err
	} else {
		*i = JsonInt32(math.Round(f))
	}
	return nil
}

type JsonFloat32 float32

func (i *JsonFloat32) UnmarshalJSON(data []byte) error {
	if f, err := getFloatFromJsonByte(data); err != nil {
		return err
	} else {
		*i = JsonFloat32(f)
	}
	return nil
}

type JsonFloat64 float64

func (i *JsonFloat64) UnmarshalJSON(data []byte) error {
	if f, err := getFloatFromJsonByte(data); err != nil {
		return err
	} else {
		*i = JsonFloat64(f)
	}
	return nil
}
