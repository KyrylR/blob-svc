package pg

import (
	"blob-svc/internal/data"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const blobsTableName = "blobs"

func NewBlobsQ(db *pgdb.DB) data.BlobsQ {
	return &blobsQ{
		db:        db.Clone(),
		sql:       sq.Select("blobs.*").From(blobsTableName),
		sqlUpdate: sq.Update(blobsTableName).Suffix("returning *"),
	}
}

type blobsQ struct {
	db        *pgdb.DB
	sql       sq.SelectBuilder
	sqlUpdate sq.UpdateBuilder
}

func (q *blobsQ) New() data.BlobsQ {
	return NewBlobsQ(q.db)
}

func (q *blobsQ) Get() (*data.Blob, error) {
	var result data.Blob
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *blobsQ) Select() ([]data.Blob, error) {
	var result []data.Blob
	err := q.db.Select(&result, q.sql)
	return result, err
}

func (q *blobsQ) Update() ([]data.Blob, error) {
	var result []data.Blob
	err := q.db.Select(&result, q.sqlUpdate)

	return result, err
}

func (q *blobsQ) Transaction(fn func(q data.BlobsQ) error) error {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *blobsQ) Insert(blob data.Blob) (data.Blob, error) {
	clauses := structs.Map(blob)
	clauses["information"] = blob.Information
	clauses["owner_address"] = blob.OwnerAddress

	var result data.Blob
	stmt := sq.Insert(blobsTableName).SetMap(clauses).Suffix("returning *")
	err := q.db.Get(&result, stmt)

	return result, err
}

func (q *blobsQ) Delete(id int64) error {
	stmt := sq.Delete(blobsTableName).Where(sq.Eq{"id": id})
	err := q.db.Exec(stmt)
	return err
}

func (q *blobsQ) FilterByID(ids ...int64) data.BlobsQ {
	q.sql = q.sql.Where(sq.Eq{"id": ids})
	return q
}
