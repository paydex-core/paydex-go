package horizon

import (
	"net/http"

	hProblem "github.com/paydex-core/paydex-go/services/horizon/internal/render/problem"
	"github.com/paydex-core/paydex-go/support/render/problem"
)

// RateLimitExceededAction renders a 429 response
type RateLimitExceededAction struct {
	Action
}

func (action RateLimitExceededAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(w, r)
	problem.Render(action.R.Context(), action.W, hProblem.RateLimitExceeded)
}
