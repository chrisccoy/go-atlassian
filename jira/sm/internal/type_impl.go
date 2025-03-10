package internal

import (
	"context"
	"fmt"
	model "github.com/chrisccoy/go-atlassian/pkg/infra/models"
	"github.com/chrisccoy/go-atlassian/service"
	"github.com/chrisccoy/go-atlassian/service/sm"
	"net/http"
	"net/url"
	"strconv"
)

func NewTypeService(client service.Client, version string) (*TypeService, error) {

	if version == "" {
		return nil, model.ErrNoVersionProvided
	}

	return &TypeService{
		internalClient: &internalTypeImpl{c: client, version: version},
	}, nil
}

type TypeService struct {
	internalClient sm.TypeConnector
}

// Search returns all customer request types used in the Jira Service Management instance,
// optionally filtered by a query string.
//
// GET /rest/servicedeskapi/requesttype
//
// https://docs.go-atlassian.io/jira-service-management-cloud/request/types#get-all-request-types
func (t *TypeService) Search(ctx context.Context, query string, start, limit int) (*model.RequestTypePageScheme, *model.ResponseScheme, error) {
	return t.internalClient.Search(ctx, query, start, limit)
}

// Gets returns all customer request types from a service desk.
//
// GET /rest/servicedeskapi/servicedesk/{serviceDeskId}/requesttype
//
// https://docs.go-atlassian.io/jira-service-management-cloud/request/types#get-request-types
func (t *TypeService) Gets(ctx context.Context, serviceDeskID, groupID, start, limit int) (*model.ProjectRequestTypePageScheme, *model.ResponseScheme, error) {
	return t.internalClient.Gets(ctx, serviceDeskID, groupID, start, limit)
}

// Create enables a customer request type to be added to a service desk based on an issue type.
//
// POST /rest/servicedeskapi/servicedesk/{serviceDeskId}/requesttype
//
// https://docs.go-atlassian.io/jira-service-management-cloud/request/types#create-request-type
func (t *TypeService) Create(ctx context.Context, serviceDeskID int, issueTypeID, name, description, helpText string) (*model.RequestTypeScheme, *model.ResponseScheme, error) {
	return t.internalClient.Create(ctx, serviceDeskID, issueTypeID, name, description, helpText)
}

// Get returns a customer request type from a service desk.
//
// GET /rest/servicedeskapi/servicedesk/{serviceDeskId}/requesttype/{requestTypeId}
//
// https://docs.go-atlassian.io/jira-service-management-cloud/request/types#get-request-type-by-id
func (t *TypeService) Get(ctx context.Context, serviceDeskID, requestTypeID int) (*model.RequestTypeScheme, *model.ResponseScheme, error) {
	return t.internalClient.Get(ctx, serviceDeskID, requestTypeID)
}

// Delete deletes a customer request type from a service desk, and removes it from all customer requests.
//
// DELETE /rest/servicedeskapi/servicedesk/{serviceDeskId}/requesttype/{requestTypeId}
//
// https://docs.go-atlassian.io/jira-service-management-cloud/request/types#delete-request-type
func (t *TypeService) Delete(ctx context.Context, serviceDeskID, requestTypeID int) (*model.ResponseScheme, error) {
	return t.internalClient.Delete(ctx, serviceDeskID, requestTypeID)
}

// Fields returns the fields for a service desk's customer request type.
//
// GET /rest/servicedeskapi/servicedesk/{serviceDeskId}/requesttype/{requestTypeId}/field
//
// https://docs.go-atlassian.io/jira-service-management-cloud/request/types#get-request-type-fields
func (t *TypeService) Fields(ctx context.Context, serviceDeskID, requestTypeID int) (*model.RequestTypeFieldsScheme, *model.ResponseScheme, error) {
	return t.internalClient.Fields(ctx, serviceDeskID, requestTypeID)
}

type internalTypeImpl struct {
	c       service.Client
	version string
}

func (i *internalTypeImpl) Search(ctx context.Context, query string, start, limit int) (*model.RequestTypePageScheme, *model.ResponseScheme, error) {

	params := url.Values{}
	params.Add("start", strconv.Itoa(start))
	params.Add("limit", strconv.Itoa(limit))

	if len(query) != 0 {
		params.Add("searchQuery", query)
	}

	endpoint := fmt.Sprintf("rest/servicedeskapi/requesttype?%v", params.Encode())

	request, err := i.c.NewRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	page := new(model.RequestTypePageScheme)
	response, err := i.c.Call(request, page)
	if err != nil {
		return nil, response, err
	}

	return page, response, nil
}

