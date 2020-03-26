package handlers

import (
	"github.com/paydex-core/paydex-go/clients/paydextoml"
	"strconv"
	"time"

	"github.com/paydex-core/paydex-go/clients/federation"
	"github.com/paydex-core/paydex-go/services/compliance/internal/config"
	"github.com/paydex-core/paydex-go/services/compliance/internal/crypto"
	"github.com/paydex-core/paydex-go/services/compliance/internal/db"
	"github.com/paydex-core/paydex-go/support/http"
)

// RequestHandler implements compliance server request handlers
type RequestHandler struct {
	Config                  *config.Config                 `inject:""`
	Client                  http.SimpleHTTPClientInterface `inject:""`
	Database                db.Database                    `inject:""`
	SignatureSignerVerifier crypto.SignerVerifierInterface `inject:""`
	PaydexTomlResolver     paydextoml.ClientInterface    `inject:""`
	FederationResolver      federation.ClientInterface     `inject:""`
	NonceGenerator          NonceGeneratorInterface        `inject:""`
}

type NonceGeneratorInterface interface {
	Generate() string
}

type NonceGenerator struct{}

func (n *NonceGenerator) Generate() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

type TestNonceGenerator struct{}

func (n *TestNonceGenerator) Generate() string {
	return "nonce"
}
