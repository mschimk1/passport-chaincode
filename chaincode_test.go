package main

import (
	"encoding/json"
	"fmt"
	"passport-chaincode/model"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(ChaincodeSuite))
	suite.Run(t, new(HandlerSuite))
}

type ChaincodeSuite struct {
	suite.Suite
	cc   *Chaincode
	stub *shim.MockStub
}

func (suite *ChaincodeSuite) SetupTest() {
	suite.cc = new(Chaincode)
	suite.cc.registerHandlers()
	suite.stub = shim.NewMockStub("mockStub", suite.cc)
}

func (suite *ChaincodeSuite) checkState(name string, value string) {
	bytes := suite.stub.State[name]
	suite.NotNil(bytes, "Failed to get state "+name)
	suite.Equal(value, string(bytes), "State value "+name+" was not as expected")
}

func (suite *ChaincodeSuite) checkStateValue(key string, name string, value string) {
	bytes := suite.stub.State[key]
	suite.NotNil(bytes, "Failed to get state for key: "+key)

	suite.Equal(value, string(bytes), "State value "+name+" was not as expected")
}

func (suite *ChaincodeSuite) checkQuery(name string, value string) {
	bytes, err := suite.stub.MockQuery("query", []string{name})
	suite.Nil(err, "Query failed")
	suite.NotNil(bytes, "Failed to get value")
	suite.Equal(value, string(bytes), "Query value "+name+"was not as expected")
}

func (suite *ChaincodeSuite) checkInvoke(function string, args []string) {
	_, err := suite.stub.MockInvoke("t1234", function, args)
	suite.Nil(err, "Invoke failed")
}

func (suite *ChaincodeSuite) TestOpenAccountValidation() {
	_, err := suite.stub.MockInvoke("t1234", "OpenAccount", []string{})
	suite.Equal(err.Error(), "Missing required account data JSON")
}

func (suite *ChaincodeSuite) TestOpenAccount() {
	testAccount := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"1234","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`
	_, err := suite.stub.MockInvoke("t1234", "OpenAccount", []string{testAccount})
	suite.Nil(err)
	suite.checkState("Account01012340", testAccount)
}

func (suite *ChaincodeSuite) TestGetAccountListValidation() {
	_, err := suite.stub.MockInvoke("t1234", "GetAccountList", []string{})
	suite.Equal(err.Error(), "Missing required customer ID")
}

func (suite *ChaincodeSuite) TestGetAccountListSingle() {
	testAccount1 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"1234","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`
	testAccountList := `{"accounts":[` + testAccount1 + "]}"
	suite.stub.MockInvoke("t1", "OpenAccount", []string{testAccount1})
	accountList, err := suite.stub.MockInvoke("t2", "GetAccountList", []string{"1"})
	suite.Nil(err)
	suite.Equal(testAccountList, string(accountList))
}

func (suite *ChaincodeSuite) TestGetAccountList() {
	testAccount1 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"1234","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`
	testAccount2 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"5678","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`
	testAccountList := `{"accounts":[` + testAccount1 + "," + testAccount2 + "]}"
	suite.stub.MockInvoke("t1", "OpenAccount", []string{testAccount1})
	suite.stub.MockInvoke("t2", "OpenAccount", []string{testAccount2})
	accountList, err := suite.stub.MockInvoke("t3", "GetAccountList", []string{"1"})
	suite.Nil(err)
	suite.Equal(testAccountList, string(accountList))
}

func (suite *ChaincodeSuite) TestGetAccountValidation() {
	_, err := suite.stub.MockInvoke("t1234", "GetAccount", []string{})
	suite.Equal(err.Error(), "Missing required customer ID and / or account ID")
}

func (suite *ChaincodeSuite) TestGetAccount() {
	testAccount := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"1234","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`
	suite.stub.MockInvoke("t1", "OpenAccount", []string{testAccount})
	account, err := suite.stub.MockInvoke("t2", "GetAccount", []string{"1", "1234"})
	suite.Nil(err)
	suite.Equal(testAccount, string(account))
}

func (suite *ChaincodeSuite) TestCloseAccountValidation() {
	_, err := suite.stub.MockInvoke("t1234", "CloseAccount", []string{})
	suite.Equal(err.Error(), "Missing required customer ID and / or account ID")
}

func (suite *ChaincodeSuite) TestCloseAccountNonExistingAccount() {
	_, err := suite.stub.MockInvoke("t1234", "CloseAccount", []string{"1", "1234"})
	suite.Equal(err.Error(), "Account with number 1234 not found.")
}

func (suite *ChaincodeSuite) TestCloseAccount() {
	testAccount := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"1234","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`
	suite.stub.MockInvoke("t1", "OpenAccount", []string{testAccount})
	suite.stub.MockInvoke("t2", "CloseAccount", []string{"1", "1234"})
	account, err := suite.stub.MockInvoke("t2", "GetAccount", []string{"1", "1234"})
	suite.Nil(err)
	actual := new(model.Account)
	json.Unmarshal(account, actual)
	suite.True(actual.Closed)
}

func (suite *ChaincodeSuite) TestTopupAccount() {
	testAccount := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"1234","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`
	suite.stub.MockInvoke("t1", "OpenAccount", []string{testAccount})
	_, err := suite.stub.MockInvoke("t2", "TopupAccount", []string{"1", "1234", "1000"})
	account, err := suite.stub.MockInvoke("t3", "GetAccount", []string{"1", "1234"})
	suite.Nil(err)
	actual := new(model.Account)
	json.Unmarshal(account, actual)
	suite.Equal(int64(2000), actual.Balance)
}

func (suite *ChaincodeSuite) TestTransferMoneyValidation() {
	_, err := suite.stub.MockInvoke("t1234", "TransferMoney", []string{})
	suite.Equal(err.Error(), "Missing transfer details JSON")
}

