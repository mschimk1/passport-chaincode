package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Transfer struct contains information about a money transfer
type Transfer struct {
	FromAccountID string            `json:"from_account"`
	ToAccountID   string            `json:"to_account"`
	Amount        int64             `json:"amount"` // amount in cents
	Fee           int64             `json:"fee"`
	Currency      Currency          `json:"currency"`
	Description   string            `json:"description"`
	Params        map[string]string `json:"params,omitempty"`
}

func (t *Transfer) validate() error {
	if t.FromAccountID == "" {
		return errors.New("Missing required from_account value")
	}
	if t.ToAccountID == "" {
		return errors.New("Missing required to_account value")
	}
	if t.Amount <= 0 {
		return fmt.Errorf("Invalid transfer amount %d", t.Amount)
	}
	if t.Currency == "" {
		return errors.New("Missing required currency value")
	}
	// TODO: check valid currency codes
	return nil
}

// TxDetails struct stores details of a transaction
type TxDetails struct {
	AccountID   string            `json:"account_id"`
	Amount      int64             `json:"amount"` // amount in cents
	Fee         int64             `json:"fee"`
	Currency    Currency          `json:"currency"`
	Created     int64             `json:"created"` // unix time
	Description string            `json:"description"`
	Params      map[string]string `json:"params,omitempty"`
}

// TxFailureCode stores allowed values for transaction failures
// Allowed values are "insufficient_funds", "account_closed"
type TxFailureCode string

// TxStatus stores allowed values for a transaction's status.
// Allowed values are "debited", "completed", "failed"
type TxStatus string

const (
	// InsufficientFunds transaction failure code
	InsufficientFunds TxFailureCode = "insufficient_funds"
	// AccountClosed transaction faiure code
	AccountClosed TxFailureCode = "account_closed"
	// Debited transaction status
	Debited TxStatus = "debited"
	// Completed transaction status
	Completed TxStatus = "completed"
	// Failed transaction status
	Failed TxStatus = "failed"
)

// Transaction data struct represents a money transfer (payer and payee sides)
type Transaction struct {
	ID string `json:"id"`
	TxDetails
	FailureCode TxFailureCode `json:"failure_code"`
	Status      TxStatus      `json:"status"`
}

//UnmarshalJSON custom unmarshalling handles time conversion
func (t *Transaction) UnmarshalJSON(data []byte) error {
	type TransactionData Transaction
	wrapper := &struct {
		Created string `json:"created"`
		*TransactionData
	}{
		TransactionData: (*TransactionData)(t),
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return err
	}
	t1, err := time.Parse(time.RFC3339, wrapper.Created)
	if err != nil {
		return err
	}
	t.Created = t1.Unix()
	return nil
}

//MarshalJSON custom marshalling handles time conversion
func (t *Transaction) MarshalJSON() ([]byte, error) {
	type TransactionData Transaction
	return json.Marshal(&struct {
		Created string `json:"created"`
		*TransactionData
	}{
		Created:         time.Unix(t.Created, 0).Format(time.RFC3339),
		TransactionData: (*TransactionData)(t),
	})
}

// TransactionList stores a list of transactions
type TransactionList struct {
	Transactions []*Transaction `json:"transactions"`
}

type ByCreated []*Transaction

func (t ByCreated) Len() int {
	return len(t)
}

func (t ByCreated) Less(i, j int) bool {
	return t[i].Created < t[j].Created
}

func (t ByCreated) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
