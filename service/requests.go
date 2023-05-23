package service

import "mocker/common"

type FlowRequest struct {
	Title      string             `yaml:"title"`
	Identifier string             `yaml:"identifier"`
	Objects    []common.JsonInt64 `yaml:"objects"`
	Config     common.FlowConfig  `yaml:"config"`
}

func (f *FlowRequest) ObjectIds() []string {
	res := []string{}
	for _, i := range f.Objects {
		res = append(res, common.Int64ToString(int64(i)))
	}
	return res
}

type ApiRequest struct {
	Route           string                   `json:"route"`
	Method          string                   `json:"method"`
	Status          int                      `json:"status"`
	ResponseHeaders common.JSONSimpleStrDict `json:"response_headers"`
	ResponseBody    string                   `json:"response_body"`
}
