package models

import (
	"github.com/soramon0/webapp/lib"
	"github.com/soramon0/webapp/utils"
	"gorm.io/gorm"
)

type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
}

func NewServices() *Services {
	db := lib.InitDB()
	us := NewUserService(db)
	gs := NewGalleryService(db)

	return &Services{
		db:      db,
		Gallery: gs,
		User:    us,
	}
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{})
}

func (s *Services) DestructiveReset() error {
	utils.Must(s.db.Migrator().DropTable(&User{}, &Gallery{}))
	return s.db.AutoMigrate(&User{}, &Gallery{})
}

// Close returns a not implemented error
func (s *Services) Close() error {
	return ErrNotImplemented
}
