package actions

import (
	"net/http"

	horizonContext "github.com/paydex-core/paydex-go/services/horizon/internal/context"
	"github.com/paydex-core/paydex-go/services/horizon/internal/db2/history"
	"github.com/paydex-core/paydex-go/support/db"
	"github.com/paydex-core/paydex-go/support/errors"
)

func historyQFromRequest(request *http.Request) (*history.Q, error) {
	ctx := request.Context()
	session, ok := ctx.Value(&horizonContext.SessionContextKey).(*db.Session)
	if !ok {
		return nil, errors.New("missing session in request context")
	}
	return &history.Q{session}, nil
}
