package vow

import "strings"

type IssueKind string

const (
	IssueKindSchema     IssueKind = "schema"
	IssueKindTransform  IssueKind = "transform"
	IssueKindValidation IssueKind = "validation"
)

// IssueKey represents a hierarchical key for identifying issues within nested data structures.
type IssueKey []string

func (k IssueKey) String() string {
	return strings.Join(k, ".")
}

type Issue struct {
	Kind     IssueKind `json:"kind"`
	Message  string    `json:"message"`
	Key      IssueKey  `json:"key"`
	Expected string    `json:"expected,omitempty"`
}

func NewIssue(kind IssueKind, message string) Issue {
	return Issue{
		Kind:    kind,
		Message: message,
	}
}

func NewIssueFromCommonSchema(schema *CommonSchema, kind IssueKind, message string) Issue {
	issue := NewIssue(kind, message)
	issue.Expected = schema.Expected

	return issue
}

func (i *Issue) PrependKey(part string) {
  i.Key = append(IssueKey{part}, i.Key...)
}

