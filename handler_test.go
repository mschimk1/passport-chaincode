package main

import (
	"reflect"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
}

// test handler function
func testFn(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return []byte("Success"), nil
}

func (suite *HandlerSuite) TestCreateHandlerMap() {
	funcMap := NewHandlerMap()
	suite.Equal(0, len(funcMap.handlers), "Expected empty handler map")
}

func (suite *HandlerSuite) TestAddOneHandler() {
	funcMap := NewHandlerMap()
	funcMap.Add("testFn", testFn)
	suite.Equal(1, len(funcMap.handlers), "Expected 1 handler in map")
}

func (suite *HandlerSuite) TestAddMultipleHandlers() {
	funcMap := NewHandlerMap()
	funcMap.Add("testFn1", testFn)
	funcMap.Add("testFn2", testFn)
	suite.Equal(2, len(funcMap.handlers), "Expected 2 handlers in map")
}

func (suite *HandlerSuite) TestHandlerIdentity() {
	funcMap := NewHandlerMap()
	funcMap.Add("testFn", testFn)
	fn1 := reflect.ValueOf(funcMap.handlers["testFn"])
	fn2 := reflect.ValueOf(testFn)
	suite.Equal(fn1.Pointer(), fn2.Pointer())
}

func (suite *HandlerSuite) TestHandleCallsHandlerFunction() {
	funcMap := NewHandlerMap()
	funcMap.Add("testFn", testFn)
	res, err := funcMap.Handle(nil, "testFn", nil)
	suite.Nil(err)
	suite.Equal("Success", string(res[:]))
}

func (suite *HandlerSuite) TestHandleWithoutHandlerFunction() {
	funcMap := NewHandlerMap()
	_, err := funcMap.Handle(nil, "testFn", nil)
	suite.NotNil(err)
}
