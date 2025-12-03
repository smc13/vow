package vow

import "context"

type CheckFunc func(context.Context, *DataValue) error
type checkSchema struct {
	*CommonSchema
	checkFunc CheckFunc
}

// Check creates a schema that runs a custom check function.
// The provided CheckFunc receives the current context and the DataValue to be validated.
// It should return an error if the validation fails and the pipeline should be halted.
func Check(check CheckFunc) Schema {
	return &checkSchema{CommonSchema: &CommonSchema{Expected: "custom check"}, checkFunc: check}
}

func (s *checkSchema) Run(ctx context.Context, dv *DataValue) error {
	return s.checkFunc(ctx, dv)
}
