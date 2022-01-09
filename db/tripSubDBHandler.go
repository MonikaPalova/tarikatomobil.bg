package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
)

type TripSubscriptionDBHandler struct {
	conn *sql.DB
}

func (tsdb *TripSubscriptionDBHandler) GetPassengersForTrip(tripID string) ([]model.UserWithoutPass, *DBError) {
	rows, err := tsdb.conn.Query(`SELECT u.name, u.email, u.phone_number, u.photo_id, u.times_passenger,
		u.times_driver FROM trip_participations t join users u on u.name = t.passenger_name WHERE t.trip_id = ?`,
		tripID)
	if err != nil {
		return nil, NewDBError(err, ErrInternal)
	}

	users := []model.UserWithoutPass{}
	for rows.Next() {
		var user model.UserWithoutPass
		if err := rows.Scan(&user.Name, &user.Email, &user.PhoneNumber, &user.PhotoID,
			&user.TimesPassenger, &user.TimesDriver); err != nil {
			return nil, NewDBError(err, ErrInternal)
		}
		users = append(users, user)
	}
	return users, nil
}

func (tsdb *TripSubscriptionDBHandler) GetFreeSpotsCount(tripID string) (int, *DBError) {
	subscribedUsers, err := tsdb.GetPassengersForTrip(tripID)
	if err != nil {
		return 0, err
	}
	row := tsdb.conn.QueryRow("SELECT max_passengers FROM trips WHERE id = ?", tripID)
	var maxTripPassengers int
	if err := row.Scan(&maxTripPassengers); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, NewDBError(err, ErrNotFound, fmt.Sprintf("trip with id %s does not exist", tripID))
		}
		return 0, NewDBError(err, ErrInternal)
	}

	return maxTripPassengers - len(subscribedUsers), nil
}

func (tsdb *TripSubscriptionDBHandler) GetTripOwner(tripID string) (string, *DBError) {
	var username string
	row := tsdb.conn.QueryRow("SELECT driver_name FROM trips WHERE id = ? ", tripID)
	if err := row.Scan(&username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return username, NewDBError(err, ErrNotFound, fmt.Sprintf("trip with id %s does not exist", tripID))
		}
		return username, NewDBError(err, ErrInternal)
	}
	return username, nil
}

func (tsdb *TripSubscriptionDBHandler) SubscribePassenger(username, tripID string) *DBError {
	if _, err := tsdb.conn.Exec("INSERT INTO trip_participations (trip_id, passenger_name) VALUES(?,?)",
		tripID, username); err != nil {
		if isDuplicateEntryError(err) {
			return NewDBError(err, ErrConflict, fmt.Sprintf("user %s is already subscribed to trip %s", username, tripID))
		}
		return NewDBError(err, ErrInternal)
	}
	return nil
}

func (tsdb *TripSubscriptionDBHandler) UnsubscribePassenger(username, tripID string) *DBError {
	result, err := tsdb.conn.Exec("DELETE FROM trip_participations WHERE trip_id = ? AND passenger_name = ?",
		tripID, username)
	if err != nil {
		return NewDBError(err, ErrInternal)
	}
	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return NewDBError(err, ErrInternal)
	}
	if rowsAffected == 0 {
		return NewDBError(err, ErrNotFound, fmt.Sprintf("user %s is not subscribed to trip %s", username, tripID))
	}
	return nil
}
