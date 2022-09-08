package data

import (
	"encoding/json"
	"fmt"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
	regources "gitlab.com/tokend/regources/generated"
)

type KeyValue struct {
	KYCRecoveryEnabled     string `fig:"kyc_recovery_enabled"`
	KYCRecoverySignerRole  string `fig:"kyc_recovery_signer_role"`
	LicenseAdminSignerRole string `fig:"license_admin_signer_role"`
	DefaultAccountRole     string `fig:"account_role"`
	DefaultSignerRole      string `fig:"signer_role"`
}

type SystemSettings interface {
	KYCRecoveryEnabled() (bool, error)
	KYCRecoverySignerRole() (uint64, error)
	LicenseAdminSignerRole() (uint64, error)
	DefaultAccountRole() (uint64, error)
	DefaultSignerRole() (uint64, error)
}

type SystemSettingsQ struct {
	Horizon  *horizon.Connector
	KeyValue KeyValue

	// Cached value of KycRecoverySignerRole
	kycRecoverySignerRole *uint64
}

func NewSystemSettingsQ(horizon *horizon.Connector, kv KeyValue) SystemSettings {
	return &SystemSettingsQ{
		Horizon:  horizon,
		KeyValue: kv,
	}
}

func (q *SystemSettingsQ) KYCRecoverySignerRole() (uint64, error) {
	if q.kycRecoverySignerRole != nil {
		return *q.kycRecoverySignerRole, nil
	}

	signerRole, err := keyValueUint64(q.Horizon, q.KeyValue.KYCRecoverySignerRole)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get role")
	}

	if signerRole == nil {
		return 0, errors.New("role id does not exist")
	}

	q.kycRecoverySignerRole = signerRole
	return *signerRole, nil
}

func (q *SystemSettingsQ) LicenseAdminSignerRole() (uint64, error) {
	signerRole, err := keyValueUint64(q.Horizon, q.KeyValue.LicenseAdminSignerRole)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get role")
	}

	if signerRole == nil {
		return 0, errors.New("role id does not exist")
	}

	return *signerRole, nil
}

func (q *SystemSettingsQ) KYCRecoveryEnabled() (bool, error) {
	value, err := keyValueUint32(q.Horizon, q.KeyValue.KYCRecoveryEnabled)
	if err != nil {
		return false, errors.Wrap(err, "failed to get value")
	}

	if value == nil {
		return false, errors.New("kyc recovery enabled value does not exist")
	}

	return *value != 0, nil
}

func (q *SystemSettingsQ) DefaultAccountRole() (uint64, error) {
	accountRole, err := keyValueUint32(q.Horizon, q.KeyValue.DefaultAccountRole)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get role")
	}

	if accountRole == nil {
		return 0, errors.New("role id does not exist")
	}
	return uint64(*accountRole), nil
}

func (q *SystemSettingsQ) DefaultSignerRole() (uint64, error) {
	signerRole, err := keyValueUint32(q.Horizon, q.KeyValue.DefaultSignerRole)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get signer role")
	}

	if signerRole == nil {
		return 0, errors.New("signer role does not exist")
	}
	return uint64(*signerRole), nil
}

func keyValueUint32(h *horizon.Connector, key string) (*uint32, error) {
	resp, err := h.Client().Get(fmt.Sprintf("/v3/key_values/%s", key))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get key value")
	}

	if resp == nil {
		return nil, nil
	}

	var result regources.KeyValueEntryResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal key value")
	}
	if t := result.Data.Attributes.Value.Type; t != xdr.KeyValueEntryTypeUint32 {
		return nil, fmt.Errorf("invalid key value type: %v", t)
	}

	return result.Data.Attributes.Value.U32, nil
}

func keyValueUint64(h *horizon.Connector, key string) (*uint64, error) {
	resp, err := h.Client().Get(fmt.Sprintf("/v3/key_values/%s", key))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get key value")
	}

	if resp == nil {
		return nil, nil
	}

	var result regources.KeyValueEntryResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal key value")
	}
	if t := result.Data.Attributes.Value.Type; t != xdr.KeyValueEntryTypeUint64 {
		return nil, fmt.Errorf("invalid key value type: %v", t)
	}

	return result.Data.Attributes.Value.U64, nil
}
