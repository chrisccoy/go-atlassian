package internal

import (
	"bytes"
	"context"
	"errors"
	model "github.com/chrisccoy/go-atlassian/pkg/infra/models"
	"github.com/chrisccoy/go-atlassian/service"
	"github.com/chrisccoy/go-atlassian/service/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_internalProjectRoleActorImpl_Add(t *testing.T) {

	payloadMocked := &struct {
		Group []string "json:\"group,omitempty\""
		Users []string "json:\"user,omitempty\""
	}{Group: []string{"jira-users"}, Users: []string{"uuid"}}

	type fields struct {
		c       service.Client
		version string
	}

	type args struct {
		ctx                context.Context
		projectKeyOrId     string
		roleId             int
		accountIds, groups []string
	}

	testCases := []struct {
		name    string
		fields  fields
		args    args
		on      func(*fields)
		wantErr bool
		Err     error
	}{
		{
			name:   "when the api version is v3",
			fields: fields{version: "3"},
			args: args{
				ctx:            context.TODO(),
				projectKeyOrId: "DUMMY",
				roleId:         10001,
				accountIds:     []string{"uuid"},
				groups:         []string{"jira-users"},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					payloadMocked).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/3/project/DUMMY/role/10001",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.ProjectRoleScheme{}).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client
			},
			wantErr: false,
			Err:     nil,
		},

		{
			name:   "when the api version is v2",
			fields: fields{version: "2"},
			args: args{
				ctx:            context.TODO(),
				projectKeyOrId: "DUMMY",
				roleId:         10001,
				accountIds:     []string{"uuid"},
				groups:         []string{"jira-users"},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					payloadMocked).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/2/project/DUMMY/role/10001",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.ProjectRoleScheme{}).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client
			},
			wantErr: false,
			Err:     nil,
		},

		{
			name:   "when the project key or id is not provided",
			fields: fields{version: "3"},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
			Err:     model.ErrNoProjectIDOrKeyError,
		},

		{
			name:   "when the project role id is not provided",
			fields: fields{version: "3"},
			args: args{
				ctx:            context.TODO(),
				projectKeyOrId: "DUMMY",
			},
			wantErr: true,
			Err:     model.ErrNoProjectRoleIDError,
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "3"},
			args: args{
				ctx:            context.TODO(),
				projectKeyOrId: "DUMMY",
				roleId:         10001,
				accountIds:     []string{"uuid"},
				groups:         []string{"jira-users"},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					payloadMocked).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/3/project/DUMMY/role/10001",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, errors.New("error, unable to create the http request"))

				fields.c = client
			},
			wantErr: true,
			Err:     errors.New("error, unable to create the http request"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			if testCase.on != nil {
				testCase.on(&testCase.fields)
			}

			newService, err := NewProjectRoleActorService(testCase.fields.c, testCase.fields.version)
			assert.NoError(t, err)

			gotResult, gotResponse, err := newService.Add(testCase.args.ctx, testCase.args.projectKeyOrId, testCase.args.roleId,
				testCase.args.accountIds, testCase.args.groups)

			if testCase.wantErr {

				if err != nil {
					t.Logf("error returned: %v", err.Error())
				}

				assert.EqualError(t, err, testCase.Err.Error())

			} else {

				assert.NoError(t, err)
				assert.NotEqual(t, gotResponse, nil)
				assert.NotEqual(t, gotResult, nil)
			}

		})
	}
}

func Test_internalProjectRoleActorImpl_Delete(t *testing.T) {

	type fields struct {
		c       service.Client
		version string
	}

	type args struct {
		ctx              context.Context
		projectKeyOrId   string
		roleId           int
		accountId, group string
	}

	testCases := []struct {
		name    string
		fields  fields
		args    args
		on      func(*fields)
		wantErr bool
		Err     error
	}{
		{
			name:   "when the api version is v3",
			fields: fields{version: "3"},
			args: args{
				ctx:            context.TODO(),
				projectKeyOrId: "DUMMY",
				roleId:         10001,
				accountId:      "uuid",
				group:          "jira-users",
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodDelete,
					"rest/api/3/project/DUMMY/role/10001?group=jira-users&user=uuid",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					nil).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client
			},
			wantErr: false,
			Err:     nil,
		},

		{
			name:   "when the api version is v2",
			fields: fields{version: "2"},
			args: args{
				ctx:            context.TODO(),
				projectKeyOrId: "DUMMY",
				roleId:         10001,
				accountId:      "uuid",
				group:          "jira-users",
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodDelete,
					"rest/api/2/project/DUMMY/role/10001?group=jira-users&user=uuid",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					nil).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client
			},
			wantErr: false,
			Err:     nil,
		},

		{
			name:   "when the project key or id is not provided",
			fields: fields{version: "3"},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
			Err:     model.ErrNoProjectIDOrKeyError,
		},

		{
			name:   "when the project role id is not provided",
			fields: fields{version: "3"},
			args: args{
				ctx:            context.TODO(),
				projectKeyOrId: "DUMMY",
			},
			wantErr: true,
			Err:     model.ErrNoProjectRoleIDError,
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "3"},
			args: args{
				ctx:            context.TODO(),
				projectKeyOrId: "DUMMY",
				roleId:         10001,
				accountId:      "uuid",
				group:          "jira-users",
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodDelete,
					"rest/api/3/project/DUMMY/role/10001?group=jira-users&user=uuid",
					nil).
					Return(&http.Request{}, errors.New("error, unable to create the http request"))

				fields.c = client
			},
			wantErr: true,
			Err:     errors.New("error, unable to create the http request"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			if testCase.on != nil {
				testCase.on(&testCase.fields)
			}

			newService, err := NewProjectRoleActorService(testCase.fields.c, testCase.fields.version)
			assert.NoError(t, err)

			gotResponse, err := newService.Delete(testCase.args.ctx, testCase.args.projectKeyOrId, testCase.args.roleId,
				testCase.args.accountId, testCase.args.group)

			if testCase.wantErr {

				if err != nil {
					t.Logf("error returned: %v", err.Error())
				}

				assert.EqualError(t, err, testCase.Err.Error())

			} else {

				assert.NoError(t, err)
				assert.NotEqual(t, gotResponse, nil)
			}

		})
	}
}
