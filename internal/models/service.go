package models

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	Name string `gorm:"size:150;not null"`
}

// Create creates a new service
func (s *Service) Create(db *gorm.DB) error {
	return db.Create(s).Error
}

// Read retrieves a service by ID
func (s *Service) Read(db *gorm.DB, id uint) error {
	return db.First(s, id).Error
}

// Update updates the service
func (s *Service) Update(db *gorm.DB) error {
	return db.Save(s).Error
}

// Delete deletes the service
func (s *Service) Delete(db *gorm.DB) error {
	return db.Delete(s).Error
}
