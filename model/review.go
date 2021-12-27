package model

type Review struct {
	ID       string `json:"id"`
	FromUser string `json:"fromUser"`
	ForUser  string `json:"forUser"`
	Rating   int    `json:"rating"`
	Comment  string `json:"comment"`
}
