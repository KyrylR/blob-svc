package requests

import (
	"blob-svc/resources"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type DeleteBlobRequest struct {
	Data resources.CreateBlob
}

func NewDeleteBlobRequest(r *http.Request) (DeleteBlobRequest, error) {
	var request DeleteBlobRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, nil
}
