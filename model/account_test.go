package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/stretchr/testify/suite"
)

type AccountSuite struct {
	suite.Suite
	testAccount *Account
}

func (suite *AccountSuite) SetupTest() {
	ts := time.Now().Unix()
	suite.testAccount = &Account{Entity{"Account"}, "1234", "1", "Test Bank", "John Smith", "", "AU", "AUD", ts, 1000, true, false, map[string]string(nil)}
}

func (suite *AccountSuite) TestGetObjectType() {
	suite.Equal(suite.testAccount.GetObjectType(), AccountObjectType)
}

func (suite *AccountSuite) TestUnmarshalAccountCreated() {
	timeStr := "2016-10-28T00:00:00+11:00"
	testTime, _ := time.Parse(time.RFC3339, timeStr)
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "{\"created\":\"%s\"}", timeStr)
	a := new(Account)
	err := json.Unmarshal(b.Bytes(), a)
	suite.Nil(err)
	suite.Equal(testTime.Unix(), a.Created)
}

func (suite *AccountSuite) TestMarshalAccountCreated() {
	timeStr := "2016-10-28T00:00:00+11:00"
	testTime, _ := time.Parse(time.RFC3339, timeStr)
	a := &Account{
		Created: testTime.Unix(),
	}
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "{\"created\":\"%s\"}", timeStr)
	accountBytes, err := json.Marshal(a)
	suite.Nil(err)
	matched, _ := regexp.Match(regexp.QuoteMeta(timeStr), accountBytes)
	suite.True(matched)
}

func (suite *AccountSuite) TestCreateAccountHappyPath() {
	accountData := []byte(`{"docType":"Account","id":"1234","customer_id":"1","bank_name":"Test Bank","account_holder":"John Smith","country":"AU","currency":"AUD","balance":1000,"default_account":true,"closed":false}`)
	a, err := CreateAccount(accountData)
	suite.Nil(err)
	suite.Equal(suite.testAccount, a)
}

func (suite *AccountSuite) TestCreateAccountMissingCustomerID() {
	accountData := "{\"bank_name\":\"Test Bank\", \"account_holder\": \"Mike\", \"country\": \"AU\", \"currency\": \"AUD\", \"balance\":100}"
	errMsg := "Missing required customer_id"
	_, err := CreateAccount([]byte(accountData))
	suite.Equal(errMsg, err.Error())
}

func (suite *AccountSuite) TestCreateAccountWithoutID() {
	accountData := []byte(`{"customer_id":"1","bank_name":"Test Bank","account_holder":"Mike","country":"AU","currency":"AUD","balance":100}`)
	acc, _ := CreateAccount(accountData)
	suite.Equal(8, len(acc.ID))
}

func (suite *AccountSuite) TestCreateAccountAddsCreated() {
	accountData := []byte(`{"customer_id":"1","bank_name":"Test Bank","account_holder":"Mike","country":"AU","currency":"AUD","balance":100}`)
	acc, _ := CreateAccount(accountData)
	var valid = regexp.MustCompile(`^[0-9]+$`)
	matched := valid.MatchString(strconv.FormatInt(acc.Created, 10))
	suite.True(matched)
}

func (suite *AccountSuite) TestDebit() {
	var amount, expected int64
	amount = 100
	expected = 900
	a := &Account{
		Balance: 1000,
	}
	a.Debit(amount)
	suite.Equal(expected, a.Balance)
}

func (suite *AccountSuite) TestCredit() {
	var amount, expected int64
	amount = 100
	expected = 1100
	a := &Account{
		Balance: 1000,
	}
	a.Credit(amount)
	suite.Equal(expected, a.Balance)
}
