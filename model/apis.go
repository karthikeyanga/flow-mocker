package model

import (
	"mocker/common"

	"github.com/jinzhu/gorm"
)

const ApisTableName = "mocker_api"

type ApisModel struct {
	BaseModel
	Route           string
	Method          string
	Status          int
	Enabled         bool
	ResponseHeaders common.JSONSimpleStrDict
	ResponseBody    string
}

func GetAllEnabledApis(db *gorm.DB) ([]ApisModel, error) {
	entries := []ApisModel{}
	res := db.Where("enabled = ?", true).Find(&entries)
	if res.RecordNotFound() {
		return nil, nil
	}
	err := res.Error
	return entries, err
}

func (ApisModel) TableName() string {
	return ApisTableName
}

func (o *ApisModel) Save(db *gorm.DB) error {
	return db.Save(o).Error
}

func GetApi(db *gorm.DB, id int64) (*ApisModel, error) {
	entry := ApisModel{}
	res := db.First(&entry, id)
	if res.RecordNotFound() {
		return nil, nil
	}
	err := res.Error
	return &entry, err
}
