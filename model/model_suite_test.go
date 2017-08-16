package model

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(RatesSuite))
	suite.Run(t, new(UserSuite))
	suite.Run(t, new(AccountSuite))
	suite.Run(t, new(TransactionSuite))
	suite.Run(t, new(TransferSuite))
}
