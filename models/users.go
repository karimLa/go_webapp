package models

import (
	"errors"

	"github.com/karimla/webapp/lib"
	"github.com/karimla/webapp/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrNotFound        = errors.New("models: resource not found")
	ErrInvalidID       = errors.New("models: ID provided was invalid")
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

type UserService struct {
	db   *gorm.DB
	hmac lib.HMAC
}

func NewUserService(db *gorm.DB) *UserService {
	h := lib.NewHMAC(utils.GetSecret())
	return &UserService{db: db, hmac: h}
}

func (us *UserService) ByID(id uint) (*User, error) {
	var u User
	db := us.db.Where("id = ?", id)
	err := first(db, &u)
	return &u, err
}

func (us *UserService) ByEmail(email string) (*User, error) {
	var u User
	db := us.db.Where("email = ?", email)
	err := first(db, &u)
	return &u, err
}

func (us *UserService) ByRemember(token string) (*User, error) {
	var u User
	hashedToken := us.hmac.Hash(token)
	db := us.db.Where("remember_hash = ?", hashedToken)
	err := first(db, &u)
	return &u, err
}

// Authenticate can be used to authenticate a user with the
// provided email address and password.
func (us *UserService) Authenticate(email, password string) (*User, error) {
	u, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	hpwBytes := []byte(u.PasswordHash)
	pwByes := []byte(password + utils.GetPepper())
	err = bcrypt.CompareHashAndPassword(hpwBytes, pwByes)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, ErrInvalidPassword
		}

		return nil, err
	}

	return u, nil
}

func (us *UserService) Create(u *User) error {
	pwBytes := []byte(u.Password + utils.GetPepper())
	hb, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hb)
	u.Password = ""

	if u.Remember == "" {
		token, err := lib.RememeberToken()
		if err != nil {
			return err
		}
		u.RememberHash = token
	} else {
		u.RememberHash = us.hmac.Hash(u.Remember)
	}

	return us.db.Create(u).Error
}

func (us *UserService) Update(u *User) error {
	if u.Remember != "" {
		u.RememberHash = us.hmac.Hash(u.Remember)
	}
	return us.db.Save(u).Error
}

func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

func (us *UserService) DestructiveReset() {
	us.db.Migrator().DropTable(&User{})
	us.db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique;index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique;index"`
}
