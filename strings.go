package vow

import (
	"context"
	"reflect"
	"regexp"
)

// All the regular expressions used in validation schemas.
var (
	EMAIL_REGEX = regexp.MustCompile(`(?i)^[\w+-]+(?:\.[\w+-]+)*@[\da-z]+(?:[.-][\da-z]+)*\.[a-z]{2,}$`)
)

type emailSchema struct{ *CommonSchema }

func Email() Schema {
	return &emailSchema{
		CommonSchema: &CommonSchema{Expected: EMAIL_REGEX.String()},
	}
}

func (s *emailSchema) Run(ctx context.Context, dv *DataValue) error {
	if dv.value.Kind() == reflect.String && !EMAIL_REGEX.MatchString(dv.value.String()) {
		dv.AddIssue(NewIssueFromCommonSchema(s.CommonSchema, IssueKindValidation, "invalid email address"))
	}

	return nil
}
