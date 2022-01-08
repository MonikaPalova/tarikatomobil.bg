package model

import "errors"

type Review struct {
	ID       string `json:"id"`
	FromUser string `json:"fromUser"`
	ForUser  string `json:"forUser"`
	Rating   int    `json:"rating"`
	Comment  string `json:"comment"`
}

func (r Review) Validate() error {
	if err := r.ValidateNotSameUser(); err != nil {
		return err
	}
	return r.ValidateRating()
}

func (r Review) ValidateNotSameUser() error {
	if r.FromUser == r.ForUser {
		return errors.New("you cannot create a review for yourself")
	}
	return nil
}

func (r Review) ValidateRating() error {
	if r.Rating < 1 || r.Rating > 5 {
		return errors.New("a review's rating must be from 1 to 5")
	}
	return nil
}
