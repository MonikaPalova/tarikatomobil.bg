package model

import "time"

type Trip struct {
	ID               string    `json:"id"`
	From             string    `json:"from"`
	To               string    `json:"to"`
	When             time.Time `json:"when"`
	DriverName       string    `json:"driverName"`
	Price            float32   `json:"price"`
	MaxPassengers    int       `json:"maxPassengers"`
	AirConditioning  bool      `json:"airConditioning"`
	Smoking          bool      `json:"smoking"`
	Pets             bool      `json:"pets"`
	Comment          string    `json:"comment"`
}

type TripParticipation struct {
	TripID      string `json:"tripId"`
	DriverID    string `json:"driverId"`
	PassengerID string `json:"passengerId"`
}
