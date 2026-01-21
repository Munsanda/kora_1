package database

import (
	"log"

	"kora_1/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {

	err := db.AutoMigrate(
		&models.Form{},
		&models.Question{},
		&models.FormQuestion{},
		&models.Submission{},
		&models.FormAnswer{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migrated successfully")
}
