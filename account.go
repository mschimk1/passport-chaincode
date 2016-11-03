package main

import (
	"encoding/json"
	"time"
)

// Country - country codes
type Country string

// Currency - currency codes
type Currency string

// Account struct holds information about a bank account
type Account struct {
	ID            string            `json:"id"`
	CustomerID    string            `json:"customer_id"`
	BankName      string            `json:"bank_name"`
	AccountHolder string            `json:"account_holder"`
	Description   string            `json:"description"`
	Country       Country           `json:"country"`
	Currency      Currency          `json:"currency"`
	Created       int64             `json:"created"` // unix time
	Balance       int64             `json:"balance"` // account balance in cents
	Default       bool              `json:"default_account"`
	Deleted       bool              `json:"deleted"`
	Params        map[string]string `json:"params,omitempty"` // additional name / value pairs
}

//UnmarshalJSON custom unmarshalling handles time conversion
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

//MarshalJSON custom marshalling handles time conversion
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

// Debit -
func (a *Account) Debit(amount int64) {
	a.Balance -= amount
}

// Credit -
func (a *Account) Credit(amount int64) {
	a.Balance += amount
}

// AccountList holds a list of bank accounts
type AccountList struct {
	Accounts []*Account `json:"accounts"`
}
