package protocols

import (
	"fmt"

	shared "github.com/paydex-core/paydex-go/services/internal/bridge-compliance-shared"
	"github.com/paydex-core/paydex-go/support/errors"
	"github.com/paydex-core/paydex-go/txnbuild"
)

// ToBaseAsset transforms Asset to github.com/paydex-core/go-paydex-base/build.Asset
func (a Asset) ToBaseAsset() txnbuild.Asset {
	if a.Code == "" && a.Issuer == "" {
		return txnbuild.NativeAsset{}
	}
	return txnbuild.CreditAsset{Code: a.Code, Issuer: a.Issuer}
}

// String returns string representation of this asset
func (a Asset) String() string {
	return fmt.Sprintf("Code: %s, Issuer: %s", a.Code, a.Issuer)
}

// Validate checks if asset params are correct.
func (a Asset) Validate() error {
	if a.Code != "" && a.Issuer != "" {
		if !shared.IsValidAssetCode(a.Code) {
			return errors.New("Invalid asset_code")
		}
		if !shared.IsValidAccountID(a.Issuer) {
			return errors.New("Invalid asset_issuer")
		}
	} else if a.Code == "" && a.Issuer == "" {
		// Native
		return nil
	} else {
		return errors.New("Asset code or issuer is missing")
	}

	return nil
}
