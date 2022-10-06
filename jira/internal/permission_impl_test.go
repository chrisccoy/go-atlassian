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

func Test_internalPermissionImpl_Gets(t *testing.T) {

	mockedBuffer := bytes.Buffer{}

	mockedBuffer.WriteString(`
	{
		"permissions": {
			"ADD_COMMENTS": {
			  "key": "ADD_COMMENTS",
			  "name": "Add Comments",
			  "type": "PROJECT",
			  "description": "Ability to comment on issues."
			}
		}
	}
`)

	type fields struct {
		c       service.Client
		version string
	}

	type args struct {
		ctx context.Context
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
				ctx: context.TODO(),
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/3/permissions",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					nil).
					Return(&model.ResponseScheme{Bytes: mockedBuffer}, nil)

				fields.c = client
			},
			wantErr: false,
			Err:     nil,
		},

		{
			name:   "when the api version is v2",
			fields: fields{version: "2"},
			args: args{
				ctx: context.TODO(),
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/2/permissions",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					nil).
					Return(&model.ResponseScheme{Bytes: mockedBuffer}, nil)

				fields.c = client
			},
			wantErr: false,
			Err:     nil,
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "3"},
			args: args{
				ctx: context.TODO(),
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/3/permissions",
					nil).
					Return(&http.Request{}, errors.New("error, unable to create the http request"))

				fields.c = client
			},
			wantErr: true,
			Err:     errors.New("error, unable to create the http request"),
		},

		{
			name:   "when the response cannot be parsed",
			fields: fields{version: "3"},
			args: args{
				ctx: context.TODO(),
			},
			on: func(fields *fields) {

				invalidBufferMocked := bytes.Buffer{}

				invalidBufferMocked.WriteString(`
					{
						"permissions": {
							"ADD_COMMENTS": {
							  "key": "ADD_COMMENTS",
							  "name": "Add Comments",
							  "type": "PROJECT",
							  "description": "Ability to comment on issues."
							}
				`)

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/3/permissions",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					nil).
					Return(&model.ResponseScheme{Bytes: invalidBufferMocked}, nil)

				fields.c = client
			},
			wantErr: true,
			Err:     errors.New("unexpected end of JSON input"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			if testCase.on != nil {
				testCase.on(&testCase.fields)
			}

			newService, err := NewPermissionService(testCase.fields.c, testCase.fields.version, nil)
			assert.NoError(t, err)

			gotResult, gotResponse, err := newService.Gets(testCase.args.ctx)

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

func Test_internalPermissionImpl_Checks(t *testing.T) {

	payloadMocked := &model.PermissionCheckPayload{
		GlobalPermissions: []string{"ADMINISTER"},
		AccountID:         "", //
		ProjectPermissions: []*model.BulkProjectPermissionsScheme{
			{
				Issues:      nil,
				Projects:    []int{10000},
				Permissions: []string{"EDIT_ISSUES"},
			},
		},
	}

	type fields struct {
		c       service.Client
		version string
	}

	type args struct {
		ctx     context.Context
		payload *model.PermissionCheckPayload
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
				ctx:     context.TODO(),
				payload: payloadMocked,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					payloadMocked).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/3/permissions/check",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.PermissionGrantsScheme{}).
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
				ctx:     context.TODO(),
				payload: payloadMocked,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					payloadMocked).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/2/permissions/check",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.PermissionGrantsScheme{}).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client
			},
			wantErr: false,
			Err:     nil,
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "3"},
			args: args{
				ctx:     context.TODO(),
				payload: payloadMocked,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					payloadMocked).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/3/permissions/check",
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

			newService, err := NewPermissionService(testCase.fields.c, testCase.fields.version, nil)
			assert.NoError(t, err)

			gotResult, gotResponse, err := newService.Check(testCase.args.ctx, testCase.args.payload)

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

func Test_internalPermissionImpl_Projects(t *testing.T) {

	payloadMocked := &struct {
		Permissions []string "json:\"permissions,omitempty\""
	}{Permissions: []string{"EDIT_ISSUES", "CREATE_ISSUES"}}

	type fields struct {
		c       service.Client
		version string
	}

	type args struct {
		ctx         context.Context
		permissions []string
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
				ctx:         context.TODO(),
				permissions: []string{"EDIT_ISSUES", "CREATE_ISSUES"},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					payloadMocked).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/3/permissions/project",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.PermittedProjectsScheme{}).
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
				ctx:         context.TODO(),
				permissions: []string{"EDIT_ISSUES", "CREATE_ISSUES"},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					payloadMocked).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/2/permissions/project",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.PermittedProjectsScheme{}).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client
			},
			wantErr: false,
			Err:     nil,
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "3"},
			args: args{
				ctx:         context.TODO(),
				permissions: []string{"EDIT_ISSUES", "CREATE_ISSUES"},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					payloadMocked).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/3/permissions/project",
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

			newService, err := NewPermissionService(testCase.fields.c, testCase.fields.version, nil)
			assert.NoError(t, err)

			gotResult, gotResponse, err := newService.Projects(testCase.args.ctx, testCase.args.permissions)

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
