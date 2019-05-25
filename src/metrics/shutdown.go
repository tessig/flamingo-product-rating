package metrics

import (
	"context"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// ShutdownMetrics handles the graceful app shutdown for the db connection
	ShutdownMetrics struct{}
)

// Notify handles the incoming event if it is an AppShutdownEvent and closes the db connection
func (s *ShutdownMetrics) Notify(_ context.Context, event flamingo.Event) {
	switch event.(type) {
	case *flamingo.ShutdownEvent:
		ticker.Stop()
	}
}
