package handlers

import (
	"blob-svc/internal/data"
	"blob-svc/internal/service/helpers"
	"blob-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetBlobList(w http.ResponseWriter, r *http.Request) {
	blobs, err := helpers.BlobsQ(r).Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get blobs")
		ape.Render(w, problems.InternalError())
		return
	}

	response := resources.BlobListResponse{
		Data: newBlobsList(blobs),
	}
	ape.Render(w, response)
}

func newBlobsList(blobs []data.Blob) []resources.Blob {
	result := make([]resources.Blob, len(blobs))
	for i, blob := range blobs {
		result[i] = resources.Blob{
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
		}
	}
	return result
}
