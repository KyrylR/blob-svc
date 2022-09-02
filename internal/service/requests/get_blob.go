package requests

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type GetBlobRequest struct {
	BlobID int64 `url:"-"`
}

func NewGetBlobRequest(r *http.Request) (GetBlobRequest, error) {
	var request GetBlobRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, nil
}
