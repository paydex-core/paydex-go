package sse

import (
	"net/http"

	"github.com/paydex-core/go-throttled"
	"github.com/paydex-core/paydex-go/services/horizon/internal/ledger"
	"github.com/paydex-core/paydex-go/support/errors"
)

// StreamHandler represents a stream handling action
type StreamHandler struct {
	RateLimiter  *throttled.HTTPRateLimiter
	LedgerSource ledger.Source
}

// GenerateEventsFunc generates a slice of sse.Event which are sent via
// streaming.
type GenerateEventsFunc func() ([]Event, error)

// ServeStream handles a SSE requests, sending data every time there is a new
// ledger.
func (handler StreamHandler) ServeStream(
	w http.ResponseWriter,
	r *http.Request,
	limit int,
	generateEvents GenerateEventsFunc,
) {
	ctx := r.Context()
	stream := NewStream(ctx, w)
	stream.SetLimit(limit)

	currentLedgerSequence := handler.LedgerSource.CurrentLedger()
	for {
		// Rate limit the request if it's a call to stream since it queries the DB every second.
		rateLimiter := handler.RateLimiter
		if rateLimiter != nil {
			limited, _, err := rateLimiter.RateLimiter.RateLimit(rateLimiter.VaryBy.Key(r), 1)
			if err != nil {
				stream.Err(errors.Wrap(err, "RateLimiter error"))
				return
			}
			if limited {
				stream.Err(ErrRateLimited)
				return
			}
		}

		events, err := generateEvents()
		if err != nil {
			stream.Err(err)
			return
		}
		for _, event := range events {
			if limit <= 0 {
				break
			}
			stream.Send(event)
			limit--
		}

		if limit <= 0 {
			stream.Done()
			return
		}

		// Manually send the preamble in case there are no data events in SSE to trigger a stream.Send call.
		// This method is called every iteration of the loop, but is protected by a sync.Once variable so it's
		// only executed once.
		stream.Init()

		select {
		case currentLedgerSequence = <-handler.LedgerSource.NextLedger(currentLedgerSequence):
			continue
		case <-ctx.Done():
			stream.Done()
			return
		}
	}
}
