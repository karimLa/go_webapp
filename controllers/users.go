package controllers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/karimla/webapp/models"
	"github.com/karimla/webapp/utils"
	"github.com/karimla/webapp/views"
)

// New Users is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewUsers(wg *sync.WaitGroup, us *models.UserService) *Users {
	return &Users{
		newView: views.NewView(wg, "bootstrap", "users/new"),
		us:      us,
		wg:      wg,
	}
}

type Users struct {
	newView *views.View
	us      *models.UserService
	wg      *sync.WaitGroup
}

// New is used to render the form where a user can
// create a new user account
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	utils.Must(u.newView.Render(w, nil))
}

type SignupForm struct {
	Name     string `schema:"name,required"`
	Email    string `schema:"email,required"`
	Password string `schema:"password,required"`
}

// Create is used to process the signup form whem a user
// submits it. This is used to create a new user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	u.wg.Add(1)
	defer u.wg.Done()

	var form SignupForm
	utils.Must(parseForm(r, &form))

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, user)
}
