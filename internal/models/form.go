package models

import "gorm.io/gorm"

type Form struct {
	gorm.Model
	Title       string `gorm:"size:150;not null"`
	Description string `gorm:"type:text"`
	ServiceId   int    `gorm:"size:100"`
	Status      int    `gorm:"default:0"`
	Version     int    `gorm:"default:1"`

	FormFields []FormFields `gorm:"foreignKey:FormID"`
}

func CreateForm(db *gorm.DB, form *Form) (*Form, error) {
	err := db.Create(form).Error
	return form, err
}

func CreateFormWithMutlipleFields(db *gorm.DB, form *Form, fields []Field) (*Form, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(form).Error; err != nil {
			return err
		}
		if err := tx.Model(form).Association("Fields").Append(fields); err != nil {
			return err
		}
		return nil
	})
	return form, err
}

// In models/form_model.go
func GetForm(db *gorm.DB, id uint) (Form, error) {
    var form Form

    result := db.
		Preload("FormFields").
		Preload("FormFields.Field").
        First(&form, id)

    if result.Error != nil {
        return Form{}, result.Error
    }

    return form, nil
}


func UpdateForm(db *gorm.DB, form *Form) error {
	return db.Save(form).Error
}

func DeleteForm(db *gorm.DB, id uint) error {
	return db.Delete(&Form{}, id).Error
}
