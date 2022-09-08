package horizon

import (
	"blob-svc/internal/data"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/regources/v2"
)

type AccountQ struct {
	log     *logan.Entry
	horizon *horizon.Connector
}

func NewAccountQ(log *logan.Entry, horizon *horizon.Connector) *AccountQ {
	return &AccountQ{
		horizon: horizon,
		log:     log,
	}
}

func (q *AccountQ) Account(address string) (*regources.Account, error) {
	endpoint := fmt.Sprintf("/v3/accounts/%s", address)
	response, err := q.horizon.Client().Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if response == nil {
		return nil, nil
	}

	var accountResponse regources.AccountResponse
	if err := json.Unmarshal(response, &accountResponse); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	return &accountResponse.Data, nil
}

func (q *AccountQ) Signers(address string) ([]data.Signer, error) {
	Signers := func(address string) ([]regources.Signer, error) {
		endpoint := fmt.Sprintf("/v3/accounts/%s/signers", address)
		response, err := q.horizon.Client().Get(endpoint)
		if err != nil {
			return nil, errors.Wrap(err, "request failed")
		}

		if response == nil {
			return nil, nil
		}

		var signerResponse regources.SignersResponse
		if err := json.Unmarshal(response, &signerResponse); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal")
		}
		return signerResponse.Data, nil
	}

	signers, err := Signers(address)
	if err != nil {
		q.log.WithError(err).Error("failed to get signers")
		return nil, err
	}

	if signers == nil {
		return nil, nil
	}

	// TODO share resource
	result := make([]data.Signer, 0, len(signers))
	for _, signer := range signers {
		roleID, err := strconv.ParseUint(signer.Relationships.Role.Data.ID, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse signer role id")
		}

		result = append(result, data.Signer{
			// TODO fixme bledb
			AccountID: signer.ID,
			Weight:    int(signer.Attributes.Weight),
			Identity:  uint32(signer.Attributes.Identity),
			Role:      roleID,
		})
	}
	return result, nil
}
