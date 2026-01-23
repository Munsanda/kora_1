package models

import "gorm.io/gorm"

type FormGroup struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	GroupName string `gorm:"size:50"`
	GroupSpan int
	GroupRow  int
}

func (FormGroup) TableName() string {
	return "form_groups"
}

func CreateFormGroup(db *gorm.DB, formGroup *FormGroup) error {
	return db.Create(formGroup).Error
}

func GetFormGroup(db *gorm.DB, id uint) (*FormGroup, error) {
	var formGroup FormGroup
	err := db.First(&formGroup, id).Error
	return &formGroup, err
}

func GetAllFormGroups(db *gorm.DB) ([]FormGroup, error) {
	var groups []FormGroup
	err := db.Find(&groups).Error
	return groups, err
}
