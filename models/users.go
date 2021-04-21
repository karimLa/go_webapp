package models

import (
	"errors"
	"regexp"
	"strings"

	"github.com/karimla/webapp/lib"
	"github.com/karimla/webapp/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrNotFound          = errors.New("models: resource not found")
	ErrIDInvalid         = errors.New("models: ID provided was invalid")
	ErrNotImplemented    = errors.New("models: not implemented")
	ErrEmailRequired     = errors.New("models: email address is required")
	ErrEmailInvalid      = errors.New("models: email address is not valid")
	ErrEmailTaken        = errors.New("models: email address is already taken")
	ErrPasswordInccorect = errors.New("models: incorrect password provided")
	ErrPasswordRequired  = errors.New("models: password is required")
	ErrPasswordTooShort  = errors.New("models: password must be at least 8 characters long")
	ErrRememberTooShort  = errors.New("models: remember token is too short")
	ErrRememberRequired  = errors.New("models: remember hash is required")
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
	// ErrNotFound, ErrPasswordInccorect, or another error if
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
// ErrNotFound, ErrPasswordInccorect, or another error if
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
			return nil, ErrPasswordInccorect
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
	hmac       lib.HMAC
	emailRegex *regexp.Regexp
}

func newUserValidator(udb UserDB) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       lib.NewHMAC(utils.GetSecret()),
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

// ByEmail will normalize the provided email and then call
// ByEmail on the subsequent UserDB layer.
func (uv *userValidator) ByEmail(email string) (*User, error) {
	u := User{Email: email}

	if err := runUserValFuncs(&u, uv.emailNormalize); err != nil {
		return nil, err
	}

	return uv.UserDB.ByEmail(u.Email)
}

// ByRemember will hash the remember token and then call
// ByRemember on the subsequent UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	u := User{Remember: token}

	if err := runUserValFuncs(&u, uv.rememberHmac); err != nil {
		return nil, err
	}

	return uv.UserDB.ByRemember(u.RememberHash)
}

// Create will hash user password and generate a remember token
// and then call Create on the subsequent UserDB layer.
func (uv *userValidator) Create(u *User) error {
	fns := []userValidatorFunc{
		uv.emailNormalize,
		uv.emailRequired,
		uv.emailIsValid,
		uv.emailIsAvail,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.passwordBcrypt,
		uv.passwordHashRequired,
		uv.rememberDefault,
		uv.rememberHmac,
		uv.rememberMinBytes,
		uv.rememberHashRequired,
	}
	if err := runUserValFuncs(u, fns...); err != nil {
		return err
	}

	return uv.UserDB.Create(u)
}

// Update generates a new remember token if necessary
// and then call Update on the subsequent UserDB layer.
func (uv *userValidator) Update(u *User) error {
	fns := []userValidatorFunc{
		uv.emailNormalize,
		uv.emailRequired,
		uv.emailIsValid,
		uv.emailIsAvail,
		uv.passwordMinLength,
		uv.passwordBcrypt,
		uv.passwordHashRequired,
		uv.rememberHmac,
		uv.rememberMinBytes,
		uv.rememberHashRequired,
	}
	if err := runUserValFuncs(u, fns...); err != nil {
		return err
	}

	return uv.UserDB.Update(u)
}

// Delete will call the subsequent UserDB layer if
// the provided id is valid. Otherwise it will return
// a ErrIDInvalid.
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
func (uv *userValidator) passwordBcrypt(u *User) error {
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

func (uv *userValidator) passwordRequired(u *User) error {
	if u.Password == "" {
		return ErrPasswordRequired
	}

	return nil
}

func (uv *userValidator) passwordHashRequired(u *User) error {
	if u.PasswordHash == "" {
		return ErrPasswordRequired
	}

	return nil
}

func (uv *userValidator) passwordMinLength(u *User) error {
	if u.Password == "" {
		return nil
	}

	if len(u.Password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}

func (uv *userValidator) rememberDefault(u *User) error {
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

func (uv *userValidator) rememberHmac(u *User) error {
	if u.Remember == "" {
		return nil
	}

	u.RememberHash = uv.hmac.Hash(u.Remember)
	return nil
}

func (uv *userValidator) rememberMinBytes(u *User) error {
	if u.Remember == "" {
		return nil
	}

	n, err := lib.NBytes(u.Remember)
	if err != nil {
		return err
	}

	if n < 32 {
		return ErrRememberTooShort
	}

	return nil
}

func (uv *userValidator) rememberHashRequired(u *User) error {
	if u.RememberHash == "" {
		return ErrRememberRequired
	}

	return nil
}

func (uv *userValidator) isGreaterThan(n uint) userValidatorFunc {
	return func(u *User) error {
		if u.ID <= n {
			return ErrIDInvalid
		}

		return nil
	}
}

func (uv *userValidator) emailNormalize(u *User) error {
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	return nil
}

func (uv *userValidator) emailRequired(u *User) error {
	if u.Email == "" {
		return ErrEmailRequired
	}

	return nil
}

func (uv *userValidator) emailIsValid(u *User) error {
	if !uv.emailRegex.MatchString(u.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvail(u *User) error {
	existing, err := uv.ByEmail(u.Email)

	if err == ErrNotFound {
		// Email address is not taken
		return nil
	}

	if err != nil {
		return err
	}

	// We found a user w/ this email address...
	// If the found user has the same ID as this user, it is
	// an update and this is the same user.
	if u.ID != existing.ID {
		return ErrEmailTaken
	}

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
