package ingest

import (
	"database/sql"
	"fmt"

	"github.com/paydex-core/paydex-go/services/horizon/internal/db2/core"
	"github.com/paydex-core/paydex-go/support/db"
	"github.com/paydex-core/paydex-go/support/errors"
)

// Load runs queries against `core` to fill in the records of the bundle.
func (lb *LedgerBundle) Load(db *db.Session) error {
	q := &core.Q{Session: db}
	// Load Header
	err := q.LedgerHeaderBySequence(&lb.Header, lb.Sequence)
	if err != nil {
		// Remove when Horizon is able to handle gaps in paydex-core DB.
		// More info:
		if err == sql.ErrNoRows {
			return errors.New(fmt.Sprintf("Gap detected in paydex-core database (ledger=%d). More information: https://www.paydex.org/developers/software/known-issues.html#gaps-detected", lb.Sequence))
		}
		return errors.Wrap(err, "failed to load header")
	}

	// Load transactions
	err = q.TransactionsByLedger(&lb.Transactions, lb.Sequence)
	if err != nil {
		return errors.Wrap(err, "failed to load transactions")
	}

	err = q.TransactionFeesByLedger(&lb.TransactionFees, lb.Sequence)
	if err != nil {
		return errors.Wrap(err, "failed to load transaction fees")
	}

	return nil
}
