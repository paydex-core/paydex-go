package paydextoml

import "net/http"

// PaydexTomlMaxSize is the maximum size of paydex.toml file
const PaydexTomlMaxSize = 100 * 1024

// WellKnownPath represents the url path at which the paydex.toml file should
// exist to conform to the federation protocol.
const WellKnownPath = "/.well-known/paydex.toml"

// DefaultClient is a default client using the default parameters
var DefaultClient = &Client{HTTP: http.DefaultClient}

// Client represents a client that is capable of resolving a Paydex.toml file
// using the internet.
type Client struct {
	// HTTP is the http client used when resolving a Paydex.toml file
	HTTP HTTP

	// UseHTTP forces the client to resolve against servers using plain HTTP.
	// Useful for debugging.
	UseHTTP bool
}

type ClientInterface interface {
	GetPaydexToml(domain string) (*Response, error)
	GetPaydexTomlByAddress(addy string) (*Response, error)
}

// HTTP represents the http client that a paydextoml resolver uses to make http
// requests.
type HTTP interface {
	Get(url string) (*http.Response, error)
}

// Response represents the results of successfully resolving a paydex.toml file
type Response struct {
	AuthServer       string `toml:"AUTH_SERVER"`
	FederationServer string `toml:"FEDERATION_SERVER"`
	EncryptionKey    string `toml:"ENCRYPTION_KEY"`
	SigningKey       string `toml:"SIGNING_KEY"`
}

// GetPaydexToml returns paydex.toml file for a given domain
func GetPaydexToml(domain string) (*Response, error) {
	return DefaultClient.GetPaydexToml(domain)
}

// GetPaydexTomlByAddress returns paydex.toml file of a domain fetched from a
// given address
func GetPaydexTomlByAddress(addy string) (*Response, error) {
	return DefaultClient.GetPyadexTomlByAddress(addy)
}
