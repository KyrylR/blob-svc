package requests

import (
	"blob-svc/resources"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type CreateBlobRequest struct {
	Data resources.CreateBlob
}

func NewCreateBlobRequest(r *http.Request) (CreateBlobRequest, error) {
	var request CreateBlobRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *CreateBlobRequest) validate() error {
	return mergeErrors(validation.Errors{
		"/data/attributes/topic": validation.Validate(&r.Data.Attributes.Information, validation.Required,
			validation.Length(3, 100)),
	}).Filter()
}

func mergeErrors(validationErrors ...validation.Errors) validation.Errors {
	result := make(validation.Errors)
	for _, errs := range validationErrors {
		for key, err := range errs {
			result[key] = err
		}
	}
	return result
}
