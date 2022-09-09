package data

import "time"

// HorizonInfo - Auto generated struct with https://mholt.github.io/json-to-go/
type HorizonInfo struct {
	LedgersState struct {
		Core struct {
			Latest                 int       `json:"latest"`
			OldestOnStart          int       `json:"oldest_on_start"`
			LastLedgerIncreaseTime time.Time `json:"last_ledger_increase_time"`
		} `json:"core"`
		History struct {
			Latest                 int       `json:"latest"`
			OldestOnStart          int       `json:"oldest_on_start"`
			LastLedgerIncreaseTime time.Time `json:"last_ledger_increase_time"`
		} `json:"history"`
		History2 struct {
			Latest                 int       `json:"latest"`
			OldestOnStart          int       `json:"oldest_on_start"`
			LastLedgerIncreaseTime time.Time `json:"last_ledger_increase_time"`
		} `json:"history_2"`
	} `json:"ledgers_state"`
	NetworkPassphrase  string `json:"network_passphrase"`
	AdminAccountID     string `json:"admin_account_id"`
	MasterExchangeName string `json:"master_exchange_name"`
	TxExpirationPeriod int    `json:"tx_expiration_period"`
	CurrentTime        int    `json:"current_time"`
	Precision          int    `json:"precision"`
	XdrRevision        string `json:"xdr_revision"`
	HorizonRevision    string `json:"horizon_revision"`
	MasterAccountID    string `json:"master_account_id"`
	EnvironmentName    string `json:"environment_name"`
}
