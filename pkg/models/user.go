package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username   string
	Role       string
	HashedPass []byte
}

func NewUser(userName string, password string, role string) (*User, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User{
		Username:   userName,
		HashedPass: hashedPass,
		Role:       role,
	}
	return user, nil
}

func (u *User) VerifyPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword(u.HashedPass, []byte(password)); err != nil {
		return false
	}
	return true
}
