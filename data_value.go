package vow

import "reflect"

type DataValue struct {
	// Represents the actual data value.
	value reflect.Value
	// Represents any issues associated with the data value.
	issues []Issue
}

func NewDataValue(val any) *DataValue {
	return &DataValue{value: reflect.ValueOf(val)}
}

// Issues returns the list of issues associated with the DataValue.
func (d *DataValue) Issues() []Issue {
	return d.issues
}

// HasIssues returns true if the DataValue has any issues.
func (d *DataValue) HasIssues() bool {
	return len(d.issues) > 0
}

func (d *DataValue) AddIssue(issue Issue) {
	d.issues = append(d.issues, issue)
}