func (suite *ChaincodeSuite) TestTransferMoneyHappyPath() {
	testAccount1 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"1234","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`
	testAccount2 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"5678","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`

	suite.stub.MockInvoke("t1", "OpenAccount", []string{testAccount1})
	suite.stub.MockInvoke("t2", "OpenAccount", []string{testAccount2})

	transfer := `{"from_customer": "1", "from_account": "1234", "to_customer": "1", "to_account":"5678", "currency":"AUD", "amount":1000}`

	_, err := suite.stub.MockInvoke("t3", "TransferMoney", []string{transfer})
	suite.Nil(err)

	account1, err := suite.stub.MockInvoke("t4", "GetAccount", []string{"1", "1234"})
	account2, err := suite.stub.MockInvoke("t4", "GetAccount", []string{"1", "5678"})

	a1 := new(model.Account)
	a2 := new(model.Account)
	json.Unmarshal(account1, a1)
	json.Unmarshal(account2, a2)
	suite.Equal(int64(0), a1.Balance)
	suite.Equal(int64(2000), a2.Balance)
}

func (suite *ChaincodeSuite) TestTransferMoneyInsufficientFunds() {
	testAccount1 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"1234","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":100,"default_account":true,"closed":false}`
	testAccount2 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"5678","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`

	suite.stub.MockInvoke("t1", "OpenAccount", []string{testAccount1})
	suite.stub.MockInvoke("t2", "OpenAccount", []string{testAccount2})

	transfer := `{"from_customer": "1", "from_account": "1234", "to_customer": "1", "to_account":"5678", "currency":"AUD", "amount":1000}`

	_, err := suite.stub.MockInvoke("t3", "TransferMoney", []string{transfer})
	suite.Equal("Insufficient funds available in account 1234", err.Error())
}

func (suite *ChaincodeSuite) TestTransferMoneyClosedFromAccount() {
	testAccount1 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"1234","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":100,"default_account":true,"closed":true}`
	testAccount2 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"5678","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`

	suite.stub.MockInvoke("t1", "OpenAccount", []string{testAccount1})
	suite.stub.MockInvoke("t2", "OpenAccount", []string{testAccount2})

	transfer := `{"from_customer": "1", "from_account": "1234", "to_customer": "1", "to_account":"5678", "currency":"AUD", "amount":1000}`

	_, err := suite.stub.MockInvoke("t2", "TransferMoney", []string{transfer})
	suite.Equal("Cannot transfer money from closed account 1234", err.Error())
}

func (suite *ChaincodeSuite) TestTransferMoneyClosedToAccount() {
	testAccount1 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"1234","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":100,"default_account":true,"closed":false}`
	testAccount2 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"5678","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":true}`

	suite.stub.MockInvoke("t1", "OpenAccount", []string{testAccount1})
	suite.stub.MockInvoke("t2", "OpenAccount", []string{testAccount2})

	transfer := `{"from_customer": "1", "from_account": "1234", "to_customer": "1", "to_account":"5678", "currency":"AUD", "amount":1000}`

	_, err := suite.stub.MockInvoke("t2", "TransferMoney", []string{transfer})
	suite.Equal("Cannot transfer money into closed account 5678", err.Error())
}

func (suite *ChaincodeSuite) TestGetTransactionListValidation() {
	_, err := suite.stub.MockInvoke("t1234", "GetTransactionList", []string{})
	suite.Equal(err.Error(), "Missing required customer ID and / or account ID")
}

func (suite *ChaincodeSuite) TestGetTransactionList() {
	testAccount1 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"1234","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`
	testAccount2 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"5678","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`

	suite.stub.MockInvoke("t1", "OpenAccount", []string{testAccount1})
	suite.stub.MockInvoke("t2", "OpenAccount", []string{testAccount2})

	transfer := `{"from_customer": "1", "from_account": "1234", "to_customer": "1", "to_account":"5678", "currency":"AUD", "amount":1000}`

	_, err := suite.stub.MockInvoke("t3", "TransferMoney", []string{transfer})
	suite.Nil(err)

	transactions, _ := suite.stub.MockInvoke("t4", "GetTransactionList", []string{"1", "1234"})
	txnList := new(model.TransactionList)
	json.Unmarshal(transactions, txnList)
	suite.Equal(txnList.Transactions[0].Status, model.Debited)

	transactions, _ = suite.stub.MockInvoke("t4", "GetTransactionList", []string{"1", "5678"})
	txnList = new(model.TransactionList)
	json.Unmarshal(transactions, txnList)
	suite.Equal(txnList.Transactions[0].Status, model.Credited)
}

func (suite *ChaincodeSuite) TestGetTransaction() {
	testAccount1 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"1234","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`
	testAccount2 := `{"created":"2017-08-15T00:00:00+10:00","docType":"Account","id":"5678","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","description":"","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`

	suite.stub.MockInvoke("t1", "OpenAccount", []string{testAccount1})
	suite.stub.MockInvoke("t2", "OpenAccount", []string{testAccount2})

	transfer := `{"from_customer": "1", "from_account": "1234", "to_customer": "1", "to_account":"5678", "currency":"AUD", "amount":1000}`

	_, err := suite.stub.MockInvoke("t3", "TransferMoney", []string{transfer})
	suite.Nil(err)

	transactions, _ := suite.stub.MockInvoke("t4", "GetTransactionList", []string{"1", "1234"})
	txnList := new(model.TransactionList)
	json.Unmarshal(transactions, txnList)

	tran, _ := suite.stub.MockInvoke("t4", "GetTransaction", []string{"1", "1234", txnList.Transactions[0].ID})
	suite.NotNil(tran)
	fmt.Println(string(tran))
}
