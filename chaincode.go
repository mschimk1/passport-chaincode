/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
This chaincode provides a very simplistic shared ledger view of cross border
financial transactions. Its main purpose is to experiment with the hyperledger
fabric blockchain service on IBM Bluemix.
*/
package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim" // v0.5
)

const (
	// AccountTable is the account ledger table name
	AccountTable = "AccountTable"
	// AccountLookupTable is used for account list range queries by customer ID
	AccountLookupTable = "AccountLookupTable"
	// TransactionTable is the transaction ledger table name
	TransactionTable = "TransactionTable"
)

var (
	// passport chaincode application logger
	logger = shim.NewLogger("passport-chaincode")
	// mapping of chaincode handler functions
	handlerMap = NewHandlerMap()
	// Ledger table names - Init method recreates these tables every time a chaincode deploy is invoked
	dataTables = []string{AccountLookupTable, AccountTable, TransactionTable}
	// mapping of table name to key columns
	keyColumnDefinitions = map[string][]string{
		AccountLookupTable: []string{"CustomerID", "ID"},
		AccountTable:       []string{"ID"},
		TransactionTable:   []string{"AccountID", "ID"},
	}
)

func main() {
	initLogging()
	logger.Infof("Starting passport chaincode for IBM Bluemix Blockchain service v0.4.3")
	cc := new(Chaincode)
	cc.registerHandlers()
	err := shim.Start(cc)
	if err != nil {
		logger.Errorf("Error starting passport chaincode: %s", err)
	}
}

// Chaincode Chaincode shim method receiver struct
type Chaincode struct{}

//------------------------
// Chaincode API functions
//------------------------

// Init called to initialize the chaincode
func (cc *Chaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	for _, name := range dataTables {
		if err := cc.deleteLedgerTable(stub, name); err != nil {
			logger.Errorf("Error deleting table %s. Error: %s", name, err)
			return nil, fmt.Errorf("Error deleting table %s", name)
		}
		if err := cc.createLedgerTable(stub, name, keyColumnDefinitions[name]); err != nil {
			logger.Errorf("Error creating table %s. Error: %s", name, err)
			return nil, fmt.Errorf("Error creating table %s", name)
		}
	}
	return nil, nil
}

// Invoke chaincode interface implementation
func (cc *Chaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return cc.handleInvocation(stub, function, args)
}

// Query chaincode interface implementation
func (cc *Chaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return cc.handleInvocation(stub, function, args)
}

func (cc *Chaincode) handleInvocation(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	logger.Debugf("Invoking chaincode handler function %s with args %v", function, args)

	res, err := handlerMap.Handle(stub, function, args)
	if err != nil {
		logger.Errorf("Error when calling handler for function %s. Error: %s", function, err)
	}
	return res, err
}

//------------------
// Handler functions
//------------------

// GetAccountList query blockchain accounts by customer ID
func (cc *Chaincode) GetAccountList(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	logger.Debugf("Entering with args %v", args)

	rows, err := cc.queryLedgerRows(stub, AccountLookupTable, args)
	if err != nil {
		return nil, fmt.Errorf("Get account list failed. Error: %s", err)
	}
	var accountIDs []string
	accountList := AccountList{}

	for i := 0; i < len(rows); i++ {
		// column 1 in account lookup table
		id := rows[i].Columns[1].GetString_()
		accountIDs = append(accountIDs, id)
	}

	for _, id := range accountIDs {
		accountData, err := cc.GetAccount(stub, []string{id})
		if err != nil {
			return nil, err
		}
		a := new(Account)
		bytesToStruct(accountData, a)
		accountList.Accounts = append(accountList.Accounts, a)
	}
	jsonList, _ := json.Marshal(accountList)

	logger.Debugf("Returning account list: %s", jsonList)
	return jsonList, nil
}

// GetAccount query blockchain account by account ID
func (cc *Chaincode) GetAccount(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	logger.Debugf("Entering with args %v", args)

	data, err := cc.queryLedger(stub, AccountTable, args)
	if err != nil {
		logger.Errorf("Failed to get account details. Error: %s", err)
		return nil, err
	}
	logger.Debugf("Returning account details: %s", data)
	return data, nil
}

// OpenAccount opens an account
func (cc *Chaincode) OpenAccount(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	logger.Debugf("Entering with args %v", args)
	if len(args) == 0 {
		return nil, errors.New("Missing required account data JSON string")
	}
	account, err := createAccount([]byte(args[0]))
	if err != nil {
		logger.Errorf("Error when creating new account. Error: %s", err)
		return nil, fmt.Errorf("Error creating new account. Error: %s", err)
	}
	accountData, _ := json.Marshal(account)
	if err := cc.updateLedger(stub, AccountTable, []string{account.ID}, accountData); err != nil {
		logger.Errorf("Error when updating updating account ledger table. Error: %s", err)
		return nil, fmt.Errorf("Error updating account ledger table. Account ID already exists.")
	}

	if err := cc.updateLedger(stub, AccountLookupTable, []string{account.CustomerID, account.ID}, []byte{}); err != nil {
		logger.Errorf("Error when updating updating account lookup ledger table. Error: %s", err)
		return nil, err
	}

	return accountData, nil
}

