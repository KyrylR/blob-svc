package types

import (
	"errors"

	"github.com/spf13/cast"

	"gitlab.com/tokend/go/strkey"
)

var (
	ErrAddressInvalid = errors.New("address is invalid")
)

type Address string

func (a Address) Validate() error {
	_, err := strkey.Decode(strkey.VersionByteAccountID, string(a))
	if err != nil {
		return ErrAddressInvalid
	}
	return nil
}

func (a Address) String() string {
	return string(a)
}

var IsAddress = &isAddress{}

type isAddress struct{}

func (ia *isAddress) Validate(value interface{}) error {
	a, err := cast.ToStringE(value)
	if err != nil {
		return err
	}
	address := Address(a)
	return address.Validate()
}
