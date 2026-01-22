package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Fields struct {
	gorm.Model
	Label string `gorm:"size:255;not null"`
	Type  string `gorm:"size:50;not null"`

	Meta datatypes.JSON `gorm:"type:jsonb"`

	Forms []FormFields

	IsRequired bool `gorm:"default:false"`
	
}

func CreateFields(db *gorm.DB, Fields *Fields) error {
	return db.Create(Fields).Error
}

func GetFields(db *gorm.DB, id uint) (*Fields, error) {
	var Fields Fields
	err := db.First(&Fields, id).Error
	return &Fields, err
}

func UpdateFields(db *gorm.DB, Fields *Fields) error {
	return db.Save(Fields).Error
}

func DeleteFields(db *gorm.DB, id uint) error {
	return db.Delete(&Fields{}, id).Error
}