// CloseAccount closes the given account
func (cc *Chaincode) CloseAccount(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	logger.Debugf("Entering with args %v", args)
	accountData, err := cc.GetAccount(stub, args)
	if err != nil {
		return nil, err
	}
	if accountData == nil {
		return nil, fmt.Errorf("Account with number %s not found.", args[0])
	}
	account := new(Account)
	bytesToStruct(accountData, account)
	account.Deleted = true
	accountData, _ = json.Marshal(account)
	if err := cc.replaceLedgerRow(stub, AccountTable, []string{account.ID}, accountData); err != nil {
		logger.Errorf("Error when updating account ledger table for account number %s. Error: %s", args[0], err)
		return nil, err
	}
	return nil, nil
}

// TransferMoney transfer money
func (cc *Chaincode) TransferMoney(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	logger.Debugf("Entering with args %v", args)

	transferData := args[0]
	t := new(Transfer)
	bytesToStruct([]byte(transferData), t)
	if err := t.validate(); err != nil {
		return nil, err
	}
	accountData, err := cc.GetAccount(stub, []string{t.FromAccountID})
	if err != nil {
		return nil, err
	}
	fromAccount := new(Account)
	bytesToStruct(accountData, fromAccount)
	accountData, err = cc.GetAccount(stub, []string{t.ToAccountID})
	if err != nil {
		return nil, err
	}
	toAccount := new(Account)
	bytesToStruct(accountData, toAccount)

	txnRef := generateHash(transferData)

	if fromAccount.Deleted || toAccount.Deleted {
		cc.recordTransaction(stub, fromAccount.ID, txnRef, t, AccountClosed, Failed)
		cc.recordTransaction(stub, toAccount.ID, txnRef, t, AccountClosed, Failed)
		return nil, nil
	}

	if fromAccount.Balance-t.Amount <= 0 {
		cc.recordTransaction(stub, fromAccount.ID, txnRef, t, InsufficientFunds, Failed)
		cc.recordTransaction(stub, toAccount.ID, txnRef, t, InsufficientFunds, Failed)
		return nil, nil
	}

	cc.recordTransaction(stub, fromAccount.ID, txnRef, t, "", Debited)
	cc.recordTransaction(stub, toAccount.ID, txnRef, t, "", Completed)

	cc.debitAccount(stub, fromAccount, t.Amount+t.Fee)
	cc.creditAccount(stub, toAccount, t.Amount)

	return nil, nil
}

// GetTransactionList query blockchain accounts by account ID
func (cc *Chaincode) GetTransactionList(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	logger.Debugf("Entering with args %v", args)

	dataCol := len(keyColumnDefinitions[TransactionTable])
	rows, err := cc.queryLedgerRows(stub, TransactionTable, args)
	if err != nil {
		return nil, fmt.Errorf("Get transaction list failed. Error: %s", err)
	}
	transactionList := TransactionList{}
	for _, row := range rows {
		t := new(Transaction)
		bytesToStruct(row.Columns[dataCol].GetBytes(), t)
		if t == nil {
			logger.Errorf("Error unmarshalling transaction data")
			continue
		}
		transactionList.Transactions = append(transactionList.Transactions, t)
	}
	sort.Sort(sort.Reverse(ByCreated(transactionList.Transactions)))
	jsonList, _ := json.Marshal(transactionList)

	logger.Debugf("Returning transaction list: %s", jsonList)
	return jsonList, nil
}

// GetTransaction query blockchain transaction by transaction ID
func (cc *Chaincode) GetTransaction(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	logger.Debugf("Entering with args %v", args)

	txnData, err := cc.queryLedger(stub, TransactionTable, args)
	if err != nil {
		logger.Errorf("Failed to get transaction details. Error: %s", err)
		return nil, err
	}
	logger.Debugf("Returning transaction details: %s", txnData)
	return txnData, nil
}

// TopupAccount update account balance
func (cc *Chaincode) TopupAccount(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	logger.Debugf("Entering with args %v", args)

	accountData, err := cc.GetAccount(stub, []string{args[0]})
	if err != nil {
		return nil, err
	}
	account := new(Account)
	bytesToStruct([]byte(accountData), account)
	if account == nil {
		return nil, errors.New("Error unmarshalling account data")
	}
	amount, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error parsing amount value %s", args[1])
	}
	account.Credit(amount)
	accountData, err = json.Marshal(account)
	if err != nil {
		return nil, errors.New("Error marshalling account data")
	}
	if err := cc.replaceLedgerRow(stub, AccountTable, []string{account.ID}, accountData); err != nil {
		logger.Errorf("Error when updating account ledger table for account number %s. Error: %s", args[0], err)
		return nil, err
	}
	return nil, nil
}

