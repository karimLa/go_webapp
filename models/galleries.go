package models

import "gorm.io/gorm"

type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not null;index"`
	Title  string `gorm:"not null"`
}

func NewGalleryService(db *gorm.DB) GalleryService {
	gg := &galleryGorm{db: db}

	return galleryService{
		GalleryDB: gg,
	}
}

type GalleryDB interface {
	Create(gallery *Gallery) error
}

type GalleryService interface {
	GalleryDB
}

type galleryService struct {
	GalleryDB
}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(g *Gallery) error {
	return nil
}
