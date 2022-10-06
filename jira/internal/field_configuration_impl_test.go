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

func Test_internalIssueFieldConfigServiceImpl_Gets(t *testing.T) {

	type fields struct {
		c       service.Client
		version string
	}

	type args struct {
		ctx                 context.Context
		ids                 []int
		isDefault           bool
		startAt, maxResults int
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
				ctx:        context.TODO(),
				ids:        []int{10000, 100001},
				isDefault:  false,
				startAt:    50,
				maxResults: 50,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/3/fieldconfiguration?id=10000&id=100001&isDefault=false&maxResults=50&startAt=50",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.FieldConfigurationPageScheme{}).
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
				ctx:        context.TODO(),
				ids:        []int{10000, 100001},
				isDefault:  false,
				startAt:    50,
				maxResults: 50,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/2/fieldconfiguration?id=10000&id=100001&isDefault=false&maxResults=50&startAt=50",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.FieldConfigurationPageScheme{}).
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
				ctx:        context.TODO(),
				ids:        []int{10000, 100001},
				isDefault:  false,
				startAt:    50,
				maxResults: 50,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/3/fieldconfiguration?id=10000&id=100001&isDefault=false&maxResults=50&startAt=50",
					nil).
					Return(&http.Request{}, errors.New("error, unable to create the http created"))

				fields.c = client

			},
			wantErr: true,
			Err:     errors.New("error, unable to create the http created"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			if testCase.on != nil {
				testCase.on(&testCase.fields)
			}

			fieldConfigService, err := NewIssueFieldConfigurationService(testCase.fields.c, testCase.fields.version, nil, nil)
			assert.NoError(t, err)

			gotResult, gotResponse, err := fieldConfigService.Gets(testCase.args.ctx, testCase.args.ids,
				testCase.args.isDefault, testCase.args.startAt, testCase.args.maxResults)

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

func Test_internalIssueFieldConfigServiceImpl_Create(t *testing.T) {

	type fields struct {
		c       service.Client
		version string
	}

	type args struct {
		ctx               context.Context
		name, description string
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
				name:        "DUMMY Field Configuration Scheme",
				description: "description sample",
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					&struct {
						Name        string "json:\"name\""
						Description string "json:\"description,omitempty\""
					}{
						Name:        "DUMMY Field Configuration Scheme",
						Description: "description sample"}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/3/fieldconfiguration",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.FieldConfigurationScheme{}).
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
				name:        "DUMMY Field Configuration Scheme",
				description: "description sample",
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					&struct {
						Name        string "json:\"name\""
						Description string "json:\"description,omitempty\""
					}{
						Name:        "DUMMY Field Configuration Scheme",
						Description: "description sample"}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/2/fieldconfiguration",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.FieldConfigurationScheme{}).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
			wantErr: false,
			Err:     nil,
		},

		{
			name:   "when the field configuration name is not provided",
			fields: fields{version: "3"},
			args: args{
				ctx:         context.TODO(),
				name:        "",
				description: "description sample",
			},
			wantErr: true,
			Err:     model.ErrNoFieldConfigurationNameError,
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "3"},
			args: args{
				ctx:         context.TODO(),
				name:        "DUMMY Field Configuration Scheme",
				description: "description sample",
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					&struct {
						Name        string "json:\"name\""
						Description string "json:\"description,omitempty\""
					}{
						Name:        "DUMMY Field Configuration Scheme",
						Description: "description sample"}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/3/fieldconfiguration",
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

			fieldConfigService, err := NewIssueFieldConfigurationService(testCase.fields.c, testCase.fields.version, nil, nil)
			assert.NoError(t, err)

			gotResult, gotResponse, err := fieldConfigService.Create(testCase.args.ctx, testCase.args.name, testCase.args.description)

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

func Test_internalIssueFieldConfigServiceImpl_Update(t *testing.T) {

	type fields struct {
		c       service.Client
		version string
	}

	type args struct {
		ctx               context.Context
		id                int
		name, description string
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
				id:          1001,
				name:        "DUMMY Field Configuration Scheme",
				description: "description sample",
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					&struct {
						Name        string "json:\"name\""
						Description string "json:\"description,omitempty\""
					}{
						Name:        "DUMMY Field Configuration Scheme",
						Description: "description sample"}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPut,
					"rest/api/3/fieldconfiguration/1001",
					bytes.NewReader([]byte{})).
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
				ctx:         context.TODO(),
				id:          1001,
				name:        "DUMMY Field Configuration Scheme",
				description: "description sample",
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					&struct {
						Name        string "json:\"name\""
						Description string "json:\"description,omitempty\""
					}{
						Name:        "DUMMY Field Configuration Scheme",
						Description: "description sample"}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPut,
					"rest/api/2/fieldconfiguration/1001",
					bytes.NewReader([]byte{})).
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
			name:   "when the field configuration name is not provided",
			fields: fields{version: "3"},
			args: args{
				ctx:         context.TODO(),
				id:          1001,
				name:        "",
				description: "description sample",
			},
			wantErr: true,
			Err:     model.ErrNoFieldConfigurationNameError,
		},

		{
			name:   "when the field configuration id is not provided",
			fields: fields{version: "3"},
			args: args{
				ctx:         context.TODO(),
				id:          0,
				name:        "field configuration id",
				description: "description sample",
			},
			wantErr: true,
			Err:     model.ErrNoFieldConfigurationIDError,
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "3"},
			args: args{
				ctx:         context.TODO(),
				id:          1001,
				name:        "DUMMY Field Configuration Scheme",
				description: "description sample",
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					&struct {
						Name        string "json:\"name\""
						Description string "json:\"description,omitempty\""
					}{
						Name:        "DUMMY Field Configuration Scheme",
						Description: "description sample"}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPut,
					"rest/api/3/fieldconfiguration/1001",
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

			fieldConfigService, err := NewIssueFieldConfigurationService(testCase.fields.c, testCase.fields.version, nil, nil)
			assert.NoError(t, err)

			gotResponse, err := fieldConfigService.Update(testCase.args.ctx, testCase.args.id, testCase.args.name,
				testCase.args.description)

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

func Test_internalIssueFieldConfigServiceImpl_Delete(t *testing.T) {

	type fields struct {
		c       service.Client
		version string
	}

	type args struct {
		ctx context.Context
		id  int
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
				id:  1001,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodDelete,
					"rest/api/3/fieldconfiguration/1001",
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
				ctx: context.TODO(),
				id:  1001,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodDelete,
					"rest/api/2/fieldconfiguration/1001",
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
			name:   "when the field configuration id is not provided",
			fields: fields{version: "3"},
			args: args{
				ctx: context.TODO(),
				id:  0,
			},
			wantErr: true,
			Err:     model.ErrNoFieldConfigurationIDError,
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "3"},
			args: args{
				ctx: context.TODO(),
				id:  1001,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodDelete,
					"rest/api/3/fieldconfiguration/1001",
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

			fieldConfigService, err := NewIssueFieldConfigurationService(testCase.fields.c, testCase.fields.version, nil, nil)
			assert.NoError(t, err)

			gotResponse, err := fieldConfigService.Delete(testCase.args.ctx, testCase.args.id)

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
