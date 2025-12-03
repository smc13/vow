package vow

import (
	"context"
	"fmt"
	"reflect"
)

type minSchema struct {
	*CommonSchema
	min float64
}

func Min(min float64) Schema {
	return &minSchema{
		CommonSchema: &CommonSchema{Expected: "minimum value"},
		min:          min,
	}
}

func (s *minSchema) Run(ctx context.Context, dv *DataValue) error {
	var value float64
	switch dv.value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = float64(dv.value.Int())
	case reflect.Float32, reflect.Float64:
		value = dv.value.Float()
	default:
		return nil
	}

	if value < s.min {
		dv.AddIssue(NewIssueFromCommonSchema(s.CommonSchema, IssueKindValidation, fmt.Sprintf("value %v is less than minimum %v", value, s.min)))
	}

	return nil
}
