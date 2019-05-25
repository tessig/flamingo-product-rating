package viewdata

import (
	"github.com/tessig/flamingo-product-rating/src/app/domain"
)

type (
	// RatingData has all data for a rating display template
	RatingData struct {
		Product   *domain.Product
		Average   *domain.RatingAverage
		Breakdown []*domain.RatingBreakdown
		Reviews   []*domain.Rating
	}
)
