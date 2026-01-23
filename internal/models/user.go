package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	FirstName  string    `gorm:"size:100"`
	MiddleName string    `gorm:"size:100"`
	Surname    string    `gorm:"size:100"`
	Dob        time.Time `gorm:"type:date"`
	Email      string    `gorm:"size:250;unique"`
	Password   string    `gorm:"size:250"`
}

func (User) TableName() string {
	return "users"
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func GetUser(db *gorm.DB, id uint) (*User, error) {
	var user User
	err := db.First(&user, id).Error
	return &user, err
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func UpdateUser(db *gorm.DB, user *User) error {
	return db.Save(user).Error
}

func DeleteUser(db *gorm.DB, id uint) error {
	return db.Delete(&User{}, id).Error
}
