package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	 "gorm.io/datatypes"

)

type BaseModel struct {
  ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Form struct {
  BaseModel

  Title       string `gorm:"size:150;not null"`
  Description string `gorm:"type:text"`
  Service     string `gorm:"size:100"`
  Status      int    `gorm:"default:1"`
  Version     int    `gorm:"default:1"`

  Questions []FormQuestion `gorm:"constraint:OnDelete:CASCADE"`
}

type Question struct {
  BaseModel

  Label string `gorm:"size:255;not null"`
  Type  string `gorm:"size:50;not null"`

  Meta datatypes.JSON `gorm:"type:jsonb"`

  Forms []FormQuestion
}

type FormQuestion struct {
  BaseModel

  FormID     uuid.UUID `gorm:"index;not null"`
  QuestionID uuid.UUID `gorm:"index;not null"`

  Position   int
  IsRequired bool

  Validations datatypes.JSON `gorm:"type:jsonb"`

  Question Question `gorm:"foreignKey:QuestionID"`
}


type Submission struct {
  BaseModel

  FormID      uuid.UUID `gorm:"index"`
  FormVersion int

  CreatedBy string `gorm:"size:100"`
  Status    string `gorm:"size:30;default:'submitted'"`

  Answers []FormAnswer `gorm:"constraint:OnDelete:CASCADE"`
}


type FormAnswer struct {
  BaseModel

  SubmissionID uuid.UUID `gorm:"index"`
  QuestionID   uuid.UUID `gorm:"index"`

  Answer     string
  AnswerJSON datatypes.JSON `gorm:"type:jsonb"`
}
