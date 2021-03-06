package resourceadapter

import (
	"context"

	protocol "github.com/paydex-core/paydex-go/protocols/horizon"
	"github.com/paydex-core/paydex-go/xdr"
)

func PopulateAsset(ctx context.Context, dest *protocol.Asset, asset xdr.Asset) error {
	return asset.Extract(&dest.Type, &dest.Code, &dest.Issuer)
}
