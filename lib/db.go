package lib

import (
	"github.com/soramon0/webapp/utils"
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

	return db
}
