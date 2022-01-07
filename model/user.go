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

// Used when updating the user's data (only these 4 fields can be updated)
type UserPatch struct {
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	PhotoID     string `json:"photoId"`
}

func (u *User) ValidateUserData() error {
	if err := validateUsername(u.Name); err != nil {
		return err
	}
	if err := validatePassword(u.Password); err != nil {
		return err
	}
	if err := validateEmail(u.Email); err != nil {
		return err
	}
	return validatePhoneNumber(u.PhoneNumber)
}

func (u *UserPatch) ValidateNonEmptyUserData() error {
	if u.Password != "" {
		if err := validatePassword(u.Password); err != nil {
			return err
		}
	}
	if u.Email != "" {
		if err := validateEmail(u.Email); err != nil {
			return err
		}
	}
	if u.PhoneNumber != "" {
		return validatePhoneNumber(u.PhoneNumber)
	}
	return nil
}

func validateUsername(username string) error {
	ok, _ := regexp.MatchString(usernameRegex, username)
	if !ok {
		return errors.New("username must begin with a letter and have only alphanumeric characters")
	}
	return nil
}

func validatePassword(password string) error {
	var (
		hasMinLen = false
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)
	if len(password) >= 7 {
		hasMinLen = true
	}
	for _, char := range password {
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

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("the provided email is not valid: %s", err.Error())
	}
	return err
}

func validatePhoneNumber(phoneNumber string) error {
	// For now just make sure it has at least 7 characters
	if len(phoneNumber) < 7 {
		return errors.New("the phone number is not valid")
	}
	return nil
}
