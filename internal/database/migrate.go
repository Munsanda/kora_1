package database

import (
	"log"

	"kora_1/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {

	err := db.AutoMigrate(
		&models.User{},
		&models.Service{},
		&models.DataType{},
		&models.Group{},
		&models.Collection{},
		&models.CollectionItem{},
		&models.FormGroup{},
		&models.ReservedName{},
		&models.Field{},
		&models.Form{},
		&models.FormFields{},
		&models.Submission{},
		&models.FormAnswer{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migrated successfully")
}
