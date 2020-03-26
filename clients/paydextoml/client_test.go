package paydextoml

import (
	"net/http"
	"strings"
	"testing"

	"github.com/paydex-core/paydex-go/support/http/httptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientURL(t *testing.T) {
	//HACK:  we're testing an internal method rather than setting up a http client
	//mock.

	c := &Client{UseHTTP: false}
	assert.Equal(t, "https://paydex.org/.well-known/paydex.toml", c.url("paydex.org"))

	c = &Client{UseHTTP: true}
	assert.Equal(t, "http://paydex.org/.well-known/paydex.toml", c.url("paydex.org"))
}

func TestClient(t *testing.T) {
	h := httptest.NewClient()
	c := &Client{HTTP: h}

	// happy path
	h.
		On("GET", "https://paydex.org/.well-known/paydex.toml").
		ReturnString(http.StatusOK,
			`FEDERATION_SERVER="https://localhost/federation"`,
		)
	stoml, err := c.GetPaydexToml("paydex.org")
	require.NoError(t, err)
	assert.Equal(t, "https://localhost/federation", stoml.FederationServer)

	// paydex.toml exceeds limit
	h.
		On("GET", "https://toobig.org/.well-known/paydex.toml").
		ReturnString(http.StatusOK,
			`FEDERATION_SERVER="https://localhost/federation`+strings.Repeat("0", PaydexTomlMaxSize)+`"`,
		)
	stoml, err = c.GetPaydexToml("toobig.org")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "paydex.toml response exceeds")
	}

	// not found
	h.
		On("GET", "https://missing.org/.well-known/paydex.toml").
		ReturnNotFound()
	stoml, err = c.GetPaydexToml("missing.org")
	assert.EqualError(t, err, "http request failed with non-200 status code")

	// invalid toml
	h.
		On("GET", "https://json.org/.well-known/paydex.toml").
		ReturnJSON(http.StatusOK, map[string]string{"hello": "world"})
	stoml, err = c.GetPaydexToml("json.org")

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "toml decode failed")
	}
}
