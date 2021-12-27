package model

type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	PhotoID        string `json:"photoId"`
	TimesPassenger int    `json:"timesPassenger"`
	TimesDriver    int    `json:"timesDriver"`
}
