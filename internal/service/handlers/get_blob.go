package handlers

import (
	"blob-svc/internal/service/helpers"
	"blob-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetBlob(w http.ResponseWriter, r *http.Request) {
	_, err := requests.NewGetBlobRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
}