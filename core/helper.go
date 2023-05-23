package core

import (
	"encoding/json"
	"html/template"
	"mocker/common"
)

func validateTemplate(ac common.AppContexter, templateStr string, renderType string) error {
	_, err := template.New("validator").Funcs(GetRespRenderFunctions(ac, renderType)).Parse(templateStr)
	return err
}

func validateFlowConfig(ac common.AppContexter, config common.FlowConfig) error {
	errs := common.NewMultiError("Config Validation failed")
	for k, v := range config {
		if err := validateApiTemplate(ac, v); err != nil {
			merr := err.(*common.ErrMultiErrors)
			for ek, es := range merr.Errors {
				errs.AddError(k+"_"+ek, es)
			}
		}
	}
	return errs.Err()
}

func validateApiTemplate(ac common.AppContexter, templates common.ApiTemplates) error {
	errs := common.NewMultiError("Api Template Validation failed")
	checkTemplate := func(k, t string) {
		if err := validateTemplate(ac, t, "RW"); err != nil {
			errs.AddError(k, err)
		}
	}

	checkTemplate("status", templates.Status)
	checkTemplate("header", templates.Header)
	checkTemplate("response_body", templates.ResponseBody)

	return errs.Err()
}

func GetRespRenderFunctions(ac common.AppContexter, renderType string) template.FuncMap {
	var funcMap = template.FuncMap{
		"Json": func(o interface{}) template.HTML {
			b, err := json.Marshal(o)
			if err != nil {
				ac.Logger().Errorln("Error while marshalling in Json obj:", "error:", err)
			}
			return template.HTML(string(b))
		},
		"Object": func(id int64) *common.Json {
			obj, err := GetObject(ac, id)
			if err != nil {
				ac.Logger().Errorln("Error while getting Object in template id=", id, "error:", err)
				return nil
			}
			if obj.IsEnabled() {
				return obj.Object()
			}
			ac.Logger().Errorln("Object not enabled id=", id)
			return nil
		},
		"Get": func(m *common.Json, keys ...interface{}) interface{} {
			return m.GetI(keys...)
		},
		"Set": func(m *common.Json, keysAndVals ...interface{}) string {
			err := m.Set(keysAndVals...)
			if err != nil {
				ac.Logger().Errorln("Error while Setting in object ", "error:", err)
			}
			return ""
		},
		"Del": func(m *common.Json, keysAndVals ...interface{}) string {
			err := m.Delete(keysAndVals...)
			if err != nil {
				ac.Logger().Errorln("Error while Deleting in object ", "error:", err)
			}
			return ""
		},
		"NewJsonArr": common.NewJsonArr,
		"NewJsonObj": common.NewJsonObj,
		"Copy": func(m *common.Json) *common.Json {
			res, err := common.NewJson(m)
			if err != nil {
				ac.Logger().Errorln("Error while Copying object ", "error:", err)
			}
			return res
		},
		"SaveObject": func(id int64) string {
			obj, err := GetObject(ac, id)
			if err != nil {
				ac.Logger().Errorln("Error while getting object id=", "error:", err)
				return ""
			}
			err = obj.Save(ac)
			if err != nil {
				ac.Logger().Errorln("Error while saving object", id, err)
			}
			return ""
		},
		"add": func(a ...interface{}) interface{} {
			var fs float64
			var is int64
			var isFloat bool
			for _, i := range a {
				switch v := i.(type) {
				case float64:
					fs += float64(v)
					isFloat = true
				case float32:
					fs += float64(v)
					isFloat = true
				case int:
					is += int64(v)
					fs += float64(v)
				case int64:
					is += int64(v)
					fs += float64(v)
				case int32:
					is += int64(v)
					fs += float64(v)
				case string:
					i:=common.StringToInt64(v)
					d:=common.StringToFloat64(v)
					isFloat = isFloat || i!=int64(d)
					is += i
					fs += d
				}
			}
			if isFloat {
				return fs
			} else {
				return is
			}
		},
		"subtract": func(a ...interface{}) interface{} {
			var fs float64
			var is int64
			var isFloat bool
			for _, i := range a {
				switch v := i.(type) {
				case float64:
					fs -= float64(v)
					isFloat = true
				case float32:
					fs -= float64(v)
					isFloat = true
				case int:
					is -= int64(v)
					fs -= float64(v)
				case int64:
					is -= int64(v)
					fs -= float64(v)
				case int32:
					is -= int64(v)
					fs -= float64(v)
				case string:
					i:=common.StringToInt64(v)
					d:=common.StringToFloat64(v)
					isFloat = isFloat || i!=int64(d)
					is -= i
					fs -= d
				}
			}
			if isFloat {
				return fs
			} else {
				return is
			}
		},
		"minus": func(a interface{}) interface{} {
			var fs float64
			var is int64
			var isFloat bool

			switch v := a.(type) {
			case float64:
				fs = float64(v)
				isFloat = true
			case float32:
				fs = float64(v)
				isFloat = true
			case int:
				is = int64(v)
			case int64:
				is = int64(v)
			case int32:
				is = int64(v)
			case string:
				i:=common.StringToInt64(v)
				d:=common.StringToFloat64(v)
				isFloat = isFloat || i!=int64(d)
				is = i
				fs = d
			}
			if isFloat {
				return -fs
			} else {
				return -is
			}
		},
		"mul": func(a ...interface{}) interface{} {
			var fs float64
			fs = 1
			var is int64
			is = 1
			var isFloat bool
			for _, i := range a {
				switch v := i.(type) {
				case float64:
					fs *= float64(v)
					isFloat = true
				case float32:
					fs *= float64(v)
					isFloat = true
				case int:
					is *= int64(v)
					fs *= float64(v)
				case int64:
					is *= int64(v)
					fs *= float64(v)
				case int32:
					is *= int64(v)
					fs *= float64(v)
				case string:
					i:=common.StringToInt64(v)
					d:=common.StringToFloat64(v)
					isFloat = isFloat || i!=int64(d)
					is *= i
					fs *= d
				}
				if fs == 0 || is == 0 {
					return 0
				}
			}
			if isFloat {
				return fs
			} else {
				return is
			}
		},
		"oneby": func(a interface{}) interface{} {
			var fs float64
			var is int64
			var isFloat bool

			switch v := a.(type) {
			case float64:
				fs = float64(v)
				isFloat = true
			case float32:
				fs = float64(v)
				isFloat = true
			case int:
				is = int64(v)
			case int64:
				is = int64(v)
			case int32:
				is = int64(v)
			case string:
				i:=common.StringToInt64(v)
				d:=common.StringToFloat64(v)
				isFloat = isFloat || i!=int64(d)
				is = i
				fs = d
			}
			if isFloat {
				return 1 / fs
			} else {
				return 1 / is
			}
		},
		"divide": func(a ...interface{}) interface{} {
			var fs float64
			fs = 1
			var is int64
			is = 1
			var isFloat bool
			for _, i := range a {
				switch v := i.(type) {
				case float64:
					fs /= float64(v)
					isFloat = true
				case float32:
					fs /= float64(v)
					isFloat = true
				case int:
					is /= int64(v)
					fs /= float64(v)
				case int64:
					is /= int64(v)
					fs /= float64(v)
				case int32:
					is /= int64(v)
					fs /= float64(v)
				case string:
					i:=common.StringToInt64(v)
					d:=common.StringToFloat64(v)
					isFloat = isFloat || i!=int64(d)
					is /= i
					fs /= d
				}
			}
			if isFloat {
				return fs
			} else {
				return is
			}
		},
	}
	if renderType == "RO" {
		delete(funcMap, "SaveObject")
		delete(funcMap, "Set")
		delete(funcMap, "Del")
		delete(funcMap, "Copy")
		funcMap["Object"] = func(id int64) *common.Json {
			obj, ok := ac.GetObject(id)
			if ok && obj != nil && obj.IsEnabled() {
				return obj.Object()
			}
			ac.Logger().Errorln("Object with id", id, "not loaded")
			return nil
		}
	}
	return funcMap
}
