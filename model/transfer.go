package model

import (
	"errors"
	"fmt"
)

// Transfer struct contains information about a money transfer
type Transfer struct {
	FromCustomerID string            `json:"from_customer"`
	FromAccountID  string            `json:"from_account"`
	ToCustomerID   string            `json:"to_customer"`
	ToAccountID    string            `json:"to_account"`
	Amount         int64             `json:"amount"` // amount in cents
	Fee            int64             `json:"fee"`
	CurrencyCode   string            `json:"currency"`
	Description    string            `json:"description"`
	Params         map[string]string `json:"params,omitempty"`
}

// Validate - checks that required are present in the transfer object
func (t *Transfer) Validate() error {
	if t.FromCustomerID == "" {
		return errors.New("Missing required from_customer value")
	}
	if t.FromAccountID == "" {
		return errors.New("Missing required from_account value")
	}
	if t.ToCustomerID == "" {
		return errors.New("Missing required to_customer value")
	}
	if t.ToAccountID == "" {
		return errors.New("Missing required to_account value")
	}
	if t.Amount <= 0 {
		return fmt.Errorf("Invalid transfer amount %d", t.Amount)
	}
	if t.CurrencyCode == "" {
		return errors.New("Missing required currency value")
	}
	// TODO: check valid currency codes
	return nil
}
