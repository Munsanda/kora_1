package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	GroupName string `gorm:"size:100;not null;unique"`
}

func CreateGroup(db *gorm.DB, group *Group) error {
	return db.Create(group).Error
}

func AddFieldsToGroup(db *gorm.DB, formID uint, groupID uint, fieldIDs []uint) error {
	// Update the group_id for the specified fields in the form
	// Only updates fields that exist in the form - fields that don't exist are simply skipped
	return db.Model(&FormFields{}).
		Where("form_id = ? AND fields_id IN ?", formID, fieldIDs).
		Update("group_id", groupID).Error
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

// GetAllGroupFields retrieves all fields for a specific group in a specific form
func GetAllGroupFields(db *gorm.DB, formID uint, groupID uint) ([]FormFields, error) {
	var formFields []FormFields
	err := db.Where("form_id = ? AND group_id = ?", formID, groupID).
		Find(&formFields).Error
	return formFields, err
}
