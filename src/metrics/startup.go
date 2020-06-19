package metrics

import (
	"context"
	"strconv"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/tessig/flamingo-product-rating/src/app/domain"
)

type (
	// StartUpMetrics subscribes to the AppStartup
	StartUpMetrics struct {
		ratingRepository domain.RatingRepository
		logger           flamingo.Logger
	}
)

// Inject dependencies
func (s *StartUpMetrics) Inject(r domain.RatingRepository, l flamingo.Logger) {
	s.ratingRepository = r
	s.logger = l
}

// Notify starts the rating metrics on the AppStartupEvent
func (s *StartUpMetrics) Notify(_ context.Context, event flamingo.Event) {
	switch event.(type) {
	case *flamingo.StartupEvent:
		s.logger.Info("Start rating metrics")
		s.ratingsMetrics()
	}
}

func (s *StartUpMetrics) ratingsMetrics() {
	ticker = time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			s.logger.Debug("collect rating metrics from DB")
			amount, errA := s.ratingRepository.Count(context.Background())
			amountsByProduct, errB := s.ratingRepository.Amounts(context.Background())
			if errA != nil || errB != nil {
				continue
			}

			ctx, _ := tag.New(context.Background(), tag.Upsert(keyProductID, "ALL"), tag.Upsert(opencensus.KeyArea, "-"))
			stats.Record(ctx, stat.M(amount))

			if amountsByProduct == nil {
				continue
			}

			for pid, amount := range *amountsByProduct {
				ctx, _ := tag.New(context.Background(), tag.Upsert(keyProductID, strconv.Itoa(pid)), tag.Upsert(opencensus.KeyArea, "-"))
				stats.Record(ctx, stat.M(int64(amount)))
			}
		}
	}()
}
