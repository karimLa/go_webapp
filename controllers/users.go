package controllers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/karimla/webapp/lib"
	"github.com/karimla/webapp/models"
	"github.com/karimla/webapp/views"
)

// New Users is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewUsers(wg *sync.WaitGroup, us models.UserService) *Users {
	return &Users{
		SignupView: views.NewView(wg, "bootstrap", "users/new"),
		LoginView:  views.NewView(wg, "bootstrap", "users/login"),
		us:         us,
		wg:         wg,
	}
}

type Users struct {
	SignupView *views.View
	LoginView  *views.View
	us         models.UserService
	wg         *sync.WaitGroup
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
	u.wg.Add(1)
	defer u.wg.Done()

	var form SignupForm
	parseForm(r, &form)

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := u.signIn(w, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
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
	u.wg.Add(1)
	defer u.wg.Done()

	var form LoginForm
	parseForm(r, &form)

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		if err == models.ErrNotFound {
			fmt.Fprintln(w, "Invalid email address.")
			return
		}

		if err == models.ErrPasswordInccorect {
			fmt.Fprintln(w, "Invalid password.")
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = u.signIn(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
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
