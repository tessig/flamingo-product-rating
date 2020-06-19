package form

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// RatingFormData represents the data of the rating form
	RatingFormData struct {
		ProductID string `form:"product_id" validate:"required" conform:"num"`
		Name      string `form:"name" validate:"required" conform:"name"`
		Stars     string `form:"stars" validate:"required" conform:"num"`
		Title     string `form:"title" validate:"required" conform:"trim"`
		Text      string `form:"text" validate:"required" conform:"trim"`
	}

	// RatingFormField represents a Field of the rating form
	RatingFormField struct {
		Type        string
		ID          string
		Name        string
		Label       string
		Value       string
		Icon        string
		HasFeedback bool
		HasError    bool
		Errors      []string
	}

	// RatingFormDataProvider is the flamingo form service implementation for the rating form
	RatingFormDataProvider struct{}
)

var (
	// RatingFormStructure defines the rating form
	RatingFormStructure = []RatingFormField{
		{
			Type:     "input",
			ID:       "name",
			Name:     "name",
			Label:    "User name",
			Value:    "",
			Icon:     "user",
			HasError: false,
			Errors:   nil,
		},
		{
			Type:     "stars",
			ID:       "stars",
			Name:     "stars",
			Label:    "Stars",
			Value:    "",
			Icon:     "",
			HasError: false,
			Errors:   nil,
		},
		{
			Type:     "input",
			ID:       "title",
			Name:     "title",
			Label:    "Title",
			Value:    "",
			Icon:     "header",
			HasError: false,
			Errors:   nil,
		},
		{
			Type:     "textarea",
			ID:       "text",
			Name:     "text",
			Label:    "Text",
			Value:    "",
			Icon:     "edit",
			HasError: false,
			Errors:   nil,
		},
	}
)

// GetFormData returns the default form values
func (p *RatingFormDataProvider) GetFormData(context.Context, *web.Request) (interface{}, error) {
	return RatingFormData{}, nil
}
