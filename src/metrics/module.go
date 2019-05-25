package metrics

import (
	"time"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/tessig/flamingo-mysql/db"
	"github.com/tessig/flamingo-product-rating/src/app"
)

var (
	stat            = stats.Int64("rating/metrics/amount", "Amount of ratings", stats.UnitDimensionless)
	keyProductID, _ = tag.NewKey("productID")
	ticker          *time.Ticker
)

type (
	// Module basic struct
	Module struct{}
)

func init() {
	opencensus.View("rating/metrics/amount", stat, view.LastValue(), keyProductID)
}

// Configure Metrics module
func (m *Module) Configure(injector *dingo.Injector) {
	flamingo.BindEventSubscriber(injector).To(&ShutdownMetrics{})
	flamingo.BindEventSubscriber(injector).To(&StartUpMetrics{})
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(app.Module),
		new(opencensus.Module),
		new(db.Module),
	}
}
