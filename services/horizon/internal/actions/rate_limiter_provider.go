package actions

import "github.com/paydex-core/go-throttled"

// RateLimiterProvider is an interface that provides access to the type's HTTPRateLimiter.
type RateLimiterProvider interface {
	GetRateLimiter() *throttled.HTTPRateLimiter
}
