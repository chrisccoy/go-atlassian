package internal

import (
	"context"
	"fmt"
	model "github.com/chrisccoy/go-atlassian/pkg/infra/models"
	"github.com/chrisccoy/go-atlassian/service"
	"github.com/chrisccoy/go-atlassian/service/jira"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func NewUserSearchService(client service.Client, version string) (*UserSearchService, error) {

	if version == "" {
		return nil, model.ErrNoVersionProvided
	}

	return &UserSearchService{
		internalClient: &internalUserSearchImpl{c: client, version: version},
	}, nil
}

type UserSearchService struct {
	internalClient jira.UserSearchConnector
}

// Projects returns a list of users who can be assigned issues in one or more projects.
//
// The list may be restricted to users whose attributes match a string.
//
// GET /rest/api/{2-3}/user/assignable/multiProjectSearch
//
// https://docs.go-atlassian.io/jira-software-cloud/users/search#find-users-assignable-to-projects
func (u *UserSearchService) Projects(ctx context.Context, accountId string, projectKeys []string, startAt, maxResults int) ([]*model.UserScheme, *model.ResponseScheme, error) {
	return u.internalClient.Projects(ctx, accountId, projectKeys, startAt, maxResults)
}

// Do return a list of users that match the search string and property.
//
//
// This operation takes the users in the range defined by startAt and maxResults, up to the thousandth user,
//
// and then returns only the users from that range that match the search string and property.
//
// This means the operation usually returns fewer users than specified in maxResults
//
// GET /rest/api/{2-3}/user/search
//
// https://docs.go-atlassian.io/jira-software-cloud/users/search#find-users
func (u *UserSearchService) Do(ctx context.Context, accountId, query string, startAt, maxResults int) ([]*model.UserScheme, *model.ResponseScheme, error) {
	return u.internalClient.Do(ctx, accountId, query, startAt, maxResults)
}

type internalUserSearchImpl struct {
	c       service.Client
	version string
}

func (i *internalUserSearchImpl) Projects(ctx context.Context, accountId string, projectKeys []string, startAt, maxResults int) ([]*model.UserScheme, *model.ResponseScheme, error) {

	if len(projectKeys) == 0 {
		return nil, nil, model.ErrNoProjectKeySliceError
	}

	params := url.Values{}
	params.Add("startAt", strconv.Itoa(startAt))
	params.Add("maxResults", strconv.Itoa(maxResults))

	if accountId != "" {
		params.Add("accountId", accountId)
	}

	if len(projectKeys) != 0 {
		params.Add("projectKeys", strings.Join(projectKeys, ","))
	}

	endpoint := fmt.Sprintf("rest/api/%v/user/assignable/multiProjectSearch?%v", i.version, params.Encode())

	request, err := i.c.NewRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var users []*model.UserScheme
	response, err := i.c.Call(request, &users)
	if err != nil {
		return nil, response, err
	}

	return users, response, nil
}

func (i *internalUserSearchImpl) Do(ctx context.Context, accountId, query string, startAt, maxResults int) ([]*model.UserScheme, *model.ResponseScheme, error) {

	params := url.Values{}
	params.Add("startAt", strconv.Itoa(startAt))
	params.Add("maxResults", strconv.Itoa(maxResults))

	if accountId != "" {
		params.Add("accountId", accountId)
	}

	if query != "" {
		params.Add("query", query)
	}

	endpoint := fmt.Sprintf("rest/api/%v/user/search?%v", i.version, params.Encode())

	request, err := i.c.NewRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var users []*model.UserScheme
	response, err := i.c.Call(request, &users)
	if err != nil {
		return nil, response, err
	}

	return users, response, nil
}
