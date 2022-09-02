package handlers

import (
	"blob-svc/internal/service/helpers"
	"blob-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func DeleteBlob(w http.ResponseWriter, r *http.Request) {
	_, err := requests.NewDeleteBlobRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
}
