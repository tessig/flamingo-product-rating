package interfaces

import (
	"context"
	"math/rand"
	"strconv"
	"time"
)

type (
	// RandomIntFunc provides a random integer template function
	RandomIntFunc struct {
		debugMode bool
	}

	// ForFunc provides a template function returning a slice of ints
	ForFunc struct{}

	// BarTypeFunc provides a template function which determines the bar type for a stars bar
	BarTypeFunc struct{}
)

// Inject dependencies
func (f *RandomIntFunc) Inject(
	c *struct {
		DebugMode bool `inject:"config:debug.mode"`
	},
) {
	f.debugMode = c.DebugMode
}

// Func returns the function returning a random int between min and max
func (f *RandomIntFunc) Func(_ context.Context) interface{} {
	return func(min, max int) int {
		if !f.debugMode {
			rand.Seed(time.Now().UnixNano())
		}

		return rand.Intn(max-min) + min
	}
}

// Func returns the function returning a slice of strings from start to end (including the edges)
func (f *ForFunc) Func(_ context.Context) interface{} {
	return func(start, end int) []string {
		if end < start {
			return []string{}
		}
		s := make([]string, end-start+1)
		for i := start; i <= end; i++ {
			s[i-start] = strconv.Itoa(i)
		}

		return s
	}
}

// Func returns the function mapping bar types to star amounts
func (f *BarTypeFunc) Func(_ context.Context) interface{} {
	return func(stars int) string {
		switch stars {
		case 5:
			return "success"
		case 4:
			return "primary"
		case 3:
			return "info"
		case 2:
			return "warning"
		case 1:
			return "danger"
		}

		return "success"
	}
}
