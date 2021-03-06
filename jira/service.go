package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/trevor-atlas/zilla/util"
)

type ClientService interface {
	GetIssues(ctx context.Context) (*JiraIssues, error)
	GetIssue(ctx context.Context, issueNumber string) (*JiraIssue, error)
	GetMappedCustomFields(ctx context.Context) (*map[string]string, error)
}

type Service struct {
	config  util.ConfigData
	client  util.RequestBuilder
	baseUrl string
}

func NewService(application *util.Zilla) ClientService {
	service := new(Service)
	service.config = *application.Config
	service.client = util.NewHTTP().
		WithHeader("Accept", "application/json")
	if service.config.Jira.CustomDomain != "" {
		service.baseUrl = service.config.Jira.CustomDomain
	} else {
		service.baseUrl = fmt.Sprintf("https://%s.atlassian.net", service.config.Jira.Orgname)
	}
	if service.config.Jira.Apikey != "" {
		service.client = service.client.WithBasicAuth(service.config.Jira.Username, service.config.Jira.Apikey)
	} else {
		service.client = service.client.WithHeader("Authorization", fmt.Sprintf("Bearer: %s", service.config.Jira.AccessToken))
	}
	return service
}

func (s *Service) GetIssues(ctx context.Context) (*JiraIssues, error) {
	var url string
	url = fmt.Sprintf("%s/rest/api/2/search?jql=assignee=currentuser()+order+by+status+asc&expand=fields", s.baseUrl)
	client := s.client.Url(url)

	body, err := client.GET()
	if err != nil {
		return nil, fmt.Errorf("there was a problem making the request to the jira API in `GetIssues`: %s", err)
	}

	parsed := JiraIssues{}
	parseError := json.Unmarshal(body, &parsed)
	if parseError != nil {
		return nil, fmt.Errorf("there was a problem parsing the jira API response:%s\n", parseError)
	}
	return &parsed, nil
}

func (s *Service) GetIssue(ctx context.Context, issueNumber string) (*JiraIssue, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s?expand=fields", s.baseUrl, issueNumber)
	client := s.client.Url(url)

	res, err := client.GET()
	if err != nil {
		return nil, fmt.Errorf("error making request: %s", err)
	}

	parsed := JiraIssue{}
	parseError := json.Unmarshal(res, &parsed)
	if parseError != nil {
		return nil, fmt.Errorf("error parsing json: %s", parseError)
	}
	return &parsed, nil
}

func (s *Service) getFieldsList(ctx context.Context) ([]Field, error) {
	url := fmt.Sprintf("%s/rest/api/2/field", s.baseUrl)
	client := s.client.Url(url)
	res, err := client.GET()
	if err != nil {
		return nil, fmt.Errorf("error making fields request")
	}

	var fieldList []Field
	parseError := json.Unmarshal(res, &fieldList)
	if parseError != nil {
		return nil, fmt.Errorf("error parsing fields json\n %s", parseError)
	}
	return fieldList, nil
}

// GetMappedCustomFields returns a map of human readable field names to their id
// { "Participants": "customfield_12345" }
// that can be later merged with issues to get the entirety of their contents.
func (s *Service) GetMappedCustomFields(ctx context.Context) (*map[string]string, error) {
	fields, err := s.getFieldsList(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the list of fields: %s", err)
	}

	var fieldMapping = make(map[string]string)
	for _, f := range fields {
		if f.Custom {
			fieldMapping[f.Name] = f.ID
		}
	}
	return &fieldMapping, nil
}
