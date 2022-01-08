package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
)

type TripsDBHandler struct {
	conn *sql.DB
}

func (tdb *TripsDBHandler) CreateTrip(trip model.Trip) *DBError {
	insertQuery := `INSERT INTO trips (id, location_from, location_to, departure_time, driver_name, price, max_passengers, air_conditioning, smoking, pets, comment)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := tdb.conn.Prepare(insertQuery)
	if err != nil {
		return NewDBError(err, ErrInternal)
	}

	if _, err = stmt.Exec(trip.ID, trip.From, trip.To, trip.When, trip.DriverName, trip.Price, trip.MaxPassengers,
		trip.AirConditioning, trip.Smoking, trip.Pets, trip.Comment); err != nil {
		return NewDBError(err, ErrInternal)
	}
	return nil
}

func (tdb *TripsDBHandler) DeleteTrip(tripID, caller string) *DBError {
	result, err := tdb.conn.Exec("DELETE FROM trips WHERE id = ? AND driver_name = ?", tripID, caller)
	if err != nil {
		if isForeignKeyError(err) {
			return NewDBError(err, ErrConflict, "cannot delete a trip that has subscribed passengers")
		}
		return NewDBError(err, ErrInternal)
	}
	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return NewDBError(err, ErrInternal)
	}
	if rowsAffected == 0 {
		return NewDBError(err, ErrNotFound, fmt.Sprintf("a trip with ID %s owned by %s does not exist",
			tripID, caller))
	}
	return nil
}

func (tdb *TripsDBHandler) GetTrips(filter model.TripFilter) ([]model.Trip, *DBError) {
	query := `SELECT id, location_from, location_to, departure_time, driver_name, price,
		max_passengers, air_conditioning, smoking, pets, comment FROM trips
		WHERE
		location_from = ? AND
		location_to = ? AND 
		TIMESTAMP(departure_time) > ? AND TIMESTAMP(departure_time) < ? AND
		price < ? AND
		air_conditioning = ? AND
		smoking = ? AND
		pets = ?`
	rows, err := tdb.conn.Query(query, filter.From, filter.To, filter.After, filter.Before, filter.MaxPrice,
		filter.AirConditioning, filter.Smoking, filter.Pets)
	if err != nil {
		return nil, NewDBError(err, ErrInternal)
	}

	trips := []model.Trip{}
	for rows.Next() {
		var trip model.Trip
		if err := rows.Scan(&trip.ID, &trip.From, &trip.To, &trip.When, &trip.DriverName, &trip.Price,
			&trip.MaxPassengers, &trip.AirConditioning, &trip.Smoking, &trip.Pets, &trip.Comment); err != nil {
			return nil, NewDBError(err, ErrInternal)
		}
		trips = append(trips, trip)
	}
	return trips, nil
}

func (tdb *TripsDBHandler) GetTrip(tripID string) (model.Trip, *DBError) {
	row := tdb.conn.QueryRow(`SELECT id, location_from, location_to, departure_time, driver_name, price,
		max_passengers, air_conditioning, smoking, pets, comment FROM trips WHERE id = ?`, tripID)
	var trip model.Trip
	if err := row.Scan(&trip.ID, &trip.From, &trip.To, &trip.When, &trip.DriverName, &trip.Price,
		&trip.MaxPassengers, &trip.AirConditioning, &trip.Smoking, &trip.Pets, &trip.Comment); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return trip, NewDBError(err, ErrNotFound, fmt.Sprintf("trip with id %s does not exist", tripID))
		}
		return trip, NewDBError(err, ErrInternal)
	}
	return trip, nil

}
