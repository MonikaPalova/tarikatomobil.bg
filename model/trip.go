package model

import (
	"errors"
	"time"
)

type Trip struct {
	ID              string    `json:"id"`
	From            string    `json:"from"`
	To              string    `json:"to"`
	When            time.Time `json:"when"`
	DriverName      string    `json:"driverName"`
	Price           float64   `json:"price"`
	MaxPassengers   int       `json:"maxPassengers"`
	AirConditioning bool      `json:"airConditioning"`
	Smoking         bool      `json:"smoking"`
	Pets            bool      `json:"pets"`
	Comment         string    `json:"comment"`
}

type TripFilter struct {
	From            string    `json:"from"`
	To              string    `json:"to"`
	Before          time.Time `json:"before"`
	After           time.Time `json:"after"`
	MaxPrice        float64   `json:"maxPrice"`
	AirConditioning bool      `json:"airConditioning"`
	Smoking         bool      `json:"smoking"`
	Pets            bool      `json:"pets"`
}

func DefaultTripFilter() TripFilter {
	return TripFilter{
		From:            "",
		To:              "",
		Before:          time.Now().Add(240 * time.Hour), // Default end-date is after 10 days
		After:           time.Now(),
		MaxPrice:        100,
		AirConditioning: false,
		Smoking:         false,
		Pets:            false,
	}
}

func (t Trip) Validate() error {
	if err := t.ValidateFrom(); err != nil {
		return err
	}
	if err := t.ValidateTo(); err != nil {
		return err
	}
	if err := t.ValidateWhen(); err != nil {
		return err
	}
	return t.ValidateMaxPassengers()
}

func (t Trip) ValidateFrom() error {
	if t.From == "" {
		return errors.New("from is a mandatory parameter of every trip")
	}
	return nil
}

func (t Trip) ValidateTo() error {
	if t.To == "" {
		return errors.New("to is a mandatory field of every trip")
	}
	return nil
}

func (t Trip) ValidateWhen() error {
	if t.When.Before(time.Now()) {
		return errors.New("when has to be a time in the future")
	}
	return nil
}

func (t Trip) ValidateMaxPassengers() error {
	if t.MaxPassengers <= 0 {
		return errors.New("maxPassengers needs to be a positive integer")
	}
	return nil
}

type TripParticipation struct {
	TripID      string `json:"tripId"`
	DriverID    string `json:"driverId"`
	PassengerID string `json:"passengerId"`
}
