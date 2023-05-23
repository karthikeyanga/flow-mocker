package core

import (
	"encoding/json"
	"html/template"
	"mocker/common"
	"mocker/model"
)

type api struct {
	api         *model.ApisModel
	resTemplate *template.Template
}

func (a *api) init() (err error) {
	if a.resTemplate == nil {
		funcMap := GetRespRenderFunctions(nil, "RW") //The ac is nil as this is just a placeholder
		// function map. The actual will passed just before execute
		a.resTemplate, err = template.New("api_" + common.Int64ToString(a.api.ID)).Funcs(funcMap).Parse(a.api.ResponseBody)
	}
	return
}

func (a *api) Id() int64 {
	return a.api.ID
}

func (a *api) ResponseTemplate() *template.Template {
	return a.resTemplate
}

func (a *api) Method() string {
	return a.api.Method
}

func (a *api) Route() string {
	return a.api.Route
}
func (a *api) ResponseBody() string {
	return a.api.ResponseBody
}
func (a *api) ResponseHeaders() common.JSONSimpleStrDict {
	return a.api.ResponseHeaders
}
func (a *api) Status() int {
	return a.api.Status
}
func (a *api) IsEnabled() bool {
	return a.api.Enabled
}
func (a *api) Enable(ac common.AppContexter) error {
	a.api.Enabled = true
	return a.Save(ac)
}
func (a *api) Disable(ac common.AppContexter) error {
	a.api.Enabled = false
	return a.Save(ac)
}

func (a *api) getRep() map[string]interface{} {
	//preparing for marshal. Not all the data would be there in data. we might have to combine them
	d := map[string]interface{}{
		"id":               a.api.ID,
		"route":            a.api.Route,
		"method":           a.api.Method,
		"status":           a.api.Status,
		"response_headers": a.api.ResponseHeaders,
		"response_body":    a.api.ResponseBody,
	}
	return d
}
func (a api) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.getRep())
}
func (a api) MarshalYAML() (interface{}, error) {
	return a.getRep(), nil
}

func (a *api) Save(ac common.AppContexter) error {
	if err := a.api.Save(ac.DB()); err != nil {
		return err
	}
	return ac.Commit()
}

func (a *api) Update(ac common.AppContexter, route, method string, status int, respHeaders common.JSONSimpleStrDict, respBody string) error {
	if err := validateTemplate(ac, respBody, "RW"); err != nil {
		return err
	}
	a.resTemplate = nil
	a.api.Route = route
	a.api.Method = method
	a.api.Status = status
	a.api.ResponseHeaders = respHeaders
	a.api.ResponseBody = respBody
	err := a.init()
	if err != nil {
		return err
	}
	return a.Save(ac)
}

func GetApi(ac common.AppContexter, id int64) (common.IApi, error) {
	apiModel, err := model.GetApi(ac.DB(), id)
	if err != nil {
		return nil, err
	}
	if apiModel == nil {
		return nil, nil
	}
	a := api{
		api: apiModel,
	}
	return &a, a.init()
}

func GetApiFromModel(apiM model.ApisModel) (common.IApi, error) {
	a := api{
		api: &apiM,
	}
	return &a, a.init()
}

func NewApi(ac common.AppContexter, route, method string, status int, respHeaders common.JSONSimpleStrDict, respBody string) (common.IApi, error) {
	if err := validateTemplate(ac, respBody, "RW"); err != nil {
		return nil, err
	}
	a := api{
		api: &model.ApisModel{
			Route:           route,
			Method:          method,
			Status:          status,
			Enabled:         true,
			ResponseHeaders: respHeaders,
			ResponseBody:    respBody,
		},
	}
	if err := a.Save(ac); err != nil {
		return nil, err
	}
	return &a, a.init()
}
