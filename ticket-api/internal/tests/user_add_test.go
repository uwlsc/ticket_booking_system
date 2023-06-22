// add_user_test.go

package tests

import (
	"testing"
)

func TestAddUser(t *testing.T) {
	userService := UserService{}
	user := User{
		Username: "john",
		Password: "secret",
	}

	userService.AddUser(user)

	if len(userService.Users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(userService.Users))
	}

	if userService.Users[0].Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, userService.Users[0].Username)
	}
}

