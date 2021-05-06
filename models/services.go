package models

import (
	"webapp/lib"
	"webapp/utils"

	"gorm.io/gorm"
)

type Services struct {
	Gallery GalleryService
	Image   ImageService
	User    UserService
	db      *gorm.DB
}

func NewServices() *Services {
	db := lib.InitDB()
	us := NewUserService(db)
	is := NewImageService()
	gs := NewGalleryService(db)

	return &Services{
		db:      db,
		Image:   is,
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
