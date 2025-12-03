package vow

import "fmt"

type ParseErr struct {
	// err is the underlying error that caused the parse to fail, if any.
	err error
	// Issues contains the list of validation issues encountered during parsing of the input.
	Issues []Issue
}

func NewParseErr(issues []Issue, err error) *ParseErr {
	return &ParseErr{Issues: issues, err: err}
}

func (e *ParseErr) Error() string {
	return fmt.Sprintf("schema validation failed with %d issues", len(e.Issues))
}

func (e *ParseErr) Unwrap() error {
	return e.err
}

func (e *ParseErr) HasIssues() bool {
	return len(e.Issues) > 0
}

func (e *ParseErr) HasIssueWithKey(key string) bool {
	for _, issue := range e.Issues {
		if issue.Key.String() == key {
			return true
		}
	}

	return false
}

func (e *ParseErr) GetIssuesForKey(key string) []Issue {
	var issues []Issue
	for _, issue := range e.Issues {
		if issue.Key.String() == key {
			issues = append(issues, issue)
		}
	}

	return issues
}
