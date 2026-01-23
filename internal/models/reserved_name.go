package models

import "gorm.io/gorm"

type ReservedName struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	ReservedName string `gorm:"size:50"`
}

func (ReservedName) TableName() string {
	return "reserved_names"
}

func CreateReservedName(db *gorm.DB, reservedName *ReservedName) error {
	return db.Create(reservedName).Error
}

func GetSimilarNames(db *gorm.DB, name string) ([]ReservedName, error) {
	var reservedNames []ReservedName
	err := db.Where("reserved_name LIKE ?", "%"+name+"%").Find(&reservedNames).Error
	return reservedNames, err
}
