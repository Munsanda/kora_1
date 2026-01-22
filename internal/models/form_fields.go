package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FormFields struct {
	gorm.Model

	FormID   uuid.UUID `gorm:"index;not null"`
	FieldsID uuid.UUID `gorm:"index;not null"`

	Position   int
	IsRequired bool

	Validations datatypes.JSON `gorm:"type:jsonb"`

	Fields Fields `gorm:"foreignKey:FieldsID"`
}

func CreateFormFields(db *gorm.DB, Fields FormFields) (FormFields, error) {
	err := db.Create(&Fields).Error
	return Fields, err
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
