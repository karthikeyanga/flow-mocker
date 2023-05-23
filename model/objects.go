package model

import (
	"mocker/common"

	"github.com/jinzhu/gorm"
)

const ObjectTableName = "mocker_object"

type ObjectModel struct {
	BaseModel
	Enabled bool
	//Format string
	Object common.Json
}

func GetObject(db *gorm.DB, id int64) (*ObjectModel, error) {
	entry := ObjectModel{}
	res := db.First(&entry, id)
	if res.RecordNotFound() {
		return nil, nil
	}
	err := res.Error
	return &entry, err
}

func (ObjectModel) TableName() string {
	return ObjectTableName
}

func (o *ObjectModel) Save(db *gorm.DB) error {
	return db.Save(o).Error
}
