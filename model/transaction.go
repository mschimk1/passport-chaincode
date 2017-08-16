package model

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"
)

// TransactionObjectType blockchain object type
const TransactionObjectType = "Transaction"

// TxDetails struct stores details of a transaction
type TxDetails struct {
	CustomerID   string            `json:"customer_id"`
	AccountID    string            `json:"account_id"`
	Amount       int64             `json:"amount"` // amount in cents
	Fee          int64             `json:"fee"`
	CurrencyCode string            `json:"currency"`
	Created      int64             `json:"created"` // unix time
	Description  string            `json:"description"`
	Params       map[string]string `json:"params,omitempty"`
}

// TxFailureCode stores allowed values for transaction failures
// Allowed values are "insufficient_funds", "account_closed"
type TxFailureCode string

// TxStatus stores allowed values for a transaction's status.
// Allowed values are "debited", "credited", "failed"
type TxStatus string

const (
	// TxFailureCodeNone successful transaction
	TxFailureCodeNone TxFailureCode = ""
	// InsufficientFunds transaction failure code
	InsufficientFunds TxFailureCode = "insufficient_funds"
	// AccountClosed transaction faiure code
	AccountClosed TxFailureCode = "account_closed"
	// Debited transaction status
	Debited TxStatus = "debited"
	// Credited transaction status
	Credited TxStatus = "credited"
	// Failed transaction status
	Failed TxStatus = "failed"
)

// Transaction data struct represents a money transfer (payer and payee sides)
type Transaction struct {
	Entity
	ID string `json:"id"`
	TxDetails
	FailureCode TxFailureCode `json:"failure_code,omitempty"`
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

// CreateTransaction a factory function for creating new Transaction entities
func CreateTransaction(customerID string, accountID string, t *Transfer, code TxFailureCode, status TxStatus) (*Transaction, error) {
	txn := &Transaction{Entity: Entity{TransactionObjectType}, FailureCode: code, Status: status}
	txn.TxDetails = TxDetails{
		CustomerID:   customerID,
		AccountID:    accountID,
		Created:      time.Now().Unix(),
		Amount:       t.Amount,
		Fee:          t.Fee,
		CurrencyCode: t.CurrencyCode,
		Description:  t.Description,
		Params:       t.Params,
	}
	transferData, _ := json.Marshal(txn)
	txn.ID = fmt.Sprintf("%x", newID(transferData))
	return txn, nil
}

func newID(data []byte) []byte {
	md5 := md5.New()
	md5.Write(data)
	return md5.Sum(nil)
}

// TransactionList stores a list of transactions
type TransactionList struct {
	Transactions []*Transaction `json:"transactions"`
}

// ByCreated sorts a list of transaction by creation timestamp
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
