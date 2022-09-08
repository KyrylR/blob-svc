package data

import (
	"blob-svc/internal/types"
	"time"
)

type Wallet struct {
	Id                int64
	AccountID         types.Address
	CurrentAccountID  types.Address
	WalletId          *string
	Username          string
	KeychainData      string
	VerificationToken string
	// Verified comes from join on email_tokens and shows if wallet email was confirmed
	Verified *bool
	// Referrer comes from join on referrals add shows who referred this wallet
	Referrer *string

	//LastSentAt comes from join on email_tokens  and shows when verified letter was send
	LastSentAt *time.Time

	LastTFACheck *time.Time
	RegisteredAt *time.Time

	// KYCBlobID comes from join on user_states
	KYCBlobID *string
	//KYCBlobValue this is addition value comes from blobs.Value for specific user
	KYCBlobValue *string
	// State comes from join on identities
	State *int
}

// IsVerified helper for handling nil attribute
func (w Wallet) IsVerified() bool {
	return w.Verified != nil && *w.Verified
}
