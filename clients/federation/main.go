package federation

import (
	"net/http"
	"net/url"

	hc "github.com/paydex-core/paydex-go/clients/horizonclient"
	"github.com/paydex-core/paydex-go/clients/paydextoml"
	proto "github.com/paydex-core/paydex-go/protocols/federation"
)

// FederationResponseMaxSize is the maximum size of response from a federation server
const FederationResponseMaxSize = 100 * 1024

// DefaultTestNetClient is a default federation client for testnet
var DefaultTestNetClient = &Client{
	HTTP:        http.DefaultClient,
	Horizon:     hc.DefaultTestNetClient,
	PaydexTOML: paydextoml.DefaultClient,
}

// DefaultPublicNetClient is a default federation client for pubnet
var DefaultPublicNetClient = &Client{
	HTTP:        http.DefaultClient,
	Horizon:     hc.DefaultPublicNetClient,
	PaydexTOML: paydextoml.DefaultClient,
}

// Client represents a client that is capable of resolving a federation request
// using the internet.
type Client struct {
	PaydexTOML PaydexTOML
	HTTP        HTTP
	Horizon     Horizon
	AllowHTTP   bool
}

type ClientInterface interface {
	LookupByAddress(addy string) (*proto.NameResponse, error)
	LookupByAccountID(aid string) (*proto.IDResponse, error)
	ForwardRequest(domain string, fields url.Values) (*proto.NameResponse, error)
}

// Horizon represents a horizon client that can be consulted for data when
// needed as part of the federation protocol
type Horizon interface {
	HomeDomainForAccount(aid string) (string, error)
}

// HTTP represents the http client that a federation client uses to make http
// requests.
type HTTP interface {
	Get(url string) (*http.Response, error)
}

// PaydexTOML represents a client that can resolve a given domain name to
// paydex.toml file.  The response is used to find the federation server that a
// query should be made against.
type PaydexTOML interface {
	GetPaydexToml(domain string) (*paydextoml.Response, error)
}

// confirm interface conformity
var _ PaydexTOML = paydextoml.DefaultClient
var _ HTTP = http.DefaultClient
var _ ClientInterface = &Client{}
