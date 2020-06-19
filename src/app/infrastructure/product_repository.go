package infrastructure

import (
	"context"
	"encoding/json"

	"go.opencensus.io/trace"

	"github.com/tessig/flamingo-product-rating/src/app/domain"
)

type (
	// ProductRepository connects to the product web service
	ProductRepository struct {
		source Source
	}

	// Source represents the product data source
	Source interface {
		Detail(ctx context.Context, id int) ([]byte, error)
		All(ctx context.Context) ([]byte, error)
	}
)

// Inject dependencies
func (p *ProductRepository) Inject(s Source) {
	p.source = s
}

// List returns all products
func (p *ProductRepository) List(ctx context.Context) ([]*domain.Product, error) {
	_, span := trace.StartSpan(ctx, "app/productrepository/List")
	defer span.End()

	data, err := p.source.All(ctx)
	if err != nil {
		return nil, err
	}

	var products []*domain.Product
	err = json.Unmarshal(data, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// Get returns a single product
func (p *ProductRepository) Get(ctx context.Context, id int) (*domain.Product, error) {
	ctx, span := trace.StartSpan(ctx, "app/productrepository/Get")
	defer span.End()

	data, err := p.source.Detail(ctx, id)
	if err != nil {
		return nil, err
	}

	product := &domain.Product{}
	err = json.Unmarshal(data, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}
