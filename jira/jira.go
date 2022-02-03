package jira

import (
	"time"
)

// Time represents the Time definition of JIRA as a time.Time of go
type Time time.Time

// Date represents the Date definition of JIRA as a time.Time of go
type Date time.Time

type IssueUser struct {
	Active       bool   `json:"active"`
	TimeZone     string `json:"timeZone"`
	DisplayName  string `json:"displayName"`
	Name         string
	EmailAddress string
	AvatarUrls   map[string]interface{} `json:"-"`
	AccountId    string
	Key          string
	Self         string
}

type IssueComment struct {
	ID           string
	Self         string
	Author       IssueUser
	Body         string
	UpdateAuthor IssueUser
	Created      *Time
	Updated      *Time
	Total        int
}

type IssueFields struct {
	Summary     string    // title of jira issue
	Created     *Time     `json:"created"`     // 2018-05-25T04:18:06.836-0500
	Updated     *Time     `json:"updated"`     // 2018-06-11T22:23:03.606-0500
	Description string    `json:"description"` // description of Jira issue
	Reporter    IssueUser `json:"reporter"`
	Assignee    IssueUser `json:"assignee"`
	Comment     IssueComments
	Priority    IssuePriority
	IssueType   IssueType
	Status      IssueStatus
	Project     IssueProject
}

type IssueComments struct {
	Comments []IssueComment
}

type IssuePriority struct {
	Name string `json:"priority"` // Medium
}

type IssueType struct {
	Name    string `json:"name"` // Bug, Task, Story
	Subtask bool   `json:"subtask"`
	IconURL string `json:"iconUrl"`
}

type IssueStatus struct {
	Description    string
	Name           string
	StatusCategory struct {
		Key  string
		Name string
		ID   int
	}
}

type IssueProject struct {
	Key  string
	Name string
}

// JiraIssue describes the response for a single jira issue
type JiraIssue struct {
	ID     string      `json:"id"`
	Self   string      `json:"self"` // url to request this issue
	Key    string      `json:"key"`  // XYZ-1234
	Fields IssueFields `json:"fields"`
}

type JiraIssues struct {
	Issues []JiraIssue
}

// UnmarshalJSON will transform the JIRA time into a time.Time
// during the transformation of the JIRA JSON response
func (t *Time) UnmarshalJSON(b []byte) error {
	// Ignore null, like in the main JSON package.
	if string(b) == "null" {
		return nil
	}
	ti, err := time.Parse("\"2006-01-02T15:04:05.999-0700\"", string(b))
	if err != nil {
		return err
	}
	*t = Time(ti)
	return nil
}

// UnmarshalJSON will transform the JIRA date into a time.Time
// during the transformation of the JIRA JSON response
func (t *Date) UnmarshalJSON(b []byte) error {
	// Ignore null, like in the main JSON package.
	if string(b) == "null" {
		return nil
	}
	ti, err := time.Parse("\"2006-01-02\"", string(b))
	if err != nil {
		return err
	}
	*t = Date(ti)
	return nil
}

// Field represents a field of a Jira issue.
type Field struct {
	ID          string      `json:"id,omitempty" structs:"id,omitempty"`
	Key         string      `json:"key,omitempty" structs:"key,omitempty"`
	Name        string      `json:"name,omitempty" structs:"name,omitempty"`
	Custom      bool        `json:"custom,omitempty" structs:"custom,omitempty"`
	Navigable   bool        `json:"navigable,omitempty" structs:"navigable,omitempty"`
	Searchable  bool        `json:"searchable,omitempty" structs:"searchable,omitempty"`
	ClauseNames []string    `json:"clauseNames,omitempty" structs:"clauseNames,omitempty"`
	Schema      FieldSchema `json:"schema,omitempty" structs:"schema,omitempty"`
}

// FieldSchema represents a schema of a Jira field.
// Documentation: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-issue-fields/#api-rest-api-2-field-get
type FieldSchema struct {
	Type     string `json:"type,omitempty" structs:"type,omitempty"`
	Items    string `json:"items,omitempty" structs:"items,omitempty"`
	Custom   string `json:"custom,omitempty" structs:"custom,omitempty"`
	System   string `json:"system,omitempty" structs:"system,omitempty"`
	CustomID int64  `json:"customId,omitempty" structs:"customId,omitempty"`
}
