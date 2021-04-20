package models_test

import (
	"testing"

	"github.com/karimla/webapp/lib"
	"github.com/karimla/webapp/models"
)

func testingUserService() *models.UserService {
	db := lib.InitDB()

	us := models.NewUserService(db)

	us.DestructiveReset()

	return us
}

func TestCreateUser(t *testing.T) {
	us := testingUserService()

	user := models.User{
		Name:     "Sam Lee",
		Email:    "sam@test.com",
		Password: "",
	}
	err := us.Create(&user)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Errorf("Expected ID > 0. Recieved %d", user.ID)
	}
}
