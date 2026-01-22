package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Submission struct {
	gorm.Model
	FormID      uuid.UUID `gorm:"index"`
	FormVersion int

	CreatedByUserId uint `gorm:"index"` 
	Status    string `gorm:"size:30;default:'submitted'"`

	Answers []FormAnswer `gorm:"constraint:OnDelete:CASCADE"`
}

func CreateSubmission(db *gorm.DB, submission *Submission) error {
	return db.Create(submission).Error
}

func GetSubmission(db *gorm.DB, id uuid.UUID) (*Submission, error) {
	var submission Submission
	err := db.First(&submission, id).Error
	return &submission, err
}

func UpdateSubmission(db *gorm.DB, submission *Submission) error {
	return db.Save(submission).Error
}

func DeleteSubmission(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&Submission{}, id).Error
}
