package products

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"

	"github.com/tessig/flamingo-product-rating/src/app/infrastructure"
	products "github.com/tessig/flamingo-product-rating/src/products/infrastructure"
)

type (
	// Module basic struct
	Module struct{}
)

// Configure product module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*infrastructure.Source)(nil)).To(new(products.Client))
}

// DefaultConfig for product module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"productservice.baseurl":          "http://localhost:8080/",
		"productservice.endpoints.list":   "products",
		"productservice.endpoints.detail": "products/id/:pid",
	}
}
