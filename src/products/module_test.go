package products_test

import (
	"testing"

	"flamingo.me/flamingo/v3/framework/config"

	"github.com/tessig/flamingo-product-rating/src/products"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(nil, new(products.Module)); err != nil {
		t.Error(err)
	}
}
