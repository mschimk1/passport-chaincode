package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// HandlerFunc is a chaincode API handler function type
type HandlerFunc func(stub shim.ChaincodeStubInterface, args []string) ([]byte, error)

// FuncMap is a mapping of function name to handler function
type FuncMap struct {
	handlers map[string]HandlerFunc
}

// NewHandlerMap creates a new handler mapping and returns a pointer
func NewHandlerMap() *FuncMap {
	return &FuncMap{make(map[string]HandlerFunc)}
}

// Add registers a handler function
func (p *FuncMap) Add(name string, handler HandlerFunc) {
	p.handlers[name] = handler
}

// Handle gets a handler function by name and invokes it
func (p *FuncMap) Handle(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	for name, handlerFunc := range p.handlers {
		if name == function {
			return handlerFunc(stub, args)
		}
	}
	return nil, fmt.Errorf("Handler function with name \"%s\" not registered.", function)
}
