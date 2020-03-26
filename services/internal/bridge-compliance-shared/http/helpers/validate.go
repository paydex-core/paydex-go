package helpers

import (
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/paydex-core/paydex-go/address"
	"github.com/paydex-core/paydex-go/amount"
	"github.com/paydex-core/paydex-go/strkey"
)

func init() {
	govalidator.CustomTypeTagMap.Set("paydex_accountid", govalidator.CustomTypeValidator(isPaydexAccountID))
	govalidator.CustomTypeTagMap.Set("paydex_seed", govalidator.CustomTypeValidator(isPaydexSeed))
	govalidator.CustomTypeTagMap.Set("paydex_asset_code", govalidator.CustomTypeValidator(isPaydexAssetCode))
	govalidator.CustomTypeTagMap.Set("paydex_address", govalidator.CustomTypeValidator(isPaydexAddress))
	govalidator.CustomTypeTagMap.Set("paydex_amount", govalidator.CustomTypeValidator(isPaydexAmount))
	govalidator.CustomTypeTagMap.Set("paydex_destination", govalidator.CustomTypeValidator(isPaydexDestination))

}

func Validate(request Request, params ...interface{}) error {
	valid, err := govalidator.ValidateStruct(request)

	if !valid {
		fields := govalidator.ErrorsByField(err)
		for field, errorValue := range fields {
			switch {
			case errorValue == "non zero value required":
				return NewMissingParameter(field)
			case strings.HasSuffix(errorValue, "does not validate as paydex_accountid"):
				return NewInvalidParameterError(field, "Account ID must start with `G` and contain 56 alphanum characters.")
			case strings.HasSuffix(errorValue, "does not validate as paydex_seed"):
				return NewInvalidParameterError(field, "Account secret must start with `S` and contain 56 alphanum characters.")
			case strings.HasSuffix(errorValue, "does not validate as paydex_asset_code"):
				return NewInvalidParameterError(field, "Asset code must be 1-12 alphanumeric characters.")
			case strings.HasSuffix(errorValue, "does not validate as paydex_address"):
				return NewInvalidParameterError(field, "Paydex address must be of form user*domain.com")
			case strings.HasSuffix(errorValue, "does not validate as paydex_destination"):
				return NewInvalidParameterError(field, "Paydex destination must be of form user*domain.com or start with `G` and contain 56 alphanum characters.")
			case strings.HasSuffix(errorValue, "does not validate as paydex_amount"):
				return NewInvalidParameterError(field, "Amount must be positive and have up to 7 decimal places.")
			default:
				return NewInvalidParameterError(field, errorValue)
			}
		}
	}

	return request.Validate(params...)
}

// These are copied from support/config. Should we move them to /strkey maybe?
func isPaydexAccountID(i interface{}, context interface{}) bool {
	enc, ok := i.(string)

	if !ok {
		return false
	}

	_, err := strkey.Decode(strkey.VersionByteAccountID, enc)
	return err == nil
}

func isPaydexSeed(i interface{}, context interface{}) bool {
	enc, ok := i.(string)

	if !ok {
		return false
	}

	_, err := strkey.Decode(strkey.VersionByteSeed, enc)
	return err == nil
}

func isPaydexAssetCode(i interface{}, context interface{}) bool {
	code, ok := i.(string)

	if !ok {
		return false
	}

	if !govalidator.IsByteLength(code, 1, 12) {
		return false
	}

	if !govalidator.IsAlphanumeric(code) {
		return false
	}

	return true
}

func isPaydexAddress(i interface{}, context interface{}) bool {
	addr, ok := i.(string)

	if !ok {
		return false
	}

	_, _, err := address.Split(addr)
	return err == nil
}

func isPaydexAmount(i interface{}, context interface{}) bool {
	am, ok := i.(string)

	if !ok {
		return false
	}

	_, err := amount.Parse(am)
	return err == nil
}

// isPaydexDestination checks if `i` is either account public key or Paydex address.
func isPaydexDestination(i interface{}, context interface{}) bool {
	dest, ok := i.(string)

	if !ok {
		return false
	}

	_, err1 := strkey.Decode(strkey.VersionByteAccountID, dest)
	_, _, err2 := address.Split(dest)

	if err1 != nil && err2 != nil {
		return false
	}

	return true
}
