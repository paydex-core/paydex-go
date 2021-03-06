package resourceadapter

import (
	"context"
	"net/url"

	"github.com/paydex-core/paydex-go/protocols/horizon"
	"github.com/paydex-core/paydex-go/services/horizon/internal/httpx"
	"github.com/paydex-core/paydex-go/services/horizon/internal/ledger"
	"github.com/paydex-core/paydex-go/support/render/hal"
)

// Populate fills in the details
func PopulateRoot(
	ctx context.Context,
	dest *horizon.Root,
	ledgerState ledger.State,
	hVersion, cVersion string,
	passphrase string,
	currentProtocolVersion int32,
	coreSupportedProtocolVersion int32,
	friendBotURL *url.URL,
	experimentalIngestionEnabled bool,
	templates map[string]string,
) {
	dest.ExpHorizonSequence = ledgerState.ExpHistoryLatest
	dest.HorizonSequence = ledgerState.HistoryLatest
	dest.HistoryElderSequence = ledgerState.HistoryElder
	dest.CoreSequence = ledgerState.CoreLatest
	dest.HorizonVersion = hVersion
	dest.PaydexCoreVersion = cVersion
	dest.NetworkPassphrase = passphrase
	dest.CurrentProtocolVersion = currentProtocolVersion
	dest.CoreSupportedProtocolVersion = coreSupportedProtocolVersion

	lb := hal.LinkBuilder{Base: httpx.BaseURL(ctx)}
	if friendBotURL != nil {
		friendbotLinkBuild := hal.LinkBuilder{Base: friendBotURL}
		l := friendbotLinkBuild.Link("{?addr}")
		dest.Links.Friendbot = &l
	}

	dest.Links.Account = lb.Link("/accounts/{account_id}")
	dest.Links.AccountTransactions = lb.PagedLink("/accounts/{account_id}/transactions")
	dest.Links.Assets = lb.Link("/assets{?asset_code,asset_issuer,cursor,limit,order}")
	dest.Links.Metrics = lb.Link("/metrics")

	if experimentalIngestionEnabled {
		accountsLink := lb.Link(templates["accounts"])
		offerLink := lb.Link("/offers/{offer_id}")
		offersLink := lb.Link(templates["offers"])
		strictReceivePaths := lb.Link(templates["strictReceivePaths"])
		strictSendPaths := lb.Link(templates["strictSendPaths"])
		dest.Links.Accounts = &accountsLink
		dest.Links.Offer = &offerLink
		dest.Links.Offers = &offersLink
		dest.Links.StrictReceivePaths = &strictReceivePaths
		dest.Links.StrictSendPaths = &strictSendPaths
	}

	dest.Links.OrderBook = lb.Link("/order_book{?selling_asset_type,selling_asset_code,selling_asset_issuer,buying_asset_type,buying_asset_code,buying_asset_issuer,limit}")
	dest.Links.Self = lb.Link("/")
	dest.Links.Transaction = lb.Link("/transactions/{hash}")
	dest.Links.Transactions = lb.PagedLink("/transactions")
}
