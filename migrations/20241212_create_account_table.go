package migrations

import (
	"myAwesomeProject/internal/entities"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateAccountMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20241212_create_accounts_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&entities.Account{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("account")
		},
	}
}
