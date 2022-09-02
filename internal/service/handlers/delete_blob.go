package handlers

import (
	"blob-svc/internal/data"
	"blob-svc/internal/service/helpers"
	"blob-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func DeleteBlob(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewDeleteBlobRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = helpers.BlobsQ(r).Transaction(func(q data.BlobsQ) error {
		err = q.Delete(request.BlobID)
		if err != nil {
			return errors.Wrap(err, "failed to delete blob")
		}

		return nil
	})
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
