package handlers

import (
	"blob-svc/internal/service/helpers"
	"blob-svc/internal/service/requests"
	"blob-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetBlob(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetBlobRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	blob, err := helpers.BlobsQ(r).FilterByID(request.BlobID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get blob from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if blob == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	result := resources.BlobResponse{
		Data: resources.Blob{
			Key: resources.NewKeyInt64(blob.ID, resources.BLOB),
			Attributes: resources.BlobAttributes{
				Information: blob.Information,
			},
			Relationships: resources.BlobRelationships{
				Owner: resources.Relation{
					Data: &resources.Key{
						ID: blob.OwnerAddress,
					},
				},
			},
		},
	}

	ape.Render(w, result)
}
