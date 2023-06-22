// delete_user_test.go

package main

import (
	"testing"
)

func TestDeleteUser(t *testing.T) {
	userService := UserService{}
	user := User{
		Username: "john",
		Password: "secret",
	}
	userService.Users = append(userService.Users, user)

	err := userService.DeleteUser("john")
	if err != nil {
		t.Errorf("Error deleting user: %s", err.Error())
	}

	if len(userService.Users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(userService.Users))
	}
}

