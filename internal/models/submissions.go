package models

import (
	"gorm.io/gorm"
)

type Submission struct {
	gorm.Model
	FormID      uint `gorm:"index"`
	FormVersion int

	CreatedByUserId uint   `gorm:"index"`
	Status          string `gorm:"size:30;default:'submitted'"`

	Answers []FormAnswer `gorm:"constraint:OnDelete:CASCADE"`
}

func CreateSubmission(db *gorm.DB, submission *Submission) error {
	return db.Create(submission).Error
}

func GetSubmission(db *gorm.DB, id uint) (*Submission, error) {
	var submission Submission
	err := db.Preload("Answers").First(&submission, id).Error
	return &submission, err
}

func GetSubmissionsByFormID(db *gorm.DB, formID uint) ([]Submission, error) {
	var submissions []Submission
	err := db.Preload("Answers").Where("form_id = ?", formID).Find(&submissions).Error
	return submissions, err
}

func UpdateSubmission(db *gorm.DB, submission *Submission) error {
	return db.Save(submission).Error
}

func DeleteSubmission(db *gorm.DB, id uint) error {
	return db.Delete(&Submission{}, id).Error
}
