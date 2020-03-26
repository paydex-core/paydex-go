package resourceadapter

import (
	protocol "github.com/paydex-core/paydex-go/protocols/horizon"
	"github.com/paydex-core/paydex-go/services/horizon/internal/db2/core"
)

func PopulateAccountFlags(dest *protocol.AccountFlags, row core.Account) {
	dest.AuthRequired = row.IsAuthRequired()
	dest.AuthRevocable = row.IsAuthRevocable()
	dest.AuthImmutable = row.IsAuthImmutable()
}
