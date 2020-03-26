package resourceadapter

import (
	"context"
	"testing"

	"github.com/paydex-core/paydex-go/protocols/horizon"
	protocol "github.com/paydex-core/paydex-go/protocols/horizon"
	"github.com/paydex-core/paydex-go/services/horizon/internal/db2/assets"
	"github.com/paydex-core/paydex-go/services/horizon/internal/db2/history"
	"github.com/paydex-core/paydex-go/xdr"
	"github.com/stretchr/testify/assert"
)

func TestLargeAmount(t *testing.T) {
	row := assets.AssetStatsR{
		SortKey:     "",
		Type:        "credit_alphanum4",
		Code:        "XIM",
		Issuer:      "GBZ35ZJRIKJGYH5PBKLKOZ5L6EXCNTO7BKIL7DAVVDFQ2ODJEEHHJXIM",
		Amount:      "100000000000000000000", // 10T
		NumAccounts: 429,
		Flags:       0,
		Toml:        "https://xim.com/.well-known/paydex.toml",
	}
	var res protocol.AssetStat
	err := PopulateAssetStat(context.Background(), &res, row)
	assert.NoError(t, err)

	assert.Equal(t, "credit_alphanum4", res.Type)
	assert.Equal(t, "XIM", res.Code)
	assert.Equal(t, "GBZ35ZJRIKJGYH5PBKLKOZ5L6EXCNTO7BKIL7DAVVDFQ2ODJEEHHJXIM", res.Issuer)
	assert.Equal(t, "10000000000000.0000000", res.Amount)
	assert.Equal(t, int32(429), res.NumAccounts)
	assert.Equal(t, "https://xim.com/.well-known/paydex.toml", res.Links.Toml.Href)
}

func TestPopulateExpAssetStat(t *testing.T) {
	row := history.ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetCode:   "XIM",
		AssetIssuer: "GBZ35ZJRIKJGYH5PBKLKOZ5L6EXCNTO7BKIL7DAVVDFQ2ODJEEHHJXIM",
		Amount:      "100000000000000000000", // 10T
		NumAccounts: 429,
	}
	issuer := history.AccountEntry{
		AccountID:  "GBZ35ZJRIKJGYH5PBKLKOZ5L6EXCNTO7BKIL7DAVVDFQ2ODJEEHHJXIM",
		Flags:      0,
		HomeDomain: "xim.com",
	}

	var res protocol.AssetStat
	err := PopulateExpAssetStat(context.Background(), &res, row, issuer)
	assert.NoError(t, err)

	assert.Equal(t, "credit_alphanum4", res.Type)
	assert.Equal(t, "XIM", res.Code)
	assert.Equal(t, "GBZ35ZJRIKJGYH5PBKLKOZ5L6EXCNTO7BKIL7DAVVDFQ2ODJEEHHJXIM", res.Issuer)
	assert.Equal(t, "10000000000000.0000000", res.Amount)
	assert.Equal(t, int32(429), res.NumAccounts)
	assert.Equal(t, horizon.AccountFlags{}, res.Flags)
	assert.Equal(t, "https://xim.com/.well-known/paydex.toml", res.Links.Toml.Href)
	assert.Equal(t, row.PagingToken(), res.PagingToken())

	issuer.HomeDomain = ""
	issuer.Flags = uint32(xdr.AccountFlagsAuthRequiredFlag) |
		uint32(xdr.AccountFlagsAuthImmutableFlag)

	err = PopulateExpAssetStat(context.Background(), &res, row, issuer)
	assert.NoError(t, err)

	assert.Equal(t, "credit_alphanum4", res.Type)
	assert.Equal(t, "XIM", res.Code)
	assert.Equal(t, "GBZ35ZJRIKJGYH5PBKLKOZ5L6EXCNTO7BKIL7DAVVDFQ2ODJEEHHJXIM", res.Issuer)
	assert.Equal(t, "10000000000000.0000000", res.Amount)
	assert.Equal(t, int32(429), res.NumAccounts)
	assert.Equal(
		t,
		horizon.AccountFlags{
			AuthRequired:  true,
			AuthImmutable: true,
		},
		res.Flags,
	)
	assert.Equal(t, "", res.Links.Toml.Href)
	assert.Equal(t, row.PagingToken(), res.PagingToken())
}
