package model

import (
	"encoding/json"
	"errors"
	"passport-chaincode/utils"
	"time"
)

// AccountObjectType blockchain object type
const AccountObjectType = "Account"

// Account struct holds information about a bank account
type Account struct {
	Entity
	ID            string            `json:"id"`
	CustomerID    string            `json:"customer_id"`
	BankName      string            `json:"bank_name"`
	AccountHolder string            `json:"account_holder"`
	Description   string            `json:"description"`
	CountryCode   string            `json:"country"`
	CurrencyCode  string            `json:"currency"`
	Created       int64             `json:"created"` // unix timestamp
	Balance       int64             `json:"balance"` // account balance in cents
	Default       bool              `json:"default_account"`
	Closed        bool              `json:"closed"`
	Params        map[string]string `json:"params,omitempty"` // additional name / value pairs
}

// AccountList holds a list of bank accounts
type AccountList struct {
	Accounts []*Account `json:"accounts"`
}

// GetObjectType returns the blockchain table name
func (a *Account) GetObjectType() string {
	return a.ObjectType
}

// UnmarshalJSON custom unmarshalling handles time conversion
func (a *Account) UnmarshalJSON(data []byte) error {
	type AccountData Account
	wrapper := &struct {
		Created string `json:"created"`
		*AccountData
	}{
		AccountData: (*AccountData)(a),
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return err
	}
	if wrapper.Created != "" {
		t1, err := time.Parse(time.RFC3339, wrapper.Created)
		if err != nil {
			return err
		}
		a.Created = t1.Unix()
	}
	return nil
}

// MarshalJSON custom marshalling handles time conversion
func (a *Account) MarshalJSON() ([]byte, error) {
	type AccountData Account
	return json.Marshal(&struct {
		Created string `json:"created"`
		*AccountData
	}{
		Created:     time.Unix(a.Created, 0).Format(time.RFC3339),
		AccountData: (*AccountData)(a),
	})
}

// CreateAccount Factory function creates a new Account struct and returns a pointer to it
func CreateAccount(accountBytes []byte) (*Account, error) {
	account := new(Account)
	if err := json.Unmarshal(accountBytes, account); err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("Error unmarshalling account data")
	}
	account.ObjectType = AccountObjectType
	if account.CustomerID == "" {
		return nil, errors.New("Missing required customer_id")
	}
	if account.ID == "" { // generate hash
		account.ID = utils.GenerateID(8)
	}
	if account.Created == 0 {
		account.Created = time.Now().Unix()
	}
	return account, nil
}

// Debit - debit the account
func (a *Account) Debit(amount int64) {
	a.Balance -= amount
}

// Credit - credit the account
func (a *Account) Credit(amount int64) {
	a.Balance += amount
}
