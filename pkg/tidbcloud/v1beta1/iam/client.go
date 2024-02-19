package iam

import (
	"fmt"

	"tidbcloud-cli/internal/config"
	ver "tidbcloud-cli/internal/version"

	"github.com/c4pt0r/go-tidbcloud-sdk-v1/client/project"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
	"github.com/pingcap/errors"
)

const (
	projectPath = "/v1beta1/projects"

	DefaultEndpoint = "https://iam.tidbapi.com"
)

type Service struct {
	client *resty.Client
	Url    string
}

func NewIamService(client *resty.Client, url string) *Service {
	return &Service{
		client,
		url,
	}
}

func (s *Service) ListProjects(params *ListProjectsParams) (*project.ListProjectsOK, error) {
	v, _ := query.Values(params)
	var result project.ListProjectsOKBody
	resp, err := s.client.R().
		SetHeader("user-agent", fmt.Sprintf("%s/%s", config.CliName, ver.Version)).
		SetResult(&result).SetQueryString(v.Encode()).
		Get(fmt.Sprintf("%s%s", s.Url, projectPath))
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, errors.Errorf("Failed to get list projects, code: %s, body: %s", resp.Status(), string(resp.Body()))
	}

	return &project.ListProjectsOK{
		Payload: &result,
	}, nil
}

type ListProjectsParams struct {
	Page     int64 `url:"page"`
	PageSize int64 `url:"page_size"`
}
