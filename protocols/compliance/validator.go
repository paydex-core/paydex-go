package compliance

import (
	"github.com/asaskevich/govalidator"
	"github.com/paydex-core/paydex-go/address"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
	govalidator.CustomTypeTagMap.Set("paydex_address", govalidator.CustomTypeValidator(isPaydexAddress))
}

func isPaydexAddress(i interface{}, context interface{}) bool {
	addr, ok := i.(string)

	if !ok {
		return false
	}

	_, _, err := address.Split(addr)

	return err == nil
}
