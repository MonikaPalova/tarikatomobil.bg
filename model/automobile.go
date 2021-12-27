package model

type Automobile struct {
	RegNumber string `json:"regNumber"`
	PhotoID   string `json:"photoId"`
	OwnerID   string `json:"ownedId"`
	Comment   string `json:"comment"`
}
