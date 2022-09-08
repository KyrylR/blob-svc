package accountcreator

import (
	"blob-svc/internal/data"
	"blob-svc/xdrbuild"
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector"
)

type AccountCreator struct {
	Tx             *xdrbuild.Transaction
	Horizon        *horizon.Connector
	SystemSettings data.SystemSettings
}

type AccountDetails struct {
	AccountID string
	Signers   []data.AccountSigner
}

func New(tx *xdrbuild.Transaction,
	horizon *horizon.Connector, systemSettings data.SystemSettings) AccountCreator {
	return AccountCreator{
		Tx:             tx,
		Horizon:        horizon,
		SystemSettings: systemSettings,
	}
}

func (c AccountCreator) CreateAccount(ctx context.Context, accountID string, signers []data.AccountSigner) error {
	tx := c.Tx

	tx, _, err := c.account(tx, accountID, signers)
	if err != nil {
		return errors.Wrap(err, "failed to craft account operation")
	}

	envelope, err := tx.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to build tx envelope")
	}

	result := c.Horizon.Submitter().Submit(ctx, envelope)
	if result.Err != nil {
		return convertTXSubmitError(result)
	}

	return nil
}

func (c AccountCreator) CreateAccountBatch(ctx context.Context, accounts []AccountDetails) error {
	if len(accounts) > 100 {
		return errors.New("too many accounts to create")
	}
	tx := c.Tx
	var err error
	for _, acc := range accounts {
		tx, _, err = c.account(tx, acc.AccountID, acc.Signers)
		if err != nil {
			return errors.Wrap(err, "failed to craft account operation")
		}
	}

	envelope, err := tx.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to build tx envelope")
	}

	result := c.Horizon.Submitter().Submit(ctx, envelope)
	if result.Err != nil {
		return convertTXSubmitError(result)
	}

	return nil
}

func (c AccountCreator) Do(ctx context.Context, wallet *data.Wallet, signers []data.AccountSigner) (uint64, error) {
	tx := c.Tx

	tx, roleID, err := c.account(tx, string(wallet.AccountID), signers)
	if err != nil {
		return 0, errors.Wrap(err, "failed to craft account operation")
	}

	//tx, err = c.invites(tx, wallet)
	//if err != nil {
	//	return 0, errors.Wrap(err, "failed to craft invites operations")
	//}

	envelope, err := tx.Marshal()
	if err != nil {
		return 0, errors.Wrap(err, "failed to build tx envelope")
	}

	result := c.Horizon.Submitter().Submit(ctx, envelope)
	if result.Err != nil {
		return 0, convertTXSubmitError(result)
	}
	return roleID, nil
}

func checkTxResult(result horizon.SubmitResult) error {
	switch result.TXCode {
	case "tx_success":
		return nil
	case "tx_failed":
		return validation.Errors{"create account tx": errors.New(fmt.Sprint("transaction failed with op codes: ", result.OpCodes))}
	default:
		return validation.Errors{"create account tx": errors.New(fmt.Sprint("transaction failed with tx code: ", result.TXCode))}
	}
}

func (c AccountCreator) account(tx *xdrbuild.Transaction, accountID string, accountSigners []data.AccountSigner) (*xdrbuild.Transaction, uint64, error) {
	accountRole, err := c.SystemSettings.DefaultAccountRole()
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get account role")
	}

	tx = tx.Op(&xdrbuild.CreateAccount{
		Destination: accountID,
		RoleID:      accountRole,
		Signers:     c.getSigners(accountSigners),
	})
	return tx, accountRole, nil
}

func (c AccountCreator) getSigners(accountSigners []data.AccountSigner) []xdrbuild.SignerData {
	var signers []xdrbuild.SignerData
	for _, accountSigner := range accountSigners {
		signers = append(signers, xdrbuild.SignerData{
			PublicKey: accountSigner.SignerID,
			RoleID:    accountSigner.RoleID,
			Weight:    accountSigner.Weight,
			Identity:  accountSigner.Identity,
			Details:   Details{},
		})
	}
	return signers
}

type Details map[string]interface{}

func (d Details) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(d))
}

func convertTXSubmitError(s horizon.SubmitResult) (err error) {
	badRequests := map[string]struct{}{
		"op_invalid_destination": {},
		"op_already_exists":      {},
	}

	if len(s.OpCodes) > 0 {
		// only create account op codes are handled
		opCode := s.OpCodes[0]
		if _, ok := badRequests[opCode]; ok {
			return validation.Errors{
				"/data/attributes/account_id": errors.New(
					fmt.Sprintf("'%s' op code received on create account op", opCode),
				),
			}
		}
	}

	return checkTxResult(s)
}
