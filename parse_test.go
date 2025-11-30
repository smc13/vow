package vow

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ParseSuite struct {
	suite.Suite
}

func TestParseSuite(t *testing.T) {
	suite.Run(t, new(ParseSuite))
}

func (s *ParseSuite) TestParse_BasicSuccess() {
	var result int
	err := Parse(s.T().Context(), Int(), 42, &result)
	s.Require().NoError(err)
	s.Require().Equal(42, result)
}

func (s *ParseSuite) TestParse_BasicFailure() {
	var result int
	err := Parse(s.T().Context(), Int(), "not an int", &result)
	expectErr := &ParseErr{}
	s.Require().ErrorAs(err, &expectErr)
	s.Require().Len(expectErr.Issues, 1)
}

type userConfig struct {
	Email string `vow:"email"`
}

type user struct {
	FirstName string `vow:"first_name"`
	LastName  string `json:"last_name"`
	Admin     bool
	Config    userConfig `vow:"config"`
}

func (s *ParseSuite) TestParse_AdvancedSuccess() {
	schema := Struct[user](Fields{
		"first_name": String(),
		"last_name":  String(),
		"Admin":      Bool(),
		"config": Struct[userConfig](Fields{
			"email": String(),
		}),
	})

	input := map[string]any{
		"first_name": "John",
		"last_name":  "Doe",
		"Admin":      true,
		"config": map[string]any{
			"email": "test@example.com",
		},
	}

	var result user
	err := Parse(s.T().Context(), schema, input, &result)
	s.Require().NoError(err)
	s.Require().Equal(user{
		FirstName: "John",
		LastName:  "Doe",
		Admin:     true,
		Config:    userConfig{Email: "test@example.com"},
	}, result)
}

func (s *ParseSuite) TestParse_AdvancedFailure() {
	schema := Struct[user](Fields{
		"first_name": String(),
		"last_name":  String(),
		"Admin":      Bool(),
		"config": Struct[userConfig](Fields{
			"email": String(),
		}),
	})

	input := map[string]any{
		"first_name": "John",
		"last_name":  123,   // invalid type
		"Admin":      "yes", // invalid type
		"config": map[string]any{
			"email": 456, // invalid type
		},
	}

	var result user
	err := Parse(s.T().Context(), schema, input, &result)
	expectErr := &ParseErr{}
	s.Require().ErrorAs(err, &expectErr)
	s.Require().Len(expectErr.Issues, 3)

	s.True(expectErr.HasIssueWithKey("last_name"))
	s.True(expectErr.HasIssueWithKey("Admin"))
	s.True(expectErr.HasIssueWithKey("config.email"))
}
