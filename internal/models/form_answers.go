package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FormAnswer struct {
	gorm.Model
	SubmissionID uuid.UUID `gorm:"index"`
	QuestionID   uuid.UUID `gorm:"index"`
	Answer       string
	AnswerJSON   datatypes.JSON `gorm:"type:jsonb"`
}

func (fa *FormAnswer) Create(db *gorm.DB) error {
	return db.Create(fa).Error
}

func (fa *FormAnswer) Read(db *gorm.DB, id uuid.UUID) error {
	return db.First(fa, "id = ?", id).Error
}

func (fa *FormAnswer) Update(db *gorm.DB) error {
	return db.Save(fa).Error
}

func (fa *FormAnswer) Delete(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&FormAnswer{}, "id = ?", id).Error
}

func (fa *FormAnswer) FindBySubmissionID(db *gorm.DB, submissionID uuid.UUID) ([]FormAnswer, error) {
	var answers []FormAnswer
	err := db.Where("submission_id = ?", submissionID).Find(&answers).Error
	return answers, err
}

func (fa *FormAnswer) FindByQuestionID(db *gorm.DB, questionID uuid.UUID) ([]FormAnswer, error) {
	var answers []FormAnswer
	err := db.Where("question_id = ?", questionID).Find(&answers).Error
	return answers, err
}
