package vow

import (
	"context"
	"fmt"
	"reflect"
)

func Parse(ctx context.Context, schema Schema, input any, out any) error {
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("output must be a non-nil pointer")
	}

	targetType := rv.Elem().Type()
	dv := NewDataValue(input)

	if err := schema.Run(ctx, dv); err != nil {
		return NewParseErr(dv.Issues(), err)
	}

	outputVal := reflect.New(targetType).Elem()
	if dv.value.Type().AssignableTo(targetType) {
		outputVal.Set(dv.value)
	}

	if !dv.value.Type().ConvertibleTo(targetType) {
		return fmt.Errorf("cannot convert value of type %s to %s", dv.value.Type(), targetType)
	}

	outputVal.Set(dv.value.Convert(targetType))
	rv.Elem().Set(outputVal)

	// doing this check last so we can return a partially fill output
	if dv.HasIssues() {
		return NewParseErr(dv.Issues(), nil)
	}

	return nil
}

// ParseAs is a convenience function that combines Parse with generic return types.
func ParseAs[out any](ctx context.Context, schema Schema, input any) (out, error) {
	var result out
	err := Parse(ctx, schema, input, &result)

	return result, err
}
