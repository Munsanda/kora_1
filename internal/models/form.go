package models

import "gorm.io/gorm"

type Form struct {
	gorm.Model
	Title       string `gorm:"size:150;not null"`
	Description string `gorm:"type:text"`
	ServiceId   int    `gorm:"size:100"`
	Status      int    `gorm:"default:1"`
	Version     int    `gorm:"default:1"`

	Fieldss []FormFields `gorm:"constraint:OnDelete:CASCADE"`
}

func CreateForm(db *gorm.DB, form *Form) (*Form, error) {
	err := db.Create(form).Error
	return form, err
}

func GetForm(db *gorm.DB, id uint) (*Form, error) {
	var form Form
	if err := db.First(&form, id).Error; err != nil {
		return nil, err
	}
	return &form, nil
}

func UpdateForm(db *gorm.DB, form *Form) error {
	return db.Save(form).Error
}

func DeleteForm(db *gorm.DB, id uint) error {
	return db.Delete(&Form{}, id).Error
}