func (i *internalTypeImpl) Gets(ctx context.Context, serviceDeskID, groupID, start, limit int) (*model.ProjectRequestTypePageScheme, *model.ResponseScheme, error) {

	if serviceDeskID == 0 {
		return nil, nil, model.ErrNoServiceDeskIDError
	}

	params := url.Values{}
	params.Add("start", strconv.Itoa(start))
	params.Add("limit", strconv.Itoa(limit))

	if groupID != 0 {
		params.Add("groupId", strconv.Itoa(groupID))
	}

	endpoint := fmt.Sprintf("rest/servicedeskapi/servicedesk/%v/requesttype?%v", serviceDeskID, params.Encode())

	request, err := i.c.NewRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	page := new(model.ProjectRequestTypePageScheme)
	response, err := i.c.Call(request, page)
	if err != nil {
		return nil, response, err
	}

	return page, response, nil
}

func (i *internalTypeImpl) Create(ctx context.Context, serviceDeskID int, issueTypeID, name, description, helpText string) (*model.RequestTypeScheme, *model.ResponseScheme, error) {

	if serviceDeskID == 0 {
		return nil, nil, model.ErrNoServiceDeskIDError
	}

	payload := struct {
		IssueTypeID string `json:"issueTypeId,omitempty"`
		HelpText    string `json:"helpText,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	}{
		IssueTypeID: issueTypeID,
		HelpText:    helpText,
		Name:        name,
		Description: description,
	}

	reader, err := i.c.TransformStructToReader(&payload)
	if err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("rest/servicedeskapi/servicedesk/%v/requesttype", serviceDeskID)

	request, err := i.c.NewRequest(ctx, http.MethodPost, endpoint, reader)
	if err != nil {
		return nil, nil, err
	}

	type_ := new(model.RequestTypeScheme)
	response, err := i.c.Call(request, type_)
	if err != nil {
		return nil, response, err
	}

	return type_, response, nil
}

func (i *internalTypeImpl) Get(ctx context.Context, serviceDeskID, requestTypeID int) (*model.RequestTypeScheme, *model.ResponseScheme, error) {

	if serviceDeskID == 0 {
		return nil, nil, model.ErrNoServiceDeskIDError
	}

	if requestTypeID == 0 {
		return nil, nil, model.ErrNoRequestTypeIDError
	}

	endpoint := fmt.Sprintf("rest/servicedeskapi/servicedesk/%v/requesttype/%v", serviceDeskID, requestTypeID)

	request, err := i.c.NewRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	type_ := new(model.RequestTypeScheme)
	response, err := i.c.Call(request, type_)
	if err != nil {
		return nil, response, err
	}

	return type_, response, nil
}

func (i *internalTypeImpl) Delete(ctx context.Context, serviceDeskID, requestTypeID int) (*model.ResponseScheme, error) {

	if serviceDeskID == 0 {
		return nil, model.ErrNoServiceDeskIDError
	}

	if requestTypeID == 0 {
		return nil, model.ErrNoRequestTypeIDError
	}

	endpoint := fmt.Sprintf("rest/servicedeskapi/servicedesk/%v/requesttype/%v", serviceDeskID, requestTypeID)

	request, err := i.c.NewRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return i.c.Call(request, nil)
}

func (i *internalTypeImpl) Fields(ctx context.Context, serviceDeskID, requestTypeID int) (*model.RequestTypeFieldsScheme, *model.ResponseScheme, error) {

	if serviceDeskID == 0 {
		return nil, nil, model.ErrNoServiceDeskIDError
	}

	if requestTypeID == 0 {
		return nil, nil, model.ErrNoRequestTypeIDError
	}

	endpoint := fmt.Sprintf("rest/servicedeskapi/servicedesk/%v/requesttype/%v/field", serviceDeskID, requestTypeID)

	request, err := i.c.NewRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	fields := new(model.RequestTypeFieldsScheme)
	response, err := i.c.Call(request, fields)
	if err != nil {
		return nil, response, err
	}

	return fields, response, nil
}
