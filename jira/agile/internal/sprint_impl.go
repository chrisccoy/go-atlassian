package internal

import (
	"context"
	"fmt"
	model "github.com/chrisccoy/go-atlassian/pkg/infra/models"
	"github.com/chrisccoy/go-atlassian/service"
	"github.com/chrisccoy/go-atlassian/service/agile"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func NewSprintService(client service.Client, version string) (*SprintService, error) {

	if version == "" {
		return nil, model.ErrNoVersionProvided
	}

	return &SprintService{
		internalClient: &internalSprintImpl{c: client, version: version},
	}, nil
}

type SprintService struct {
	internalClient agile.SprintConnector
}

// Get Returns the sprint for a given sprint ID.
//
// The sprint will only be returned if the user can view the board that the sprint was created on,
//
// or view at least one of the issues in the sprint.
//
// GET /rest/agile/1.0/sprint/{sprintId}
//
// https://docs.go-atlassian.io/jira-agile/sprints#get-sprint
func (s *SprintService) Get(ctx context.Context, sprintID int) (*model.SprintScheme, *model.ResponseScheme, error) {
	return s.internalClient.Get(ctx, sprintID)
}

// Create creates a future sprint.
//
// Sprint name and origin board id are required.
//
// Start date, end date, and goal are optional.
//
// POST /rest/agile/1.0/sprint
//
// https://docs.go-atlassian.io/jira-agile/sprints#create-print
func (s *SprintService) Create(ctx context.Context, payload *model.SprintPayloadScheme) (*model.SprintScheme, *model.ResponseScheme, error) {
	return s.internalClient.Create(ctx, payload)
}

// Update Performs a full update of a sprint.
//
// A full update means that the result will be exactly the same as the request body.
//
// Any fields not present in the request JSON will be set to null.
//
// PUT /rest/agile/1.0/sprint/{sprintId}
//
// https://docs.go-atlassian.io/jira-agile/sprints#update-sprint
func (s *SprintService) Update(ctx context.Context, sprintID int, payload *model.SprintPayloadScheme) (*model.SprintScheme, *model.ResponseScheme, error) {
	return s.internalClient.Update(ctx, sprintID, payload)
}

// Path Performs a partial update of a sprint.
//
// A partial update means that fields not present in the request JSON will not be updated.
//
// POST /rest/agile/1.0/sprint/{sprintId}
//
// https://docs.go-atlassian.io/jira-agile/sprints#partially-update-sprint
func (s *SprintService) Path(ctx context.Context, sprintID int, payload *model.SprintPayloadScheme) (*model.SprintScheme, *model.ResponseScheme, error) {
	return s.internalClient.Path(ctx, sprintID, payload)
}

// Delete deletes a sprint.
//
// Once a sprint is deleted, all open issues in the sprint will be moved to the backlog.
//
// DELETE /rest/agile/1.0/sprint/{sprintId}
//
// https://docs.go-atlassian.io/jira-agile/sprints#delete-sprint
func (s *SprintService) Delete(ctx context.Context, sprintID int) (*model.ResponseScheme, error) {
	return s.internalClient.Delete(ctx, sprintID)
}

// Issues returns all issues in a sprint, for a given sprint ID.
//
// This only includes issues that the user has permission to view.
//
// By default, the returned issues are ordered by rank.
//
// GET /rest/agile/1.0/sprint/{sprintId}/issue
//
// https://docs.go-atlassian.io/jira-agile/sprints#get-issues-for-sprint
func (s *SprintService) Issues(ctx context.Context, sprintID int, opts *model.IssueOptionScheme, startAt, maxResults int) (*model.SprintIssuePageScheme, *model.ResponseScheme, error) {
	return s.internalClient.Issues(ctx, sprintID, opts, startAt, maxResults)
}

// Start initiate the Sprint
//
// PUT /rest/agile/1.0/sprint/{sprintId}
//
// https://docs.go-atlassian.io/jira-agile/sprints#start-sprint
func (s *SprintService) Start(ctx context.Context, sprintID int) (*model.ResponseScheme, error) {
	return s.internalClient.Start(ctx, sprintID)
}

// Close closes the Sprint
//
// PUT /rest/agile/1.0/sprint/{sprintId}
//
// https://docs.go-atlassian.io/jira-agile/sprints#close-sprint
func (s *SprintService) Close(ctx context.Context, sprintID int) (*model.ResponseScheme, error) {
	return s.internalClient.Close(ctx, sprintID)
}

type internalSprintImpl struct {
	c       service.Client
	version string
}

func (i *internalSprintImpl) Get(ctx context.Context, sprintID int) (*model.SprintScheme, *model.ResponseScheme, error) {

	if sprintID == 0 {
		return nil, nil, model.ErrNoSprintIDError
	}

	endpoint := fmt.Sprintf("rest/agile/%v/sprint/%v", i.version, sprintID)

	request, err := i.c.NewRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	sprint := new(model.SprintScheme)
	response, err := i.c.Call(request, sprint)
	if err != nil {
		return nil, response, err
	}

	return sprint, response, nil
}

func (i *internalSprintImpl) Create(ctx context.Context, payload *model.SprintPayloadScheme) (*model.SprintScheme, *model.ResponseScheme, error) {

	reader, err := i.c.TransformStructToReader(payload)
	if err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("rest/agile/%v/sprint", i.version)

	request, err := i.c.NewRequest(ctx, http.MethodPost, endpoint, reader)
	if err != nil {
		return nil, nil, err
	}

	sprint := new(model.SprintScheme)
	response, err := i.c.Call(request, sprint)
	if err != nil {
		return nil, response, err
	}

	return sprint, response, nil
}

func (i *internalSprintImpl) Update(ctx context.Context, sprintID int, payload *model.SprintPayloadScheme) (*model.SprintScheme, *model.ResponseScheme, error) {

	if sprintID == 0 {
		return nil, nil, model.ErrNoSprintIDError
	}

	reader, err := i.c.TransformStructToReader(payload)
	if err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("rest/agile/%v/sprint/%v", i.version, sprintID)

	request, err := i.c.NewRequest(ctx, http.MethodPut, endpoint, reader)
	if err != nil {
		return nil, nil, err
	}

	sprint := new(model.SprintScheme)
	response, err := i.c.Call(request, sprint)
	if err != nil {
		return nil, response, err
	}

	return sprint, response, nil
}

func (i *internalSprintImpl) Path(ctx context.Context, sprintID int, payload *model.SprintPayloadScheme) (*model.SprintScheme, *model.ResponseScheme, error) {

	if sprintID == 0 {
		return nil, nil, model.ErrNoSprintIDError
	}

	reader, err := i.c.TransformStructToReader(payload)
	if err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("rest/agile/%v/sprint/%v", i.version, sprintID)

	request, err := i.c.NewRequest(ctx, http.MethodPost, endpoint, reader)
	if err != nil {
		return nil, nil, err
	}

	sprint := new(model.SprintScheme)
	response, err := i.c.Call(request, sprint)
	if err != nil {
		return nil, response, err
	}

	return sprint, response, nil
}

func (i *internalSprintImpl) Delete(ctx context.Context, sprintID int) (*model.ResponseScheme, error) {

	if sprintID == 0 {
		return nil, model.ErrNoSprintIDError
	}

	endpoint := fmt.Sprintf("rest/agile/%v/sprint/%v", i.version, sprintID)

	request, err := i.c.NewRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return i.c.Call(request, nil)
}

func (i *internalSprintImpl) Issues(ctx context.Context, sprintID int, opts *model.IssueOptionScheme, startAt, maxResults int) (*model.SprintIssuePageScheme, *model.ResponseScheme, error) {

	if sprintID == 0 {
		return nil, nil, model.ErrNoSprintIDError
	}

	params := url.Values{}
	params.Add("startAt", strconv.Itoa(startAt))
	params.Add("maxResults", strconv.Itoa(maxResults))

	if opts != nil {

		if !opts.ValidateQuery {
			params.Add("validateQuery", "false")
		}

		if len(opts.JQL) != 0 {
			params.Add("jql", opts.JQL)
		}

		if len(opts.Expand) != 0 {
			params.Add("expand", strings.Join(opts.Expand, ","))
		}

		if len(opts.Fields) != 0 {
			params.Add("fields", strings.Join(opts.Fields, ","))
		}
	}

	endpoint := fmt.Sprintf("rest/agile/%v/sprint/%v/issue?%v", i.version, sprintID, params.Encode())

	request, err := i.c.NewRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	page := new(model.SprintIssuePageScheme)
	response, err := i.c.Call(request, page)
	if err != nil {
		return nil, response, err
	}

	return page, response, nil
}

func (i *internalSprintImpl) Start(ctx context.Context, sprintID int) (*model.ResponseScheme, error) {

	if sprintID == 0 {
		return nil, model.ErrNoSprintIDError
	}

	payload := model.SprintPayloadScheme{
		State: "Active",
	}

	reader, err := i.c.TransformStructToReader(payload)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("rest/agile/%v/sprint/%v", i.version, sprintID)

	request, err := i.c.NewRequest(ctx, http.MethodPost, endpoint, reader)
	if err != nil {
		return nil, err
	}

	return i.c.Call(request, nil)
}

func (i *internalSprintImpl) Close(ctx context.Context, sprintID int) (*model.ResponseScheme, error) {

	if sprintID == 0 {
		return nil, model.ErrNoSprintIDError
	}

	payload := model.SprintPayloadScheme{
		State: "Closed",
	}

	reader, err := i.c.TransformStructToReader(payload)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("rest/agile/%v/sprint/%v", i.version, sprintID)

	request, err := i.c.NewRequest(ctx, http.MethodPost, endpoint, reader)
	if err != nil {
		return nil, err
	}

	return i.c.Call(request, nil)
}
