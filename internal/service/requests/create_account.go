package requests

import (
	"blob-svc/resources"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

type CreateAccountRequest struct {
	Data       resources.CreateAccount
	Included   resources.Included
	SignerData IncludedSignerData
}

type IncludedSignerData []struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		RoleID   int `json:"role_id"`
		Weight   int `json:"weight"`
		Identity int `json:"identity"`
	} `json:"attributes"`
}

func NewCreateAccountRequest(r *http.Request) (CreateAccountRequest, error) {
	var request CreateAccountRequest
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
