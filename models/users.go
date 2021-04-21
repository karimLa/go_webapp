package models

import (
	"errors"
	"strings"

	"github.com/karimla/webapp/lib"
	"github.com/karimla/webapp/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrNotFound        = errors.New("models: resource not found")
	ErrInvalidID       = errors.New("models: ID provided was invalid")
	ErrInvalidPassword = errors.New("models: incorrect password provided")
	ErrNotImplemented  = errors.New("not implemented")
)

// User represents the user model stored in our database
// This is used for user accounts, storing both an email
// address and a password so users can log in and gain
// access to their content
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique;index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique;index"`
}

// UserDB is used to interact with the users database.
//
// For pretty much all single user queries:
// If the user is found, we will return a nil error
// If the user is not found, we will return ErrNotFound
// If there is another error, we will return an error with
// more information about what went wrong. This may not be
// an error generated by the models package.
//
// For single user queries, any error but ErrNotFound should
// probably result in a 500 error.
type UserDB interface {
	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Used to close a DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// UserService is a set of mthods used to manipulate and
// work with the user model
type UserService interface {
	// Authenticate will verify the provided email address and
	// password are correct. If they are correct, the user
	// corresponding to that email will be returned, Otherwise
	// it returns either:
	// ErrNotFound, ErrInvalidPassword, or another error if
	// something goes wrong.
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(db *gorm.DB) UserService {
	ug := newUserGorm(db)
	uv := newUserValidator(ug)

	return &userService{
		UserDB: uv,
	}
}

type userService struct {
	UserDB
}

// Authenticate will verify the provided email address and
// password are correct. If they are correct, the user
// corresponding to that email will be returned, Otherwise
// it returns either:
// ErrNotFound, ErrInvalidPassword, or another error if
// something goes wrong.
func (us *userService) Authenticate(email, password string) (*User, error) {
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

type userValidatorFunc func(*User) error

// runUserValFuncs runs the given fns passing user to each one.
// If it encountres an error, it returns it and breaks.
func runUserValFuncs(user *User, fns ...userValidatorFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

type userValidator struct {
	UserDB
	hmac lib.HMAC
}

func newUserValidator(ug UserDB) *userValidator {
	h := lib.NewHMAC(utils.GetSecret())
	return &userValidator{
		hmac:   h,
		UserDB: ug,
	}
}

// ByEmail will normalize the provided email and then call
// ByEmail on the subsequent UserDB layer.
func (uv *userValidator) ByEmail(email string) (*User, error) {
	u := User{Email: email}

	if err := runUserValFuncs(&u, uv.normalizeEmail); err != nil {
		return nil, err
	}

	return uv.UserDB.ByEmail(u.Email)
}

// ByRemember will hash the remember token and then call
// ByRemember on the subsequent UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	u := User{Remember: token}

	if err := runUserValFuncs(&u, uv.hmacRemember); err != nil {
		return nil, err
	}

	return uv.UserDB.ByRemember(u.RememberHash)
}

// Create will hash user password and generate a remember token
// and then call Create on the subsequent UserDB layer.
func (uv *userValidator) Create(u *User) error {
	fns := []userValidatorFunc{uv.normalizeEmail, uv.bcryptPassword, uv.defaultRemember, uv.hmacRemember}
	if err := runUserValFuncs(u, fns...); err != nil {
		return err
	}

	return uv.UserDB.Create(u)
}

// Update generates a new remember token if necessary
// and then call Update on the subsequent UserDB layer.
func (uv *userValidator) Update(u *User) error {
	if err := runUserValFuncs(u, uv.normalizeEmail, uv.bcryptPassword, uv.hmacRemember); err != nil {
		return err
	}

	return uv.UserDB.Update(u)
}

// Delete will call the subsequent UserDB layer if
// the provided id is valid. Otherwise it will return
// a ErrInvalidID.
func (uv *userValidator) Delete(id uint) error {
	u := User{Model: gorm.Model{ID: id}}

	if err := runUserValFuncs(&u, uv.isGreaterThan(0)); err != nil {
		return err
	}

	return uv.UserDB.Delete(id)
}

// bcryptPassword will hash a user's password with a
// predefined pepper and bcrypt if the Password field
// if not empty string
func (uv *userValidator) bcryptPassword(u *User) error {
	if u.Password == "" {
		return nil
	}

	pwBytes := []byte(u.Password + utils.GetPepper())
	hb, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hb)
	u.Password = ""

	return nil
}

func (uv *userValidator) defaultRemember(u *User) error {
	if u.Remember != "" {
		return nil
	}

	token, err := lib.RememberToken()
	if err != nil {
		return err
	}
	u.RememberHash = token

	return nil
}

func (uv *userValidator) hmacRemember(u *User) error {
	if u.Remember == "" {
		return nil
	}

	u.RememberHash = uv.hmac.Hash(u.Remember)
	return nil
}

func (uv *userValidator) isGreaterThan(n uint) userValidatorFunc {
	return func(u *User) error {
		if u.ID <= n {
			return ErrInvalidID
		}

		return nil
	}
}

func (uv *userValidator) normalizeEmail(u *User) error {
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	return nil
}

type userGorm struct {
	db *gorm.DB
}

func newUserGorm(db *gorm.DB) *userGorm {
	return &userGorm{db: db}
}

// ByID looks up a user with the provided ID.
// If the user is found, we will return a nil error
// If the user is found, we will return a nil error
// if the user is not found, we will return ErrNotFound
// If there is another error, we will return an error with
// more information about what went wrong. This may not be
// an error generated by the models package.
//
// As a general rule, any error but ErrNotFound should
// probably result in a 500 error
func (ug *userGorm) ByID(id uint) (*User, error) {
	var u User
	db := ug.db.Where("id = ?", id)
	err := first(db, &u)
	return &u, err
}

// ByEmail looks up a user with the given email address and
// returns that user.
// If the user is found, we will return a nil error
// if the user is not found, we will return ErrNotFound
// If there is another error, we will return an error with
// more information about what went wrong. This may not be
// an error generated by the models package.
//
// As a general rule, any error but ErrNotFound should
// probably result in a 500 error
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var u User
	db := ug.db.Where("email = ?", email)
	err := first(db, &u)
	return &u, err
}

// ByRemember looks up a user with the given remember token
// and returns that user. This method expects the remember
// token to be already hashed
//
// Errors are the same as ByEmail
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var u User
	db := ug.db.Where("remember_hash = ?", rememberHash)
	err := first(db, &u)
	return &u, err
}

// Create will create the provided user and backfill data
// Like the ID, CreatedAt, and UpdatedAt fields.
func (ug *userGorm) Create(u *User) error {
	return ug.db.Create(u).Error
}

// Update will update the provided user with all of the data
// in the provided user object.
func (ug *userGorm) Update(u *User) error {
	return ug.db.Save(u).Error
}

// Delete will delete the user with the provided ID
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Close returns a not implemented error
func (ug *userGorm) Close() error {
	return ErrNotImplemented
}

func (ug *userGorm) AutoMigrate() error {
	return ug.db.AutoMigrate(&User{})
}

func (ug *userGorm) DestructiveReset() error {
	utils.Must(ug.db.Migrator().DropTable(&User{}))
	return ug.db.AutoMigrate(&User{})
}

// first will query using the provided gorm.DB and it will
// get the first item returned amd place it into dst. If
// nothing is found in the query, it will return ErrNotFound
//
// NOTE: dst should be a pointer so that it populates
// the refrenced variable
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
