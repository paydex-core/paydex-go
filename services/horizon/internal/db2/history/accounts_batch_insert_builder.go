package history

import "github.com/paydex-core/paydex-go/xdr"

func (i *accountsBatchInsertBuilder) Add(account xdr.AccountEntry, lastModifiedLedger xdr.Uint32) error {
	return i.builder.Row(accountToMap(account, lastModifiedLedger))
}

func (i *accountsBatchInsertBuilder) Exec() error {
	return i.builder.Exec()
}
