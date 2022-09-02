package requests

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type DeleteBlobRequest struct {
	BlobID int64 `url:"-"`
}

func NewDeleteBlobRequest(r *http.Request) (DeleteBlobRequest, error) {
	var request DeleteBlobRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, nil
}
