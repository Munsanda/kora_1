package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	GroupName string `gorm:"size:100;not null;unique"`
}

func CreateGroup(db *gorm.DB, group *Group) error {
	return db.Create(group).Error
}

func GetGroupByID(db *gorm.DB, id uint) (*Group, error) {
	var group Group
	err := db.First(&group, id).Error
	return &group, err
}

func GetAllGroups(db *gorm.DB) ([]Group, error) {
	var groups []Group
	err := db.Find(&groups).Error
	return groups, err
}

// UpdateGroup updates the details of an existing group
func UpdateGroup(db *gorm.DB, group *Group) error {
	return db.Save(group).Error
}

func DeleteGroup(db *gorm.DB, id uint) error {
	return db.Delete(&Group{}, id).Error
}
