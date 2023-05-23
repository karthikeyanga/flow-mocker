package core

import (
	"encoding/json"
	"mocker/common"
	"mocker/model"
)

type object struct {
	object *model.ObjectModel
}

func (o *object) Id() int64 {
	return o.object.ID
}

func (o *object) Object() *common.Json {
	return &o.object.Object
}
func (o *object) Update(ac common.AppContexter, m *common.Json) error {
	o.object.Object = *m
	return o.Save(ac)
}
func (o *object) IsEnabled() bool {
	return o.object.Enabled
}
func (o *object) Enable(ac common.AppContexter) error {
	o.object.Enabled = true
	return o.Save(ac)
}
func (o *object) Disable(ac common.AppContexter) error {
	o.object.Enabled = false
	ac.UnloadObject(o.Id())
	return o.Save(ac)
}
func (o *object) Save(ac common.AppContexter) error {
	if err := o.object.Save(ac.DB()); err != nil {
		return err
	}
	return ac.Commit()
}
func (o *object) getRep() map[string]interface{} {
	//preparing for marshal. Not all the data would be there in data. we might have to combine them
	d := map[string]interface{}{
		"id":     o.object.ID,
		"object": o.object.Object,
	}
	return d
}
func (o object) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.getRep())
}

func (o object) MarshalYAML() (interface{}, error) {
	return o.getRep(), nil
}

func GetObject(ac common.AppContexter, id int64) (common.IObject, error) {
	if obj, ok := ac.GetObject(id); ok {
		return obj, nil
	}
	obj, err := model.GetObject(ac.DB(), id)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, nil
	}
	return &object{object: obj}, nil
}

func NewObject(ac common.AppContexter, m *common.Json) (common.IObject, error) {
	if m.Type() == common.JSON_NULL {
		return nil, nil
	}
	obj := object{
		object: &model.ObjectModel{
			Enabled: true,
			Object:  *m,
		},
	}
	if err := obj.Save(ac); err != nil {
		return nil, err
	}
	return &obj, nil
}