// Creates a new Account struct and returns a pointer to it
func createAccount(accountData []byte) (*Account, error) {

	account := new(Account)
	if err := bytesToStruct([]byte(accountData), account); err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("Error unmarshalling account data")
	}
	if account.CustomerID == "" {
		return nil, errors.New("Missing required customer_id")
	}
	if account.ID == "" { // generate hash
		account.ID = generateID(8)
	}
	if account.Created == 0 {
		account.Created = time.Now().Unix()
	}
	if account.Balance <= 0 {
		logger.Warningf("Initial account balance is %d", account.Balance)
	}
	return account, nil
}

func (cc *Chaincode) recordTransaction(stub *shim.ChaincodeStub, accountID string, id string, t *Transfer, code TxFailureCode, status TxStatus) error {

	txn := &Transaction{ID: id, FailureCode: code, Status: status}
	txn.TxDetails = TxDetails{
		AccountID:   accountID,
		Created:     time.Now().Unix(),
		Amount:      t.Amount,
		Fee:         t.Fee,
		Currency:    t.Currency,
		Description: t.Description,
		Params:      t.Params,
	}
	txnData, err := json.Marshal(txn)
	if err != nil {
		return fmt.Errorf("Error marshalling transaction data. Error: %s", err)
	}
	if err := cc.updateLedger(stub, TransactionTable, []string{accountID, id}, txnData); err != nil {
		logger.Errorf("Error when updating updating transaction ledger table. Error: %s", err)
		return fmt.Errorf("Error updating transaction ledger table. Transaction ID already exists.")
	}
	return nil
}

func (cc *Chaincode) debitAccount(stub *shim.ChaincodeStub, a *Account, amount int64) error {
	a.Debit(amount)
	buff, _ := json.Marshal(a)
	err := cc.replaceLedgerRow(stub, AccountTable, []string{a.ID}, buff)
	if err != nil {
		logger.Errorf("Error when updating updating account ledger table. Error: %s", err)
		return fmt.Errorf("{\"Error\":\"%s\"}", err)
	}
	return nil
}

func (cc *Chaincode) creditAccount(stub *shim.ChaincodeStub, a *Account, amount int64) error {
	a.Credit(amount)
	buff, _ := json.Marshal(a)
	err := cc.replaceLedgerRow(stub, AccountTable, []string{a.ID}, buff)
	if err != nil {
		logger.Errorf("Error when updating updating account ledger table. Error: %s", err)
		return fmt.Errorf("{\"Error\":\"%s\"}", err)
	}
	return nil
}

//-------------------------------------------------
// Helpers
//-------------------------------------------------
func initLogging() {
	logger.SetLevel(shim.LogInfo)
	logLevel, _ := shim.LogLevel(os.Getenv("SHIM_LOGGING_LEVEL"))
	shim.SetLoggingLevel(logLevel)
}

// Registers handler function mappings
func (cc *Chaincode) registerHandlers() {
	handlerMap.Add("OpenAccount", cc.OpenAccount)
	handlerMap.Add("CloseAccount", cc.CloseAccount)
	handlerMap.Add("GetAccount", cc.GetAccount)
	handlerMap.Add("GetAccountList", cc.GetAccountList)
	handlerMap.Add("TransferMoney", cc.TransferMoney)
	handlerMap.Add("TopupAccount", cc.TopupAccount)
	handlerMap.Add("GetTransaction", cc.GetTransaction)
	handlerMap.Add("GetTransactionList", cc.GetTransactionList)
}

func (cc *Chaincode) queryLedger(stub *shim.ChaincodeStub, name string, keys []string) ([]byte, error) {
	logger.Debugf("Entering with table name %s and key args %v", name, keys)

	var columns []shim.Column
	nCols := min(len(keyColumnDefinitions[name]), len(keys))
	for i := 0; i < nCols; i++ {
		columns = append(columns, *createKeyColumn(keys[i]))
	}
	row, err := stub.GetRow(name, columns)
	if err != nil {
		return nil, err
	}
	if len(row.Columns) == 0 {
		logger.Debugf("No data available for table name %s and key args %v", name, keys)
		return nil, nil
	}
	logger.Debugf("Fetched %d columns from table %s", len(row.Columns), name)
	// last column stores the data object
	result := row.Columns[nCols].GetBytes()
	logger.Debugf("Query ledger result %s", result)
	return result, nil
}

