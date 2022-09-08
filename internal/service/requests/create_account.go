package requests

import (
	"blob-svc/resources"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func NewCreateAccountRequest(r *http.Request) (resources.CreateAccountResponse, error) {
	var request resources.CreateAccountResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, ValidateCreateAccountRequest(request)
}

func ValidateCreateAccountRequest(r resources.CreateAccountResponse) error {
	errs := validation.Errors{
		"/data/":                      validation.Validate(r.Data, validation.Required),
		"/data/relationships/signers": validation.Validate(r.Data.Relationships.Signers),
	}

	return errs.Filter()
}
