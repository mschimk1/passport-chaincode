package main

import (
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// test handler function
func testFn(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return []byte("Success"), nil
}

func TestCreateHandlerMap(t *testing.T) {
	funcMap := NewHandlerMap()
	if len(funcMap.handlers) != 0 {
		t.Errorf("Expected empty handler map, but got %v", funcMap.handlers)
	}
}

func TestAddOneHandler(t *testing.T) {
	funcMap := NewHandlerMap()
	funcMap.Add("testFn", testFn)
	if len(funcMap.handlers) != 1 {
		t.Errorf("Expected 1 handler function, but got %d", len(funcMap.handlers))
	}
}

func TestAddMultipleHandlers(t *testing.T) {
	funcMap := NewHandlerMap()
	funcMap.Add("testFn1", testFn)
	funcMap.Add("testFn2", testFn)
	if len(funcMap.handlers) != 2 {
		t.Errorf("Expected 1 handler function, but got %d", len(funcMap.handlers))
	}
}

func TestHandlerIdentity(t *testing.T) {
	funcMap := NewHandlerMap()
	funcMap.Add("testFn", testFn)
	fn1 := reflect.ValueOf(funcMap.handlers["testFn"])
	fn2 := reflect.ValueOf(testFn)
	if fn1.Pointer() != fn2.Pointer() {
		t.Error("Expected same handler function, but got different function pointers")
	}
}

func TestHandleCallsHandlerFunction(t *testing.T) {
	funcMap := NewHandlerMap()
	funcMap.Add("testFn", testFn)
	res, err := funcMap.Handle(nil, "testFn", nil)
	if err != nil {
		t.Errorf("Expected handler function to be called, but got error: %s", err)
	}
	if string(res[:]) != "Success" {
		t.Errorf("Expected 'Success' function result, but got %s", res)
	}
}

func TestHandleWithoutHandlerFunction(t *testing.T) {
	funcMap := NewHandlerMap()
	_, err := funcMap.Handle(nil, "testFn", nil)
	if err == nil {
		t.Error("Expected error when trying to handle unknown function")
	}
}