func (cc *Chaincode) queryLedgerRows(stub *shim.ChaincodeStub, name string, keys []string) ([]shim.Row, error) {
	logger.Debugf("Entering with table name %s and key args %v", name, keys)

	var columns []shim.Column
	for _, key := range keys {
		columns = append(columns, *createKeyColumn(key))
	}

	rowChan, err := stub.GetRows(name, columns)
	if err != nil {
		return nil, err
	}
	var rows []shim.Row
	for rowChan != nil {
		select {
		case row, ok := <-rowChan:
			if !ok {
				rowChan = nil
			} else {
				rows = append(rows, row)
			}
		}
	}
	logger.Debugf("Fetched %d rows from ledger table %s", len(rows), name)
	return rows, nil
}

func (cc *Chaincode) createLedgerTable(stub *shim.ChaincodeStub, name string, keys []string) error {
	logger.Debugf("Entering with table name %s and key args %v", name, keys)

	var columnDefs []*shim.ColumnDefinition
	var colDef shim.ColumnDefinition
	for i := range keys {
		colDef := shim.ColumnDefinition{Name: "Key" + strconv.Itoa(i), Type: shim.ColumnDefinition_STRING, Key: true}
		columnDefs = append(columnDefs, &colDef)
	}
	colDef = shim.ColumnDefinition{Name: name + "TableData", Type: shim.ColumnDefinition_BYTES, Key: false}
	columnDefs = append(columnDefs, &colDef)

	logger.Debugf("Creating table %s with column spec %v", name, keys)
	if err := stub.CreateTable(name, columnDefs); err != nil {
		return err
	}
	return nil
}

func (cc *Chaincode) deleteLedgerTable(stub *shim.ChaincodeStub, name string) error {
	logger.Debugf("Deleting ledger table %s", name)
	return stub.DeleteTable(name)
}

func (cc *Chaincode) updateLedger(stub *shim.ChaincodeStub, name string, keys []string, args []byte) error {
	logger.Debugf("Entering with table name %s and key args %v", name, keys)
	columns := createColumnSpec(keys, args)
	row := shim.Row{Columns: columns}
	ok, err := stub.InsertRow(name, row)
	if err != nil {
		return err
	}
	if !ok {
		logger.Errorf("Row with given key already exists")
		return errors.New("Row with given key already exists")
	}
	logger.Debugf("Insert row into table %s successful.", name)

	return nil
}

func (cc *Chaincode) replaceLedgerRow(stub *shim.ChaincodeStub, name string, keys []string, args []byte) error {
	logger.Debugf("Entering with table name %s and key args %v", name, keys)

	columns := createColumnSpec(keys, args)
	row := shim.Row{Columns: columns}
	ok, err := stub.ReplaceRow(name, row)
	if err != nil {
		logger.Errorf("Error replacing ledger table table row. Error: %s", err)
		return err
	}
	if !ok {
		return errors.New("Error replacing ledger table table row.")
	}
	logger.Debugf("Successfully replaced row in ledger table %s.", name)
	return nil
}

// Deletes a row from the given ledger table
func (cc *Chaincode) deleteFromLedger(stub *shim.ChaincodeStub, name string, keys []string) error {
	logger.Debugf("Entering with table name %s and key args %v", name, keys)

	var columns []shim.Column
	for _, key := range keys {
		columns = append(columns, *createKeyColumn(key))
	}
	if err := stub.DeleteRow(name, columns); err != nil {
		return err
	}
	logger.Debugf("Successfully deleted row from ledger table %s.", name)
	return nil
}

func createKeyColumn(key string) *shim.Column {
	return &shim.Column{Value: &shim.Column_String_{String_: key}}
}

// Creates key column definitions - only supports string keys for now
func createKeyColumnSpec(keys []string) []*shim.Column {
	var columns []*shim.Column
	for _, key := range keys {
		columns = append(columns, createKeyColumn(key))
	}
	return columns
}

// Creates column definitions - only supports string keys for now
func createColumnSpec(keys []string, data []byte) []*shim.Column {
	columns := createKeyColumnSpec(keys)
	if data != nil {
		colBytes := shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(data)}}
		columns = append(columns, &colBytes)
	}
	return columns
}

// Unmarshals byte slice into given data type
func bytesToStruct(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		logger.Errorf("Error unmarshalling data for type %T", v)
		return err
	}
	return nil
}

// Generates a sha256 hash of the given message
func generateHash(msg string) string {
	msgHash := sha256.Sum256([]byte(msg))
	return fmt.Sprintf("%x", msgHash)
}

// Generates a fixed length string of random digits
func generateID(length int) string {
	rand.Seed(time.Now().UnixNano())
	r := []rune("1234567890")
	b := make([]rune, length)
	for i := range b {
		b[i] = r[rand.Intn(len(r))]
	}
	return string(b)
}

// min helper for int types
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
