package model

import (
	"encoding/json"

	"github.com/stretchr/testify/suite"
)

type RatesSuite struct {
	suite.Suite
	testRate      Rates
	testRateBytes []byte
}

func (suite *RatesSuite) SetupTest() {
	suite.testRate = Rates{
		Entity: Entity{"Rates"},
		Base:   "AUD",
		Date:   "2017-08-14",
		Currencies: struct {
			HKD float64
			NZD float64
			SGD float64
		}{
			HKD: 6.1484,
			NZD: 1.0793,
			SGD: 1.0729,
		},
	}
	suite.testRateBytes = []byte(`{"docType":"Rates","base":"AUD","date":"2017-08-14","rates":{"HKD":6.1484,"NZD":1.0793,"SGD":1.0729}}`)
}

func (suite *RatesSuite) TestGetObjectType() {
	suite.Equal(suite.testRate.GetObjectType(), "Rates")
}

func (suite *RatesSuite) TestUnmarshalRates() {
	expected := suite.testRate
	actual := Rates{}
	json.Unmarshal(suite.testRateBytes, &actual)
	suite.Equal(expected, actual, "Rates should be equal")
}

func (suite *RatesSuite) TestMarshalRates() {
	expected := suite.testRateBytes
	actual, _ := json.Marshal(suite.testRate)
	suite.Equal(expected, actual)
}
