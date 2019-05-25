package domain

type (
	// Product represents a simple product information model
	Product struct {
		ID   int
		Name string
	}

	// ProductRepository can retrieve product information
	ProductRepository interface {
		List() ([]*Product, error)
		Get(id int) (*Product, error)
	}
)
