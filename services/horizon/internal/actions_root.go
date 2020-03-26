package horizon

import (
	"github.com/paydex-core/paydex-go/protocols/horizon"
	"github.com/paydex-core/paydex-go/services/horizon/internal/actions"
	"github.com/paydex-core/paydex-go/services/horizon/internal/ledger"
	"github.com/paydex-core/paydex-go/services/horizon/internal/resourceadapter"
	"github.com/paydex-core/paydex-go/support/render/hal"
)

// Interface verification
var _ actions.JSONer = (*RootAction)(nil)

// RootAction provides a summary of the horizon instance and links to various
// useful endpoints
type RootAction struct {
	Action
}

// JSON renders the json response for RootAction
func (action *RootAction) JSON() error {
	var res horizon.Root
	templates := map[string]string{
		"accounts":           actions.AccountsQuery{}.URITemplate(),
		"offers":             actions.OffersQuery{}.URITemplate(),
		"strictReceivePaths": StrictReceivePathsQuery{}.URITemplate(),
		"strictSendPaths":    FindFixedPathsQuery{}.URITemplate(),
	}
	resourceadapter.PopulateRoot(
		action.R.Context(),
		&res,
		ledger.CurrentState(),
		action.App.horizonVersion,
		action.App.coreVersion,
		action.App.config.NetworkPassphrase,
		action.App.currentProtocolVersion,
		action.App.coreSupportedProtocolVersion,
		action.App.config.FriendbotURL,
		action.App.config.EnableExperimentalIngestion,
		templates,
	)

	hal.Render(action.W, res)
	return action.Err
}
