package models

import "gorm.io/gorm"

type FormAnswer struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	FormFieldID  *uint  `gorm:"index"`
	Answer       string `gorm:"size:250"`
	SubmissionID *uint  `gorm:"index"`

	// Associations
	FormField  *FormFields `gorm:"foreignKey:FormFieldID"`
	Submission *Submission `gorm:"foreignKey:SubmissionID"`
}

func (FormAnswer) TableName() string {
	return "form_answers"
}

func CreateFormAnswer(db *gorm.DB, answer *FormAnswer) error {
	return db.Create(answer).Error
}

func GetFormAnswer(db *gorm.DB, id uint) (*FormAnswer, error) {
	var answer FormAnswer
	err := db.Preload("FormField").Preload("Submission").First(&answer, id).Error
	return &answer, err
}

func UpdateFormAnswer(db *gorm.DB, answer *FormAnswer) error {
	return db.Save(answer).Error
}

func DeleteFormAnswer(db *gorm.DB, id uint) error {
	return db.Delete(&FormAnswer{}, id).Error
}
