package domain

import (
	"context"
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
		Count(ctx context.Context) (int64, error)
		CountByProductID(ctx context.Context, pid int) (int64, error)
		List(ctx context.Context) ([]*Rating, error)
		ListByProductID(ctx context.Context, pid int) ([]*Rating, error)
		Get(ctx context.Context, id int) (*Rating, error)
		Save(ctx context.Context, rating *Rating) error
		Delete(ctx context.Context, rating *Rating) error
		Amounts(ctx context.Context) (*RatingAmounts, error)
		Average(ctx context.Context) (*RatingAverage, error)
		AverageByProductID(ctx context.Context, pid int) (*RatingAverage, error)
		Breakdown(ctx context.Context) ([]*RatingBreakdown, error)
		BreakdownByProductID(ctx context.Context, pid int) ([]*RatingBreakdown, error)
	}
)
