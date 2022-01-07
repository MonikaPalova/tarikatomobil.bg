package model

type Automobile struct {
	RegNumber string `json:"regNumber"`
	PhotoID   string `json:"photoId"`
	OwnerName string `json:"ownerName"`
	Comment   string `json:"comment"`
}
