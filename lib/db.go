package lib

import (
	"github.com/karimla/webapp/models"
	"github.com/karimla/webapp/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	psqlDialector := postgres.New(postgres.Config{
		DSN:                  utils.GetDB(),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	})
	db, err := gorm.Open(psqlDialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	migrate(db)

	return db
}

func migrate(db *gorm.DB) {
	utils.Must(db.AutoMigrate(&models.User{}))
}
