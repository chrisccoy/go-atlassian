package jira

import (
	"context"
	model "github.com/chrisccoy/go-atlassian/pkg/infra/models"
)

type ServerConnector interface {

	// Info returns information about the Jira instance
	//
	// GET /rest/api/{2-3}/serverInfo
	//
	// https://docs.go-atlassian.io/jira-software-cloud/server#get-jira-instance-info
	Info(ctx context.Context) (*model.ServerInformationScheme, *model.ResponseScheme, error)
}
