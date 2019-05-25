package infrastructure

import (
	"encoding/json"

	"github.com/tessig/flamingo-product-rating/src/app/domain"
)

type (
	// ProductRepository connects to the product web service
	ProductRepository struct {
		source Source
	}

	// Source represents the product data source
	Source interface {
		Detail(int) ([]byte, error)
		All() ([]byte, error)
	}
)

// Inject dependencies
func (p *ProductRepository) Inject(s Source) {
	p.source = s
}

// List returns all products
func (p *ProductRepository) List() ([]*domain.Product, error) {
	data, err := p.source.All()
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
func (p *ProductRepository) Get(id int) (*domain.Product, error) {
	data, err := p.source.Detail(id)
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
