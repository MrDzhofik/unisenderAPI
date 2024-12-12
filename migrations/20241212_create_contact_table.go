package migrations

import (
	"myAwesomeProject/internal/entities"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateContactMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20241212_create_contacts_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&entities.Contact{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("contact")
		},
	}
}
