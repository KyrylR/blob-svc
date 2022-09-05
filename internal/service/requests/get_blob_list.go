package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type GetBlobListRequest struct {
	pgdb.OffsetPageParams
	FilterOwner []string `filter:"token"`
}

func NewGetBlobListRequest(r *http.Request) (GetBlobListRequest, error) {
	var request GetBlobListRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
