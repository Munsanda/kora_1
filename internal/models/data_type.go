package models

import "gorm.io/gorm"

type DataType struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	DataType string `gorm:"size:50"`
}

func (DataType) TableName() string {
	return "data_types"
}

func CreateDataType(db *gorm.DB, dataType *DataType) error {
	return db.Create(dataType).Error
}

func GetDataType(db *gorm.DB, id uint) (*DataType, error) {
	var dataType DataType
	err := db.First(&dataType, id).Error
	return &dataType, err
}

func GetAllDataTypes(db *gorm.DB) ([]DataType, error) {
	var dataTypes []DataType
	err := db.Find(&dataTypes).Error
	return dataTypes, err
}
