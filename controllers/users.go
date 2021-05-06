package controllers

import (
	"log"
	"net/http"

	"webapp/lib"
	"webapp/models"
	"webapp/views"

	"github.com/gorilla/mux"
)

// New Users is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewUsers(us models.UserService, r *mux.Router, l *log.Logger) *Users {
	return &Users{
		SignupView: views.NewView("bootstrap", "users/new"),
		LoginView:  views.NewView("bootstrap", "users/login"),
		us:         us,
		r:          r,
		l:          l,
	}
}

type Users struct {
	SignupView *views.View
	LoginView  *views.View
	us         models.UserService
	r          *mux.Router
	l          *log.Logger
}

type SignupForm struct {
	Name     string `schema:"name,required"`
	Email    string `schema:"email,required"`
	Password string `schema:"password,required"`
}

// Signup is used to process the signup form whem a user
// submits it. This is used to create a new user account.
//
// POST /signup
func (u *Users) Signup(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.SignupView.Render(w, r, vd)
		return
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.SignupView.Render(w, r, vd)
		return
	}

	if err := u.signIn(w, &user); err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	path := Reverse(GalleriesIndexURL, "/", u.r)
	http.Redirect(w, r, path, http.StatusFound)
}

type LoginForm struct {
	Email    string `schema:"email,required"`
	Password string `schema:"password,required"`
}

// Login is used to verify the provided email address and
// password and then log the user in if they are correct.
//
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}
	var form LoginForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		if err == models.ErrNotFound {
			vd.AlertError("Invalid email address")
			u.LoginView.Render(w, r, vd)
			return
		}

		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	if err = u.signIn(w, user); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	path := Reverse(GalleriesIndexURL, "/", u.r)
	http.Redirect(w, r, path, http.StatusFound)
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := lib.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	c := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &c)

	return nil
}
