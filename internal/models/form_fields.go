package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FormFields struct {
	gorm.Model

	FormID      uint           `gorm:"index;not null"`
	FieldsID    uint           `gorm:"not null;index"`
	Validations datatypes.JSON `gorm:"type:jsonb"`
}

func CreateFormFields(db *gorm.DB, Fields *FormFields) error {
	return db.Create(&Fields).Error
}

func GetFormFields(db *gorm.DB, id uuid.UUID) (FormFields, error) {
	var Fields FormFields
	err := db.First(&Fields, "Fields_id = ?", id).Error
	return Fields, err
}

func UpdateFormFields(db *gorm.DB, Fields FormFields) error {
	return db.Save(&Fields).Error
}

func DeleteFormFields(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&FormFields{}, id).Error
}

func UpdateFormStatus(db *gorm.DB, formId uint, status bool) error {
	return db.Model(&Form{}).Where("id = ?", formId).Update("status", status).Error
}
