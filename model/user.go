package model

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"unicode"
)

const (
	usernameRegex = `^[a-zA-Z]+[a-zA-Z0-9]*$`
)

type User struct {
	Name           string `json:"name"`
	Password       string `json:"password"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	PhotoID        string `json:"photoId"`
	TimesPassenger int    `json:"timesPassenger"`
	TimesDriver    int    `json:"timesDriver"`
}

// Used when GET-ting the user (so that the password is not returned by the GET API)
type UserWithoutPass struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	PhotoID        string `json:"photoId"`
	TimesPassenger int    `json:"timesPassenger"`
	TimesDriver    int    `json:"timesDriver"`
}

func (u *User) ValidateUserData() error {
	if err := u.ValidateUsername(); err != nil {
		return err
	}
	if err := u.ValidatePassword(); err != nil {
		return err
	}
	return u.ValidateEmail()
}

func (u *User) ValidateUsername() error {
	ok, _ := regexp.MatchString(usernameRegex, u.Name)
	if !ok {
		return errors.New("username must begin with a letter and have only alphanumeric characters")
	}
	return nil
}

func (u *User) ValidatePassword() error {
	var (
		hasMinLen = false
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)
	if len(u.Password) >= 7 {
		hasMinLen = true
	}
	for _, char := range u.Password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}
	if !(hasMinLen && hasUpper && hasLower && hasNumber) {
		return errors.New(`password must be at least 7 characters long and contain at least one of the following: 
			capital letter, lower letter, number`)
	}
	return nil
}

func (u *User) ValidateEmail() error {
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return fmt.Errorf("the provided email is not valid: %s", err.Error())
	}
	return err
}
