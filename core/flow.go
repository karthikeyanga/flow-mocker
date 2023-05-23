package core

import (
	"encoding/json"
	"html/template"
	"mocker/common"
	"mocker/model"
	"time"
)

type flow struct {
	flow           *model.FlowModel
	identifier     *template.Template
	lastAccessedAt time.Time
}

func (f *flow) init(ac common.AppContexter) (err error) {
	if f.identifier == nil {
		funcMap := GetRespRenderFunctions(nil, "RO") //The ac is nil as this is just a placeholder
		// function map. The actual will passed just before execute
		f.identifier, err = template.New("flow_identifier_" + common.Int64ToString(f.flow.ID)).Funcs(funcMap).Parse(f.flow.Identifier)
		f.lastAccessedAt = time.Now()
	}
	return
}

func (f *flow) Id() int64 {
	return f.flow.ID
}

func (f *flow) Accessed() {
	f.lastAccessedAt = time.Now()
}

func (f *flow) Config() common.FlowConfig {
	return f.flow.Config
}

func (f *flow) LastAccessedAt() time.Time {
	return f.lastAccessedAt
}
func (f *flow) ObjectIds() []string {
	return f.flow.ObjectList.A
}
func (f *flow) IdentifierString() string {
	return f.flow.Identifier
}
func (f *flow) Title() string {
	return f.flow.Title
}
func (f *flow) IdentifierTemplate() *template.Template {
	return f.identifier
}
func (f *flow) Load(ac common.AppContexter) error {
	ac.LoadFlow(f)
	for _, sid := range f.ObjectIds() {
		id := common.StringToInt64(sid)
		obj, _ := GetObject(ac, id)
		if obj != nil && obj.IsEnabled() {
			ac.SetObject(id, obj)
		}
	}
	return nil
}
func (f *flow) Unload(ac common.AppContexter) error {
	ac.UnloadFlow(f)
	for _, sid := range f.ObjectIds() {
		id := common.StringToInt64(sid)
		ac.UnloadObject(id)
	}
	return nil
}

func (f *flow) getRep() map[string]interface{} {
	//preparing for marshal. Not all the data would be there in data. we might have to combine them
	d := map[string]interface{}{
		"id":         f.flow.ID,
		"title":      f.flow.Title,
		"identifier": f.flow.Identifier,
		"objects":    f.flow.ObjectList.A,
		"config":     f.flow.Config,
	}
	return d
}

func (f flow) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.getRep())
}

func (f flow) MarshalYAML() (interface{}, error) {
	return f.getRep(), nil
}

func (f *flow) Save(ac common.AppContexter) error {
	if err := f.flow.Save(ac.DB()); err != nil {
		return err
	}
	return ac.Commit()
}

func (f *flow) Update(ac common.AppContexter, title, identifier string, objIds []string, config common.FlowConfig) error {
	if err := validateFlow(ac, identifier, config); err != nil {
		return err
	}
	f.identifier = nil
	f.flow.Title = title
	f.flow.Identifier = identifier
	f.flow.Config = config
	f.flow.ObjectList = common.CSVArray{A: objIds}
	if err := f.init(ac); err != nil {
		_ = f.Unload(ac)
		return err
	}
	return f.Save(ac)
}

func GetFlow(ac common.AppContexter, id int64) (common.IFlow, error) {
	flowModel, err := model.GetFlow(ac.DB(), id)
	if err != nil {
		return nil, err
	}
	if flowModel == nil {
		return nil, nil
	}
	f := flow{
		flow: flowModel,
	}
	return &f, f.init(ac)
}

func NewFlow(ac common.AppContexter, title, identifier string, objIds []string, config common.FlowConfig) (common.IFlow, error) {
	if err := validateFlow(ac, identifier, config); err != nil {
		return nil, err
	}
	f := flow{
		flow: &model.FlowModel{
			Title:      title,
			Identifier: identifier,
			ObjectList: common.CSVArray{A: objIds},
			Config:     config,
		},
	}
	if err := f.Save(ac); err != nil {
		return nil, err
	}
	return &f, f.init(ac)
}
func validateFlow(ac common.AppContexter, identifier string, config common.FlowConfig) error {
	if err := validateTemplate(ac, identifier, "RO"); err != nil {
		return err
	}
	if err := validateFlowConfig(ac, config); err != nil {
		return err
	}
	return nil
}
