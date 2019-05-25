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

	injector.Bind((*domain.RatingRepository)(nil)).To(new(infrastructure.RatingRepository))
	injector.Bind((*domain.ProductRepository)(nil)).To(new(infrastructure.ProductRepository))
	web.BindRoutes(injector, new(routes))
}

// Routes served by this module
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.Route("/", "home")
	registry.HandleGet("home", r.home.Home)

	registry.Route("/products", "products")
	registry.HandleGet("products", r.home.ProductList)

	registry.Route("/rating/$pid<[0-9]+>", "rating")
	registry.HandleGet("rating", r.rating.View)

	registry.Route("/rating/new", "rating.product.form")
	registry.HandleGet("rating.product.form", r.rating.ProductForm)

	registry.Route("/rating/$pid<[0-9]+>/new", "rating.new")
	registry.HandleGet("rating.new", r.rating.Form)

	registry.Route("/rating/", "rating.post")
	registry.HandlePost("rating.post", r.rating.FormPost)

	registry.Route("/rating/success", "rating.success")
	registry.HandleGet("rating.success", r.rating.Success)
}
