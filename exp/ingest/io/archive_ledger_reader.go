package io

import "github.com/paydex-core/paydex-go/xdr"

// ArchiveLedgerReader placeholder
type ArchiveLedgerReader interface {
	GetSequence() uint32
	Read() (bool, xdr.Transaction, xdr.TransactionResult, error)
}
