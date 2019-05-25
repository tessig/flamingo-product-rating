package controller

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/form/application"

	"github.com/tessig/flamingo-product-rating/src/app/domain"
	"github.com/tessig/flamingo-product-rating/src/app/interfaces/controller/form"
	"github.com/tessig/flamingo-product-rating/src/app/interfaces/controller/viewdata"
)

type (
	// RatingController provides the rating related actions
	RatingController struct {
		responder              *web.Responder
		ratingRepo             domain.RatingRepository
		productRepo            domain.ProductRepository
		formHandlerFactory     application.FormHandlerFactory
		ratingFormDataProvider *form.RatingFormDataProvider
	}

	// FormView is the data object for templates with the rating form
	FormView struct {
		Product   *domain.Product
		Structure []form.RatingFormField
		IsValid   bool
	}
)

// Inject dependencies
func (c *RatingController) Inject(
	r *web.Responder,
	ratingRepository domain.RatingRepository,
	productRepository domain.ProductRepository,
	formHandlerFactory application.FormHandlerFactory,
	ratingFormDataProvider *form.RatingFormDataProvider,
) *RatingController {
	c.responder = r
	c.ratingRepo = ratingRepository
	c.productRepo = productRepository
	c.formHandlerFactory = formHandlerFactory
	c.ratingFormDataProvider = ratingFormDataProvider

	return c
}

// View shows the ratings for a specific product
func (c *RatingController) View(_ context.Context, r *web.Request) web.Result {
	pidStr := r.Params["pid"]
	pid, _ := strconv.Atoi(pidStr)

	product, err := c.productRepo.Get(pid)
	if err != nil {
		return c.responder.ServerError(err)
	}

	average, err := c.ratingRepo.AverageByProductID(pid)
	if err != nil {
		return c.responder.ServerError(err)
	}

	breakdown, err := c.ratingRepo.BreakdownByProductID(pid)
	if err != nil {
		return c.responder.ServerError(err)
	}

	reviews, err := c.ratingRepo.ListByProductID(pid)
	if err != nil {
		return c.responder.ServerError(err)
	}

	return c.responder.Render(
		"rating/view",
		&viewdata.RatingData{
			Product:   product,
			Average:   average,
			Breakdown: breakdown,
			Reviews:   reviews,
		},
	)
}

// ProductForm shows the product selection form
func (c *RatingController) ProductForm(_ context.Context, r *web.Request) web.Result {
	if pid, err := r.Query1("pid"); err == nil {
		return c.responder.RouteRedirect("rating.new", map[string]string{"pid": pid})
	}

	products, err := c.productRepo.List()
	if err != nil {
		return c.responder.ServerError(err)
	}

	return c.responder.Render("rating/productform", products)
}

// Form shows the rating form
func (c *RatingController) Form(ctx context.Context, r *web.Request) web.Result {
	pidStr := r.Params["pid"]
	pid, _ := strconv.Atoi(pidStr)

	product, err := c.productRepo.Get(pid)
	if err != nil {
		return c.responder.ServerError(err)
	}

	formHandler := c.formHandlerFactory.CreateSimpleFormHandler()
	// HandleUnsubmittedForm provides default domain.Form instance without performing
	// http request body processing and form data validation
	f, err := formHandler.HandleUnsubmittedForm(ctx, r)
	fmt.Println(f)

	return c.responder.Render("rating/form", &FormView{Product: product, Structure: form.RatingFormStructure})
}

// FormPost receives the form data and saves the entity
func (c *RatingController) FormPost(ctx context.Context, r *web.Request) web.Result {

	formHandler := c.formHandlerFactory.CreateFormHandlerWithFormService(c.ratingFormDataProvider)
	// HandleSubmittedForm provides domain.Form instance after performing
	// http request body processing and form data validation
	f, err := formHandler.HandleSubmittedForm(ctx, r)
	// return on parse error (template need to handle error display)
	if err != nil {
		return c.responder.Render(
			"rating/form",
			&FormView{
				Structure: form.RatingFormStructure,
				IsValid:   f.ValidationInfo.IsValid(),
			},
		)
	}
	ratingFormData := f.Data.(form.RatingFormData)

	structure := make([]form.RatingFormField, len(form.RatingFormStructure))

	// take values into structure
	copy(structure, form.RatingFormStructure)
	for i, entry := range structure {
		structure[i].HasFeedback = true
		structure[i].Value, _ = r.Form1(entry.Name)
	}

	// take errors into structure
	for field, e := range f.ValidationInfo.GetErrorsForAllFields() {
		var index int
		for i, f := range structure {
			if f.Name == field {
				index = i
				break
			}
		}
		structure[index].HasError = true
		structure[index].Errors = make([]string, len(e))
		for i, msg := range e {
			structure[index].Errors[i] = msg.MessageKey
		}
	}

	pid, err := strconv.Atoi(ratingFormData.ProductID)
	if err != nil {
		return c.responder.ServerError(err)
	}
	product, err := c.productRepo.Get(pid)
	if err != nil {
		return c.responder.ServerError(err)
	}

	renderErrorResponse := c.responder.Render(
		"rating/form",
		&FormView{
			Product:   product,
			Structure: structure,
			IsValid:   f.ValidationInfo.IsValid(),
		},
	)

	if !f.IsValidAndSubmitted() {
		return renderErrorResponse
	}

	stars, err := strconv.Atoi(ratingFormData.Stars)
	if err != nil {
		return renderErrorResponse
	}
	rating := &domain.Rating{
		UserName:  ratingFormData.Name,
		ProductID: product.ID,
		CreatedAt: time.Now(),
		Title:     ratingFormData.Title,
		Text:      ratingFormData.Text,
		Stars:     stars,
	}

	err = c.ratingRepo.Save(rating)
	if err != nil {
		return renderErrorResponse
	}

	return c.responder.RouteRedirect("rating.success", nil)
}

// Success shows the rating form success page
func (c *RatingController) Success(_ context.Context, _ *web.Request) web.Result {
	return c.responder.Render("rating/success", nil)
}
