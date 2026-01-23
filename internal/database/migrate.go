package database

import (
	"log"

	"kora_1/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {

	// Clean up orphaned records before adding constraints
	if db.Migrator().HasTable("form_fields") && db.Migrator().HasTable("fields") {
		db.Exec("DELETE FROM form_fields WHERE fields_id NOT IN (SELECT id FROM fields)")
	}

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
