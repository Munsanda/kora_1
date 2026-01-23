package models

import "gorm.io/gorm"

type Field struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Label        string `gorm:"size:50;not null"`
	DataTypeID   uint   `gorm:"not null"`
	GroupID      *uint  `gorm:"index"`
	CollectionID *uint  `gorm:"index"`
	Status       *bool  `gorm:"default:null"` // Using pointer for nullable boolean

	// Associations
	DataType   DataType    `gorm:"foreignKey:DataTypeID"`
	Group      *Group      `gorm:"foreignKey:GroupID"`
	Collection *Collection `gorm:"foreignKey:CollectionID"`
}

func (Field) TableName() string {
	return "fields"
}

func CreateFields(db *gorm.DB, field *Field) error {
	return db.Create(field).Error
}

func GetFields(db *gorm.DB, id uint) (*Field, error) {
	var field Field
	err := db.Preload("DataType").Preload("Group").Preload("Collection").First(&field, id).Error
	return &field, err
}

func UpdateFields(db *gorm.DB, field *Field) error {
	return db.Save(field).Error
}

func DeleteFields(db *gorm.DB, id uint) error {
	return db.Delete(&Field{}, id).Error
}
