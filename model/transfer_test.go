package model

import "github.com/stretchr/testify/suite"

type TransferSuite struct {
	suite.Suite
}

func (suite *TransferSuite) TestValidateHappyPath() {
	transfer := &Transfer{"1", "1234", "2", "5678", 100, 0, "AUD", "", map[string]string(nil)}
	err := transfer.Validate()
	suite.Nil(err)
}

func (suite *TransferSuite) TestMissingFromCustomer() {
	transfer := &Transfer{"", "1234", "2", "5678", 100, 0, "AUD", "", map[string]string(nil)}
	err := transfer.Validate()
	suite.NotNil(err)
}

func (suite *TransactionSuite) TestValidateMissingFromAccount() {
	transfer := &Transfer{"1", "", "2", "5678", 100, 0, "AUD", "", map[string]string(nil)}
	err := transfer.Validate()
	suite.NotNil(err)
}

func (suite *TransferSuite) TestMissingToCustomer() {
	transfer := &Transfer{"1", "1234", "", "5678", 100, 0, "AUD", "", map[string]string(nil)}
	err := transfer.Validate()
	suite.NotNil(err)
}

func (suite *TransactionSuite) TestValidateMissingToAccount() {
	transfer := &Transfer{"1", "1234", "2", "", 100, 0, "AUD", "", map[string]string(nil)}
	err := transfer.Validate()
	suite.NotNil(err)
}

func (suite *TransactionSuite) TestValidateIncorrectAmount() {
	transfer := &Transfer{"1", "1234", "2", "5678", 0, 0, "AUD", "", map[string]string(nil)}
	err := transfer.Validate()
	suite.NotNil(err)
}

func (suite *TransactionSuite) TestValidateMissingCurrency() {
	transfer := &Transfer{"1", "1234", "2", "5678", 100, 0, "", "", map[string]string(nil)}
	err := transfer.Validate()
	suite.NotNil(err)
}
