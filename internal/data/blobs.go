package data

import "encoding/json"

type BlobsQ interface {
	New() BlobsQ

	Get() (*Blob, error)
	Select() ([]Blob, error)
	Update() ([]Blob, error)

	Transaction(fn func(q BlobsQ) error) error

	Insert(data Blob) (Blob, error)
	Delete(id int64) error

	FilterByID(id ...int64) BlobsQ
}

type Blob struct {
	ID          int64           `db:"id" structs:"-"`
	Information json.RawMessage `db:"information" structs:"information"`
}
