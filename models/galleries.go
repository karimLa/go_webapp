package models

import "gorm.io/gorm"

type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not null;index"`
	Title  string `gorm:"not null"`
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

func NewGalleryService(db *gorm.DB) GalleryService {
	gg := newGalleryGorm(db)
	gv := newGalleryValidator(gg)

	return &galleryService{
		GalleryDB: gv,
	}
}

type galleryValidator struct {
	GalleryDB
}

func newGalleryValidator(gg *galleryGorm) *galleryValidator {
	return &galleryValidator{
		GalleryDB: gg,
	}
}

type galleryGorm struct {
	db *gorm.DB
}

func newGalleryGorm(db *gorm.DB) *galleryGorm {
	return &galleryGorm{db: db}
}

// Create will create the provided gallery and backfill data
// Like the ID, CreatedAt, and UpdatedAt fields.
func (gg *galleryGorm) Create(g *Gallery) error {
	return gg.db.Create(g).Error
}
