package data

import (
	"encoding/json"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type BlobsQ interface {
	New() BlobsQ

	Get() (*Blob, error)
	Select() ([]Blob, error)
	Update() ([]Blob, error)

	Transaction(fn func(q BlobsQ) error) error

	Insert(data Blob) (Blob, error)
	Delete(id int64) error

	Page(pageParams pgdb.OffsetPageParams) BlobsQ

	FilterByID(ids ...int64) BlobsQ
	FilterByOwnerAddress(ownerAddresses ...string) BlobsQ
}

type Blob struct {
	ID           int64           `db:"id" structs:"-"`
	Information  json.RawMessage `db:"information" structs:"information"`
	OwnerAddress string          `db:"owner_address" structs:"owner_address"`
}
