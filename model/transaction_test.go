package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/stretchr/testify/suite"
)

type TransactionSuite struct {
	suite.Suite
}

func (suite *TransactionSuite) TestGetObjectType() {
	txPtr := &Transaction{Entity: Entity{"Transaction"}}
	suite.Equal("Transaction", txPtr.GetObjectType())
}

func (suite *TransactionSuite) TestUnmarshalTransactionCreated() {
	timeStr := "2016-10-28T00:00:00+11:00"
	testTime, _ := time.Parse(time.RFC3339, timeStr)
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "{\"created\":\"%s\"}", timeStr)
	tx := new(Transaction)
	err := json.Unmarshal(b.Bytes(), tx)
	suite.Nil(err)
	suite.Equal(testTime.Unix(), tx.Created)
}

func (suite *TransactionSuite) TestMarshalTransactionCreated() {
	timeStr := "2016-10-28T00:00:00+11:00"
	testTime, _ := time.Parse(time.RFC3339, timeStr)
	details := TxDetails{Created: testTime.Unix()}
	tptr := &Transaction{TxDetails: details}
	txBytes, err := json.Marshal(tptr)
	suite.Nil(err)
	matched, _ := regexp.Match(regexp.QuoteMeta(timeStr), txBytes)
	suite.True(matched)
}

func (suite *TransactionSuite) TestCreateTransaction() {
	tPtr := &Transfer{"1", "1234", "2", "5678", 100, 0, "AUD", "", map[string]string(nil)}
	txn, _ := CreateTransaction("1", "1234", tPtr, "", Credited)
	suite.Equal(32, len(txn.ID))
}

func (suite *TransactionSuite) TestTransactionListSort() {
	now := time.Now()
	earlier := now.Add(-1 * time.Minute)
	t1 := Transaction{
		ID:        "1",
		TxDetails: TxDetails{Created: earlier.Unix()},
	}
	t2 := Transaction{
		ID:        "2",
		TxDetails: TxDetails{Created: now.Unix()},
	}
	transactionList := TransactionList{}
	transactionList.Transactions = append(transactionList.Transactions, &t1)
	transactionList.Transactions = append(transactionList.Transactions, &t2)
	sort.Sort(sort.Reverse(ByCreated(transactionList.Transactions)))

	suite.Equal("2", transactionList.Transactions[0].ID)
}
