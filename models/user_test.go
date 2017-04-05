package models_test

import (
	"testing"
)

func TestAuthenticateFailsWithIncorrectPassword(t *testing.T) {
	user := NewUser("go@appreciate.com", "password")
	if user.Authenticate("wrong") {
		t.Errorf("False")
	}
}
