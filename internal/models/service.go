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
func (s *Service) GetServiceByID(db *gorm.DB, id uint) error {
	return db.First(s, id).Error
}

// List retrieves all services
func (s *Service) ListServices(db *gorm.DB) ([]Service, error) {
	var services []Service
	return services, db.Find(&services).Error
}

// Update updates the service
func (s *Service) Update(db *gorm.DB) error {
	return db.Save(s).Error
}

// Delete deletes the service
func (s *Service) Delete(db *gorm.DB) error {
	return db.Delete(s).Error
}

// GetServiceByID retrieves a service by ID (standalone function)
func GetServiceByID(db *gorm.DB, id uint) (*Service, error) {
	var service Service
	if err := db.First(&service, id).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

// CreateService creates a new service (standalone function)
func CreateService(db *gorm.DB, name string) (*Service, error) {
	service := &Service{Name: name}
	if err := db.Create(service).Error; err != nil {
		return nil, err
	}
	return service, nil
}

// ListAllServices retrieves all services (standalone function)
func ListAllServices(db *gorm.DB) ([]Service, error) {
	var services []Service
	return services, db.Find(&services).Error
}

// UpdateService updates a service (standalone function)
func UpdateService(db *gorm.DB, service *Service) error {
	return db.Save(service).Error
}

// DeleteService deletes a service (standalone function)
func DeleteService(db *gorm.DB, id uint) error {
	return db.Delete(&Service{}, id).Error
}
