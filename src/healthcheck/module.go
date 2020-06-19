package healthcheck

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
)

type (
	// Module basic struct
	Module struct{}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMap(new(healthcheck.Status), "db").To(new(DB))
}
