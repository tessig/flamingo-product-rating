package products_test

import (
	"testing"

	"flamingo.me/dingo"

	"github.com/tessig/flamingo-product-rating/src/products"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(products.Module)); err != nil {
		t.Error(err)
	}
}
