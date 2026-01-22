package database

import (
	"log"

	"kora_1/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {

	err := db.AutoMigrate(
		&models.Form{},
		&models.Field{},
		&models.FormFields{},
		&models.Submission{},
		&models.FormAnswer{},
		&models.Group{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migrated successfully")
}
