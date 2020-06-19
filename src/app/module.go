package app

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/gotemplate"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	"github.com/tessig/flamingo-product-rating/src/app/domain"
	"github.com/tessig/flamingo-product-rating/src/app/infrastructure"
	"github.com/tessig/flamingo-product-rating/src/app/interfaces"
	"github.com/tessig/flamingo-product-rating/src/app/interfaces/controller"
)

type (
	// Module basic struct
	Module struct{}

	routes struct {
		home   *controller.HomeController
		rating *controller.RatingController
	}
)

// Inject dependencies
func (r *routes) Inject(
	homeController *controller.HomeController,
	ratingController *controller.RatingController,
) {
	r.home = homeController
	r.rating = ratingController
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(gotemplate.Module),
	}
}

// Configure Rating module
func (m *Module) Configure(injector *dingo.Injector) {
	flamingo.BindTemplateFunc(injector, "random", new(interfaces.RandomIntFunc))
	flamingo.BindTemplateFunc(injector, "for", new(interfaces.ForFunc))
	flamingo.BindTemplateFunc(injector, "barType", new(interfaces.BarTypeFunc))

	injector.Bind(new(domain.RatingRepository)).To(new(infrastructure.RatingRepository))
	web.BindRoutes(injector, new(routes))
}

// Routes served by this module
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/", "home")
	registry.HandleGet("home", r.home.Home)

	registry.MustRoute("/products", "products")
	registry.HandleGet("products", r.home.ProductList)

	registry.MustRoute("/rating/$pid<[0-9]+>", "rating")
	registry.HandleGet("rating", r.rating.View)

	registry.MustRoute("/rating/new", "rating.product.form")
	registry.HandleGet("rating.product.form", r.rating.ProductForm)

	registry.MustRoute("/rating/$pid<[0-9]+>/new", "rating.new")
	registry.HandleGet("rating.new", r.rating.Form)

	registry.MustRoute("/rating/", "rating.post")
	registry.HandlePost("rating.post", r.rating.FormPost)

	registry.MustRoute("/rating/success", "rating.success")
	registry.HandleGet("rating.success", r.rating.Success)
}
