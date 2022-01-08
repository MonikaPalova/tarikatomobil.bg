package model

type Automobile struct {
	RegNumber string `json:"regNumber"`
	PhotoID   string `json:"photoId"`
	OwnerName string `json:"ownerName"`
	Comment   string `json:"comment"`
}

type AutomobilePatch struct {
	PhotoID   string `json:"photoId"`
	Comment   string `json:"comment"`
}