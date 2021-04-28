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

func (gv *galleryValidator) Create(g *Gallery) error {
	fns := []galleryValidatorFunc{gv.userIDRequired, gv.titleRequired}
	if err := runGalleryValFuncs(g, fns...); err != nil {
		return err
	}
	return gv.GalleryDB.Create(g)
}

func (gv *galleryValidator) userIDRequired(g *Gallery) error {
	if g.UserID <= 0 {
		return ErrUserIDRequired
	}

	return nil
}

func (gv *galleryValidator) titleRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}

	return nil
}

type galleryValidatorFunc func(*Gallery) error

// runGalleryValFuncs runs the given fns passing gallery to each one.
// If it encountres an error, it returns it and breaks.
func runGalleryValFuncs(g *Gallery, fns ...galleryValidatorFunc) error {
	for _, fn := range fns {
		if err := fn(g); err != nil {
			return err
		}
	}
	return nil
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
