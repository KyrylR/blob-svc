package data

import (
	"blob-svc/xdrbuild"
	"gitlab.com/tokend/keypair"
	"gitlab.com/tokend/regources"
)

//go:generate mockery -case underscore -name Info

type Infobuilder func(info Info, source keypair.Address) *xdrbuild.Transaction

type Info interface {
	Info() (*regources.Info, error)
}
