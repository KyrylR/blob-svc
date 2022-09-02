package requests

import (
	"blob-svc/resources"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type GetBlobListRequest struct {
	Data resources.CreateBlob
}

func NewGetBlobListRequest(r *http.Request) (GetBlobListRequest, error) {
	var request GetBlobListRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, nil
}
