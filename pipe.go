package vow

import "context"

type pipeSchema struct {
	schemas []Schema
}

func Pipe(first Schema, rest ...Schema) Schema {
	return &pipeSchema{
		schemas: append([]Schema{first}, rest...),
	}
}

func (p *pipeSchema) Run(ctx context.Context, dv *DataValue) error {
	for _, schema := range p.schemas {
		if err := schema.Run(ctx, dv); err != nil {
			return err
		}
	}

	return nil
}
