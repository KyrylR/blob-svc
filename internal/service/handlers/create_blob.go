package handlers

import (
	"blob-svc/internal/data"
	"blob-svc/internal/service/helpers"
	"blob-svc/internal/service/requests"
	"blob-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func CreateBlob(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateBlobRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var resultBlob data.Blob

	err = helpers.BlobsQ(r).Transaction(func(q data.BlobsQ) error {
		blob := data.Blob{
			Information: request.Data.Attributes.Information,
		}

		resultBlob, err = q.Insert(blob)
		if err != nil {
			return errors.Wrap(err, "failed to insert blob")
		}

		return nil
	})
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.BlobResponse{
		Data: resources.Blob{
			Key: resources.NewKeyInt64(resultBlob.ID, resources.BLOB),
			Attributes: resources.BlobAttributes{
				Information: resultBlob.Information,
			},
		},
	}
	ape.Render(w, result)
}
