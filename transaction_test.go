package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"testing"
	"time"
)

func TestUnmarshalTransaction(t *testing.T) {
	timeStr := "2016-10-28T00:00:00+11:00"
	testTime, _ := time.Parse(time.RFC3339, timeStr)
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "{\"created\":\"%s\"}", timeStr)
	tx := new(Transaction)
	err := json.Unmarshal(b.Bytes(), tx)
	if err != nil {
		t.Errorf("Expected transaction struct, but got %s", err)
	}
	if tx.Created != testTime.Unix() {
		t.Errorf("Expected unix epoch %d, but got %d", testTime.Unix(), tx.Created)
	}
}

func TestMarshalTransaction(t *testing.T) {
	timeStr := "2016-10-28T00:00:00+11:00"
	testTime, _ := time.Parse(time.RFC3339, timeStr)
	details := TxDetails{Created: testTime.Unix()}
	tptr := &Transaction{TxDetails: details}
	txBytes, err := json.Marshal(tptr)
	if err != nil {
		t.Errorf("Expected account byte slice, but got %s", err)
	}
	matched, _ := regexp.Match(regexp.QuoteMeta(timeStr), txBytes)
	if !matched {
		t.Error("Expected regexp match on time string")
	}
}

func TestTransactionListSort(t *testing.T) {
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

	if transactionList.Transactions[0].ID != "2" {
		t.Errorf("Expected transaction %s, but got transaction %s", t2.ID, t1.ID)
	}
}
