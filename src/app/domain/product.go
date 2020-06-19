package domain

import (
	"context"
)

type (
	// Product represents a simple product information model
	Product struct {
		ID   int
		Name string
	}

	// ProductRepository can retrieve product information
	ProductRepository interface {
		List(ctx context.Context) ([]*Product, error)
		Get(ctx context.Context, id int) (*Product, error)
	}
)
