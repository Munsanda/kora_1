package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Field struct {
	gorm.Model
	Label string `gorm:"size:255;not null"`
	Type  string `gorm:"size:50;not null"`

	Meta datatypes.JSON `gorm:"many2many:form_fields;"`
	
	IsRequired bool `gorm:"default:false"`
}

func CreateFields(db *gorm.DB, Fields *Field) error {
	return db.Create(Fields).Error
}

func GetFields(db *gorm.DB, id uint) (*Field, error) {
	var Fields Field
	err := db.First(&Fields, id).Error
	return &Fields, err
}

func UpdateFields(db *gorm.DB, Fields *Field) error {
	return db.Save(Fields).Error
}

func DeleteFields(db *gorm.DB, id uint) error {
	return db.Delete(&Field{}, id).Error
}
