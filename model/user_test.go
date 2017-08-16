package model

import (
	"encoding/json"

	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	testUser      User
	testUserBytes []byte
}

func (suite *UserSuite) SetupTest() {
	suite.testUser = User{Entity{"User"}, "1234", "John Smith"}
	suite.testUserBytes = []byte(`{"docType":"User","id":"1234","name":"John Smith"}`)
}

func (suite *UserSuite) TestGetObjectType() {
	suite.Equal("User", suite.testUser.GetObjectType())
}

func (suite *UserSuite) TestUnmarshalUser() {
	expected := suite.testUser
	actual := User{}
	err := json.Unmarshal(suite.testUserBytes, &actual)
	suite.Nil(err)
	suite.Equal(expected, actual, "Users should be equal")
}

func (suite *UserSuite) TestMarshalUser() {
	expected := suite.testUserBytes
	actual, _ := json.Marshal(suite.testUser)
	suite.Equal(expected, actual)
}
