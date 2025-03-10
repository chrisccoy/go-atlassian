package jira

import (
	"context"
	model "github.com/chrisccoy/go-atlassian/pkg/infra/models"
)

type ResolutionConnector interface {

	// Gets returns a list of all issue resolution values.
	//
	// GET /rest/api/{2-3}/resolution
	//
	// https://docs.go-atlassian.io/jira-software-cloud/issues/resolutions#get-resolutions
	Gets(ctx context.Context) ([]*model.ResolutionScheme, *model.ResponseScheme, error)

	// Get returns an issue resolution value.
	//
	//
	// GET /rest/api/{2-3}/resolution/{id}
	//
	// https://docs.go-atlassian.io/jira-software-cloud/issues/resolutions#get-resolution
	Get(ctx context.Context, resolutionId string) (*model.ResolutionScheme, *model.ResponseScheme, error)
}
