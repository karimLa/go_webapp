package models

import "gorm.io/gorm"

type Gallery struct {
	gorm.Model
	UserID uint     `gorm:"not null;index"`
	Title  string   `gorm:"not null"`
	Images []string `gorm:"-"`
}

func (g *Gallery) ImageSplitN(n int) [][]string {
	ret := make([][]string, n)
	for i := 0; i < n; i++ {
		ret[i] = make([]string, 0)
	}

	for i, img := range g.Images {
		bucket := i % n
		ret[bucket] = append(ret[bucket], img)
	}
	return ret
}

type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	ByUserID(id uint) ([]Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(id uint) error
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

func (gv *galleryValidator) Update(g *Gallery) error {
	fns := []galleryValidatorFunc{gv.userIDRequired, gv.titleRequired}
	if err := runGalleryValFuncs(g, fns...); err != nil {
		return err
	}
	return gv.GalleryDB.Update(g)
}

func (gv *galleryValidator) Delete(id uint) error {
	g := Gallery{Model: gorm.Model{ID: id}}

	if err := runGalleryValFuncs(&g, gv.isGreaterThan(0)); err != nil {
		return err
	}

	return gv.GalleryDB.Delete(id)
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

func (gv *galleryValidator) isGreaterThan(n uint) galleryValidatorFunc {
	return func(g *Gallery) error {
		if g.ID <= n {
			return ErrIDInvalid
		}

		return nil
	}
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

func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var g Gallery
	db := gg.db.Where("id = ?", id)
	err := first(db, &g)
	return &g, err
}

func (gg *galleryGorm) ByUserID(userID uint) ([]Gallery, error) {
	var galleries []Gallery
	gg.db.Where("user_id = ?", userID).Find(&galleries)
	return galleries, nil
}

// Create will create the provided gallery and backfill data
// Like the ID, CreatedAt, and UpdatedAt fields.
func (gg *galleryGorm) Create(g *Gallery) error {
	return gg.db.Create(g).Error
}

func (gg *galleryGorm) Update(g *Gallery) error {
	return gg.db.Save(g).Error
}

func (gg *galleryGorm) Delete(id uint) error {
	g := Gallery{Model: gorm.Model{ID: id}}
	return gg.db.Delete(&g).Error
}
