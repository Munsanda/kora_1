package models

import "gorm.io/gorm"

type ReservedName struct {
	gorm.Model
	Name string `gorm:"size:100;not null;unique"`
}

func CreateReservedName(db *gorm.DB, reservedName *ReservedName) error {
	return db.Create(reservedName).Error
}

func GetSimilarNames(db *gorm.DB, name string) ([]ReservedName, error) {
	var reservedNames []ReservedName
	err := db.Where("name LIKE ?", "%"+name+"%").Find(&reservedNames).Error
	return reservedNames, err
}
