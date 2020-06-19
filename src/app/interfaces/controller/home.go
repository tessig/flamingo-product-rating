package controller

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
	"go.opencensus.io/trace"

	"github.com/tessig/flamingo-product-rating/src/app/domain"
	"github.com/tessig/flamingo-product-rating/src/app/interfaces/controller/viewdata"
)

type (
	// HomeController provides the start page
	HomeController struct {
		responder   *web.Responder
		ratingRepo  domain.RatingRepository
		productRepo domain.ProductRepository
	}
)

// Inject the dependencies
func (c *HomeController) Inject(
	r *web.Responder,
	ratingRepository domain.RatingRepository,
	productRepository domain.ProductRepository,
) {
	c.responder = r
	c.ratingRepo = ratingRepository
	c.productRepo = productRepository
}

// Home provides the general overview
func (c *HomeController) Home(ctx context.Context, _ *web.Request) web.Result {
	ctx, span := trace.StartSpan(ctx, "app/controller/Home")
	defer span.End()

	average, err := c.ratingRepo.Average(ctx)
	if err != nil {
		return c.responder.ServerError(err)
	}

	breakdown, err := c.ratingRepo.Breakdown(ctx)
	if err != nil {
		return c.responder.ServerError(err)
	}

	reviews, err := c.ratingRepo.List(ctx)
	if err != nil {
		return c.responder.ServerError(err)
	}

	return c.responder.Render("index", &viewdata.RatingData{
		Average:   average,
		Breakdown: breakdown,
		Reviews:   reviews,
	})
}

// ProductList shows a list of all products with links to their review overview pages
func (c *HomeController) ProductList(ctx context.Context, _ *web.Request) web.Result {
	ctx, span := trace.StartSpan(ctx, "app/controller/ProductList")
	defer span.End()

	products, err := c.productRepo.List(ctx)
	if err != nil {
		return c.responder.ServerError(err)
	}

	amounts, err := c.ratingRepo.Amounts(ctx)
	if err != nil {
		return c.responder.ServerError(err)
	}

	return c.responder.Render(
		"products", struct {
			Products []*domain.Product
			Amounts  *domain.RatingAmounts
		}{
			Products: products,
			Amounts:  amounts,
		},
	)
}
