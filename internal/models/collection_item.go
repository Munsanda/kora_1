package models

import "gorm.io/gorm"

type CollectionItem struct {
	ID                        uint   `gorm:"primaryKey;autoIncrement"`
	CollectionID              *uint  `gorm:"index"` // Nullable to match schema 'int' without 'not null' constraint, though practically FK usually implies existence
	CollectionItem            string `gorm:"size:50"`
	RelationCollectionItemsID *uint  `gorm:"index"`

	// Associations - optional but helpful, matching foreign keys
	Collection              *Collection     `gorm:"foreignKey:CollectionID"`
	RelationCollectionItems *CollectionItem `gorm:"foreignKey:RelationCollectionItemsID"`
}

func (CollectionItem) TableName() string {
	return "collection_items"
}

func CreateCollectionItem(db *gorm.DB, item *CollectionItem) error {
	return db.Create(item).Error
}

func GetCollectionItem(db *gorm.DB, id uint) (*CollectionItem, error) {
	var item CollectionItem
	err := db.Preload("Collection").Preload("RelationCollectionItems").First(&item, id).Error
	return &item, err
}
