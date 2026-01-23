package models

import (
	"time"

	"gorm.io/gorm"
)

type Submission struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	ServicesID *uint     `gorm:"index"`
	CreatedBy  *uint     `gorm:"index"`
	CreatedOn  time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	// Associations
	Service *Service `gorm:"foreignKey:ServicesID"`
	User    *User    `gorm:"foreignKey:CreatedBy"`
}

func (Submission) TableName() string {
	return "submissions"
}

func CreateSubmission(db *gorm.DB, submission *Submission) error {
	return db.Create(submission).Error
}

func GetSubmission(db *gorm.DB, id uint) (*Submission, error) {
	var submission Submission
	err := db.Preload("Service").Preload("User").First(&submission, id).Error
	return &submission, err
}

func UpdateSubmission(db *gorm.DB, submission *Submission) error {
	return db.Save(submission).Error
}

func DeleteSubmission(db *gorm.DB, id uint) error {
	return db.Delete(&Submission{}, id).Error
}
