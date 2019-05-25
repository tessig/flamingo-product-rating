package main

import (
	"net/http"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/cache"
	"flamingo.me/flamingo/v3/core/gotemplate"
	"flamingo.me/flamingo/v3/core/locale"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/form"

	"github.com/tessig/flamingo-mysql/db"
	"github.com/tessig/flamingo-mysql/migration"
	"github.com/tessig/flamingo-product-rating/src/app"
	"github.com/tessig/flamingo-product-rating/src/metrics"
	"github.com/tessig/flamingo-product-rating/src/products"
)

type (
	application   struct{}
	defaultRoutes struct{}
)

func (a *application) Configure(injector *dingo.Injector) {
	injector.Bind((*cache.Backend)(nil)).ToInstance(cache.NewInMemoryCache())
	web.BindRoutes(injector, &defaultRoutes{})
}

// Routes
func (a *defaultRoutes) Routes(registry *web.RouterRegistry) {
	registry.Route("/static/*n", "_static")
	registry.HandleGet(
		"_static",
		web.WrapHTTPHandler(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))),
	)
}

func main() {
	flamingo.App(
		[]dingo.Module{
			new(locale.Module),
			new(form.Module),
			new(gotemplate.Module),
			new(db.Module),
			new(migration.Module),
			new(products.Module),
			new(opencensus.Module),
			new(metrics.Module),
			new(app.Module),
			new(application),
		},
	)
}
