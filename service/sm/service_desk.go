package sm

import (
	"context"
	model "github.com/chrisccoy/go-atlassian/pkg/infra/models"
	"io"
)

type ServiceDeskConnector interface {

	// Gets returns all the service desks in the Jira Service Management instance that the user has permission to access.
	//
	// GET /rest/servicedeskapi/servicedesk
	//
	// https://docs.go-atlassian.io/jira-service-management-cloud/request/service-desk#get-service-desks
	Gets(ctx context.Context, start, limit int) (*model.ServiceDeskPageScheme, *model.ResponseScheme, error)

	// Get returns a service desk.
	//
	// Use this method to get service desk details whenever your application component is passed a service desk ID
	//
	// but needs to display other service desk details.
	//
	// GET /rest/servicedeskapi/servicedesk/{serviceDeskId}
	//
	// https://docs.go-atlassian.io/jira-service-management-cloud/request/service-desk#get-service-desk-by-id
	Get(ctx context.Context, serviceDeskID int) (*model.ServiceDeskScheme, *model.ResponseScheme, error)

	// Attach one temporary attachments to a service desk
	//
	// POST /rest/servicedeskapi/servicedesk/{serviceDeskId}/attachTemporaryFile
	//
	// https://docs.go-atlassian.io/jira-service-management-cloud/request/service-desk#attach-temporary-file
	Attach(ctx context.Context, serviceDeskID int, fileName string, file io.Reader) (*model.ServiceDeskTemporaryFileScheme, *model.ResponseScheme, error)
}
