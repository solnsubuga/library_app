package models

import (
	"golang.org/x/crypto/bcrypt"
)

//Person ...
type Person struct {
	Username string
	Password string
}

//NewUser ...
func NewUser(un, pw string) *Person {
	pass, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return &Person{
		Username: un,
		Password: string(pass),
	}
}

//Authenticate ..
func (u *User) Authenticate(pw string) bool {
	return false //bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw)) == nil
}
