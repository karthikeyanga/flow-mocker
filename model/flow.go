package model

import (
	"mocker/common"

	"github.com/jinzhu/gorm"
)

const FlowTableName = "mocker_flow"

type FlowModel struct {
	BaseModel
	Title      string
	Identifier string
	ObjectList common.CSVArray
	Config     common.FlowConfig
}

func GetFlow(db *gorm.DB, id int64) (*FlowModel, error) {
	entry := FlowModel{}
	res := db.First(&entry, id)
	if res.RecordNotFound() {
		return nil, nil
	}
	err := res.Error
	return &entry, err
}
func (FlowModel) TableName() string {
	return FlowTableName
}

func (o *FlowModel) Save(db *gorm.DB) error {
	return db.Save(o).Error
}
