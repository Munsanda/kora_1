package models

import "gorm.io/gorm"

type FormFields struct { // Name matches 'form_fields' table but struct convention usually Singular. However keeping as FormFields to match file/usage context or TableName. I will use 'FormField' singular for the struct type if possible, but keep filenames. The previous file used FormFields. Let's start using singular 'FormField' for struct but map to 'form_fields'.

	ID          uint   `gorm:"primaryKey;autoIncrement"`
	FormID      uint   `gorm:"not null"`
	FieldID     uint   `gorm:"not null"`
	FieldName   string `gorm:"size:50"`
	FormGroupID *uint  `gorm:"index"`
	Validation  string `gorm:"size:250"`
	FieldSpan   int
	FieldRow    int

	// Associations
	Form      Form       `gorm:"foreignKey:FormID"`
	Field     Field      `gorm:"foreignKey:FieldID"`
	FormGroup *FormGroup `gorm:"foreignKey:FormGroupID"`
}

func (FormFields) TableName() string {
	return "form_fields"
}

func CreateFormFields(db *gorm.DB, formField *FormFields) error {
	return db.Create(formField).Error
}

func GetFormFields(db *gorm.DB, id uint) (*FormFields, error) {
	var formField FormFields
	err := db.Preload("Form").Preload("Field").Preload("FormGroup").First(&formField, id).Error
	return &formField, err
}

func UpdateFormFields(db *gorm.DB, formField *FormFields) error {
	return db.Save(formField).Error
}

func DeleteFormFields(db *gorm.DB, id uint) error {
	return db.Delete(&FormFields{}, id).Error
}
