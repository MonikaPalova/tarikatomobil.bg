package db

import (
	"database/sql"
	"fmt"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
)

type ReviewsDBHandler struct {
	conn *sql.DB
}

func (rdb *ReviewsDBHandler) GetReviewsForUser(username string) ([]model.Review, *DBError) {
	return nil, nil
}

func (rdb *ReviewsDBHandler) GetReviewsFromUser(username string) ([]model.Review, *DBError) {
	return nil, nil
}

func (rdb *ReviewsDBHandler) CreateReview(review model.Review) *DBError {
	insertQuery := `INSERT INTO reviews (id, from_user, for_user, rating, comment) VALUES (?, ?, ?, ?, ?)`
	stmt, err := rdb.conn.Prepare(insertQuery)
	if err != nil {
		return NewDBError(err, ErrInternal)
	}

	if _, err = stmt.Exec(review.ID, review.FromUser, review.ForUser, review.Rating, review.Comment); err != nil {
		if isForeignKeyError(err) {
			return NewDBError(err, ErrNotFound, fmt.Sprintf("user %s does not exist", review.ForUser))
			// the HTTP handler already validated that review.FromUser is a valid user
		}
		return NewDBError(err, ErrInternal)
	}
	return nil
}

func (rdb *ReviewsDBHandler) DeleteReview(reviewID string) *DBError {
	return nil
}
