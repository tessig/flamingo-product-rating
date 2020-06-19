package products

import (
	"flamingo.me/dingo"

	"github.com/tessig/flamingo-product-rating/src/app"
	"github.com/tessig/flamingo-product-rating/src/app/domain"
	products "github.com/tessig/flamingo-product-rating/src/products/infrastructure"
)

type (
	// Module basic struct
	Module struct{}
)

// Configure product module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(domain.ProductRepository)).To(new(products.Client))
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(app.Module),
	}
}

// CueConfig for the module
func (m *Module) CueConfig() string {
	return `
productservice: {
	baseurl: string | *"http://localhost:8080/"
	endpoints: {
		list: string | *"products"
		detail: string | *"products/id/:pid"
	}
}
`
}
