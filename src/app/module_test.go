package app_test

import (
	"testing"

	"flamingo.me/dingo"

	"github.com/tessig/flamingo-product-rating/src/app"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(app.Module)); err != nil {
		t.Error(err)
	}
}
