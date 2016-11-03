package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"
)

func TestUnmarshalAccount(t *testing.T) {
	timeStr := "2016-10-28T00:00:00+11:00"
	testTime, _ := time.Parse(time.RFC3339, timeStr)
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "{\"created\":\"%s\"}", timeStr)
	a := new(Account)
	err := json.Unmarshal(b.Bytes(), a)
	if err != nil {
		t.Errorf("Expected account struct, but got %s", err)
	}
	if a.Created != testTime.Unix() {
		t.Errorf("Expected unix time %d, but got %d", testTime.Unix(), a.Created)
	}
}

func TestMarshalAccount(t *testing.T) {
	timeStr := "2016-10-28T00:00:00+11:00"
	testTime, _ := time.Parse(time.RFC3339, timeStr)
	a := &Account{
		Created: testTime.Unix(),
	}
	accountBytes, err := json.Marshal(a)
	if err != nil {
		t.Errorf("Expected account byte slice, but got %s", err)
	}
	matched, _ := regexp.Match(regexp.QuoteMeta(timeStr), accountBytes)
	if !matched {
		t.Error("Expected regexp match on time string")
	}
}

func TestCreateAccountHappyPath(t *testing.T) {
	accountData := "{\"customer_id\":\"12345\", \"bank_name\":\"Test Bank\", \"account_holder\": \"Mike\", \"country\": \"AU\", \"currency\": \"AUD\", \"balance\":1000}"
	a, err := createAccount([]byte(accountData))
	if err != nil {
		t.Errorf("Unexpected error when creating account with valid data. Error: %s", err)
	}
	if a.ID == "" {
		t.Errorf("Expected generated account ID but got %s", a.ID)
	}
	if a.CustomerID != "12345" {
		t.Errorf("Expected customer ID to be 12345, but got %s", a.CustomerID)
	}
	if a.Balance != 1000 {
		t.Errorf("Expected initial account balance to be 100, but got %d", a.Balance)
	}
	if a.Created == 0 {
		t.Error("Expected account created to be greater zero")
	}

}

func TestCreateAccountMissingCustomerID(t *testing.T) {
	accountData := "{\"bank_name\":\"Test Bank\", \"account_holder\": \"Mike\", \"country\": \"AU\", \"currency\": \"AUD\", \"balance\":100}"
	errMsg := "Missing required customer_id"
	_, err := createAccount([]byte(accountData))
	if err.Error() != errMsg {
		t.Errorf("Expected error message to be %s", errMsg)
	}
}

func TestCreateAccountIncorrectBalanceType(t *testing.T) {
	accountData := "{\"bank_name\":\"Test Bank\", \"account_holder\": \"Mike\", \"country\": \"AU\", \"currency\": \"AUD\", \"balance\":\"100\"}"
	errMsg := "json: cannot unmarshal string into Go value of type int64"
	_, err := createAccount([]byte(accountData))
	if err.Error() != errMsg {
		t.Errorf("Expected error message to be %s, but was %s", errMsg, err)
	}
}

func TestDebit(t *testing.T) {
	var amount, expected int64
	amount = 100
	expected = 900
	a := &Account{
		Balance: 1000,
	}
	a.Debit(amount)
	if a.Balance != expected {
		t.Errorf("Expected account balance of %d after debit, but got %d", expected, a.Balance)
	}
}

func TestCredit(t *testing.T) {
	var amount, expected int64
	amount = 100
	expected = 1100
	a := &Account{
		Balance: 1000,
	}
	a.Credit(amount)
	if a.Balance != expected {
		t.Errorf("Expected account balance of %d after credit, but got %d", expected, a.Balance)
	}
}
