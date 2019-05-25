package domain

import (
	"time"
)

type (
	// Rating represents a single product rating
	Rating struct {
		ID        int       `db:"id"`
		UserName  string    `db:"name"`
		ProductID int       `db:"product_id"`
		CreatedAt time.Time `db:"created_at"`
		Title     string    `db:"title"`
		Text      string    `db:"text"`
		Stars     int       `db:"stars"`
	}

	// RatingAverage represents the mean of a set of ratings
	RatingAverage struct {
		Value  float64
		Amount int
		Max    int
		Stars  int
	}

	// RatingBreakdown represents the amount of n-stars ratings for a set of ratings
	RatingBreakdown struct {
		Stars      int
		Amount     int
		Total      int
		Percentage float64
	}

	// RatingAmounts represents the number of ratings by product ID
	RatingAmounts map[int]int

	// RatingRepository provides all functions for Rating storage
	RatingRepository interface {
		Count() (int64, error)
		CountByProductID(pid int) (int64, error)
		List() ([]*Rating, error)
		ListByProductID(pid int) ([]*Rating, error)
		Get(id int) (*Rating, error)
		Save(rating *Rating) error
		Delete(rating *Rating) error
		Amounts() (*RatingAmounts, error)
		Average() (*RatingAverage, error)
		AverageByProductID(pid int) (*RatingAverage, error)
		Breakdown() ([]*RatingBreakdown, error)
		BreakdownByProductID(pid int) ([]*RatingBreakdown, error)
	}
)
