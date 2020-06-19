package main

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/cache"
	"flamingo.me/flamingo/v3/core/gotemplate"
	coreHealthcheck "flamingo.me/flamingo/v3/core/healthcheck"
	"flamingo.me/flamingo/v3/core/locale"
	framework "flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/form"
	"github.com/tessig/flamingo-mysql/db"
	"github.com/tessig/flamingo-mysql/migration"

	"github.com/tessig/flamingo-product-rating/src/app"
	"github.com/tessig/flamingo-product-rating/src/healthcheck"
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
	registry.MustRoute("/static/*name", `flamingo.static.file(name,dir?="static")`)
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
			new(framework.SessionModule),
			new(coreHealthcheck.Module),
			new(healthcheck.Module),
			new(metrics.Module),
			new(app.Module),
			new(application),
		},
	)
}
