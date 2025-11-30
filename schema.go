package vow

import (
	"context"
	"fmt"
	"reflect"
)

type Schema interface {
	// Run executes the schema and modifies the provided DataValue accordingly.
	// It should return an error only if the execution itself fails and further processing should be halted.
	Run(context.Context, *DataValue) error
}

type CommonSchema struct {
	Expected string
}

type typeSchema struct {
	*CommonSchema
	valueType reflect.Kind
}

func (s *typeSchema) Run(ctx context.Context, dv *DataValue) error {
	if dv.value.Kind() != s.valueType {
		dv.AddIssue(NewIssueFromCommonSchema(s.CommonSchema, IssueKindSchema, fmt.Sprintf("expected type %s but got %T", s.Expected, dv.value)))
	}

	return nil
}

func NewTypeSchema(valueType reflect.Kind) Schema {
	return &typeSchema{valueType: valueType, CommonSchema: &CommonSchema{Expected: valueType.String()}}
}

func Bool() Schema {
	return NewTypeSchema(reflect.Bool)
}

func String() Schema {
	return NewTypeSchema(reflect.String)
}

func Float64() Schema {
	return NewTypeSchema(reflect.Float64)
}

func Int() Schema {
	return NewTypeSchema(reflect.Int)
}

func Map() Schema {
	return NewTypeSchema(reflect.Map)
}

func Slice() Schema {
	return NewTypeSchema(reflect.Slice)
}

type Fields map[string]Schema

type structField struct {
	index  int
	key    string
	schema Schema
}

type structSchema struct {
	*CommonSchema
	t      reflect.Type
	fields []structField
}

func Struct[T any](fields Fields) Schema {
	t := reflect.TypeOf((*T)(nil)).Elem()
	schema := &structSchema{
		t:            t,
		CommonSchema: &CommonSchema{Expected: t.Name()},
	}

	lookup := structToFieldLookup(t)
	for name, subSchema := range fields {
		sf, ok := lookup[name]
		if !ok {
			panic(fmt.Sprintf("field %s not found in struct %s, found fields: %+v", name, t.Name(), lookup))
		}

		sf.schema = subSchema
		schema.fields = append(schema.fields, sf)
	}

	return schema
}

func (s *structSchema) Run(ctx context.Context, dv *DataValue) error {
	var asMap map[string]reflect.Value

	switch dv.value.Kind() {
	case reflect.Struct:
		asMap := make(map[string]reflect.Value, dv.value.NumField())
		for i := 0; i < dv.value.NumField(); i++ {
			f := s.t.Field(i)
			key := f.Tag.Get("vow")
			if key == "" {
				key = f.Tag.Get("json")
			}

			if key == "" {
				key = f.Name
			}

			asMap[key] = dv.value.Field(i)
		}

	case reflect.Map:
		if dv.value.Type().Key().Kind() != reflect.String {
			dv.AddIssue(NewIssueFromCommonSchema(s.CommonSchema, IssueKindSchema, fmt.Sprintf("expected map with string keys but got map with %s keys", dv.value.Type().Key().Kind())))
			return nil
		}

		asMap = make(map[string]reflect.Value, dv.value.Len())
		for _, keyVal := range dv.value.MapKeys() {
			keyStr := keyVal.String()
			asMap[keyStr] = dv.value.MapIndex(keyVal)
		}

	default:
		dv.AddIssue(NewIssueFromCommonSchema(s.CommonSchema, IssueKindSchema, fmt.Sprintf("expected struct or map but got %T", dv.value.Interface())))
		return nil
	}

	out := reflect.New(s.t).Elem()
	for _, field := range s.fields {
		val, ok := asMap[field.key]
		if !ok {
			val = reflect.Zero(out.Field(field.index).Type())
		}

		fieldDv := NewDataValue(val.Interface())
		if err := field.schema.Run(ctx, fieldDv); err != nil {
			return err
		}

		if fieldDv.HasIssues() {
			// move the issues up to the parent dv with a prefix on the key
			// indicating the field name
			for _, issue := range fieldDv.Issues() {
				newIssue := issue
				newIssue.PrependKey(field.key)
				dv.AddIssue(newIssue)
			}
		} else {
			out.Field(field.index).Set(fieldDv.value)
		}
	}

	dv.value = out
	return nil
}

// builds a lookup map from field name to struct field
// taking into account `vow` and `json` tags
func structToFieldLookup(t reflect.Type) map[string]structField {
	lookup := make(map[string]structField)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" && !f.Anonymous {
			// unexported field
			continue
		}

		keys := []string{f.Tag.Get("vow"), f.Tag.Get("json"), f.Name}
		for _, key := range keys {
			if key != "" {
				lookup[key] = structField{
					index:  i,
					key:    key,
					schema: nil,
				}
			}
		}
	}

	return lookup
}
