package handlers

import (
	"github.com/paydex-core/paydex-go/clients/federation"
	hc "github.com/paydex-core/paydex-go/clients/horizonclient"
	"github.com/paydex-core/paydex-go/clients/paydextoml"
	"github.com/paydex-core/paydex-go/services/bridge/internal/config"
	"github.com/paydex-core/paydex-go/services/bridge/internal/db"
	"github.com/paydex-core/paydex-go/services/bridge/internal/listener"
	"github.com/paydex-core/paydex-go/services/bridge/internal/submitter"
	"github.com/paydex-core/paydex-go/support/http"
)

// RequestHandler implements bridge server request handlers
type RequestHandler struct {
	Config               *config.Config                          `inject:""`
	Client               http.SimpleHTTPClientInterface          `inject:""`
	Horizon              hc.ClientInterface                      `inject:""`
	Database             db.Database                             `inject:""`
	PaydexTomlResolver   paydextoml.ClientInterface              `inject:""`
	FederationResolver   federation.ClientInterface              `inject:""`
	TransactionSubmitter submitter.TransactionSubmitterInterface `inject:""`
	PaymentListener      *listener.PaymentListener               `inject:""`
}
