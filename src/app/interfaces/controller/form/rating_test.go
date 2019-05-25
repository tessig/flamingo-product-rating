package form_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/form/application"
	"flamingo.me/form/domain/formdata"
	"github.com/go-test/deep"

	"github.com/tessig/flamingo-product-rating/src/app/interfaces/controller/form"
)

func TestRatingFormService_ParseFormData(t *testing.T) {
	type args struct {
		formValues url.Values
	}
	tests := []struct {
		name string
		args args
		want form.RatingFormData
	}{
		{
			name: "valid case",
			args: args{
				formValues: url.Values{
					"product_id": []string{"7"},
					"name":       []string{"Chris"},
					"stars":      []string{"5"},
					"title":      []string{"A title"},
					"text":       []string{"The text"},
				},
			},
			want: form.RatingFormData{
				ProductID: "7",
				Name:      "Chris",
				Stars:     "5",
				Title:     "A title",
				Text:      "The text",
			},
		},
		{
			name: "work for conform",
			args: args{
				formValues: url.Values{
					"product_id": []string{" tra7tra "},
					"name":       []string{"•?((¯°·._.•Chris•._.·°¯))؟•"},
					"stars":      []string{"  give 5 stars please!! ★★★★★"},
					"title":      []string{"   A title which must be trimmed  !    "},
					"text":       []string{"   The text  to be trimmed   "},
				},
			},
			want: form.RatingFormData{
				ProductID: "7",
				Name:      "Chris",
				Stars:     "5",
				Title:     "A title which must be trimmed  !",
				Text:      "The text  to be trimmed",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formHandlerFactory := new(application.FormHandlerFactoryImpl)
			validatorProvider := new(application.ValidatorProviderImpl)
			validatorProvider.Inject(nil, nil)
			formHandlerFactory.Inject(
				nil,
				nil,
				nil,
				nil,
				nil,
				new(formdata.DefaultFormDataProviderImpl),
				new(formdata.DefaultFormDataDecoderImpl),
				new(formdata.DefaultFormDataValidatorImpl),
				validatorProvider,
				flamingo.NullLogger{},
			)
			formHandler := formHandlerFactory.CreateFormHandlerWithFormService(new(form.RatingFormDataProvider))
			request := new(http.Request)
			request.PostForm = tt.args.formValues
			f, err := formHandler.HandleSubmittedForm(
				context.Background(),
				web.CreateRequest(request, web.EmptySession()),
			)
			if err != nil {
				t.Fatal("error on form handling:", err)
			}
			got := f.Data.(form.RatingFormData)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}
