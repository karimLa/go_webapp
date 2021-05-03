package models_test

import (
	"testing"

	"github.com/nicholasjackson/env"
	"github.com/soramon0/webapp/models"
	"github.com/soramon0/webapp/utils"
)

func testingUserService() models.UserService {
	utils.Must(env.Parse())
	s := models.NewServices()
	s.DestructiveReset()
	return s.User
}

func TestCreateUser(t *testing.T) {
	us := testingUserService()

	user := models.User{
		Name:     "Sam Lee",
		Email:    "sam@test.com",
		Password: "password",
	}
	err := us.Create(&user)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Errorf("Expected ID > 0. Recieved %d", user.ID)
	}
}
