/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateAccount struct {
	Key
	Relationships CreateAccountRelationships `json:"relationships"`
}
type CreateAccountResponse struct {
	Data     CreateAccount `json:"data"`
	Included Included      `json:"included"`
}

type CreateAccountListResponse struct {
	Data     []CreateAccount `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
}

// MustCreateAccount - returns CreateAccount from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateAccount(key Key) *CreateAccount {
	var createAccount CreateAccount
	if c.tryFindEntry(key, &createAccount) {
		return &createAccount
	}
	return nil
}
