package models

import "gorm.io/gorm"

type Collection struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	CollectionName string `gorm:"size:50"`
}

func (Collection) TableName() string {
	return "collections"
}

func CreateCollection(db *gorm.DB, collection *Collection) error {
	return db.Create(collection).Error
}

func GetCollection(db *gorm.DB, id uint) (*Collection, error) {
	var collection Collection
	err := db.First(&collection, id).Error
	return &collection, err
}

func GetAllCollections(db *gorm.DB) ([]Collection, error) {
	var collections []Collection
	err := db.Find(&collections).Error
	return collections, err
}
