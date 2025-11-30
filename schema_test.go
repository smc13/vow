package vow

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testTypeSchema(t *testing.T, schema Schema, validValue any, invalidValue any) {
	dv := &DataValue{value: reflect.ValueOf(validValue)}
	require.NoError(t, schema.Run(t.Context(), dv))
	assert.False(t, dv.HasIssues())

	dv = &DataValue{value: reflect.ValueOf(invalidValue)}
	require.NoError(t, schema.Run(t.Context(), dv))
	assert.True(t, dv.HasIssues())
	assert.Equal(t, IssueKindSchema, dv.Issues()[0].Kind)
}

func TestBoolSchema(t *testing.T) {
	testTypeSchema(t, Bool(), true, "not a bool")
}

func TestStringSchema(t *testing.T) {
	testTypeSchema(t, String(), "a string", 123)
}

func TestFloat64Schema(t *testing.T) {
	testTypeSchema(t, Float64(), 3.14, "not a float")
}

func TestIntSchema(t *testing.T) {
	testTypeSchema(t, Int(), 42, 3.14)
}

func TestMapSchema(t *testing.T) {
	testTypeSchema(t, Map(), map[string]int{"a": 1}, []int{1, 2, 3})
}

func TestSliceSchema(t *testing.T) {
	testTypeSchema(t, Slice(), []int{1, 2, 3}, map[string]int{"a": 1})
}
