package internal

import (
	"bytes"
	"context"
	"errors"
	model "github.com/chrisccoy/go-atlassian/pkg/infra/models"
	"github.com/chrisccoy/go-atlassian/service"
	"github.com/chrisccoy/go-atlassian/service/jira"
	"github.com/chrisccoy/go-atlassian/service/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFilterService_Create(t *testing.T) {

	type fields struct {
		c       service.Client
		version string
		share   jira.FilterSharingConnector
	}

	type args struct {
		ctx     context.Context
		payload *model.FilterPayloadScheme
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
			name:   "when the api version is v2",
			fields: fields{version: "2"},
			args: args{
				ctx:     context.Background(),
				payload: &model.FilterPayloadScheme{},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					&model.FilterPayloadScheme{}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/2/filter",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.FilterScheme{}).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the api version is v3",
			fields: fields{version: "3"},
			args: args{
				ctx:     context.Background(),
				payload: &model.FilterPayloadScheme{},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					&model.FilterPayloadScheme{}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/3/filter",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.FilterScheme{}).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the payload is not provided",
			fields: fields{version: "2"},
			args: args{
				ctx:     context.Background(),
				payload: nil,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					(*model.FilterPayloadScheme)(nil)).
					Return(nil, model.ErrNilPayloadError)

				fields.c = client

			},
			wantErr: true,
			Err:     model.ErrNilPayloadError,
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "2"},
			args: args{
				ctx:     context.Background(),
				payload: &model.FilterPayloadScheme{},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					&model.FilterPayloadScheme{}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPost,
					"rest/api/2/filter",
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

			applicationService, err := NewFilterService(testCase.fields.c, testCase.fields.version, testCase.fields.share)
			assert.NoError(t, err)

			gotResult, gotResponse, err := applicationService.Create(testCase.args.ctx, testCase.args.payload)

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

func TestFilterService_Favorite(t *testing.T) {

	type fields struct {
		c       service.Client
		version string
		share   jira.FilterSharingConnector
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
			name:   "when the api version is v2",
			fields: fields{version: "2"},
			args: args{
				ctx: context.Background(),
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/2/filter/favourite",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					[]*model.FilterScheme(nil)).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the api version is v3",
			fields: fields{version: "3"},
			args: args{
				ctx: context.Background(),
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/3/filter/favourite",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					[]*model.FilterScheme(nil)).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the request cannot be created",
			fields: fields{version: "2"},
			args: args{
				ctx: context.Background(),
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/2/filter/favourite",
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

			applicationService, err := NewFilterService(testCase.fields.c, testCase.fields.version, testCase.fields.share)
			assert.NoError(t, err)

			gotResult, gotResponse, err := applicationService.Favorite(testCase.args.ctx)

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

func TestFilterService_My(t *testing.T) {

	type fields struct {
		c       service.Client
		version string
		share   jira.FilterSharingConnector
	}

	type args struct {
		ctx       context.Context
		favorites bool
		expand    []string
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
			name:   "when the api version is v2",
			fields: fields{version: "2"},
			args: args{
				ctx:       context.Background(),
				favorites: true,
				expand:    []string{"subscriptions"},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/2/filter/my?expand=subscriptions&includeFavourites=true",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					[]*model.FilterScheme(nil)).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the api version is v3",
			fields: fields{version: "3"},
			args: args{
				ctx:       context.Background(),
				favorites: true,
				expand:    []string{"subscriptions"},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/3/filter/my?expand=subscriptions&includeFavourites=true",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					[]*model.FilterScheme(nil)).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the request cannot be created",
			fields: fields{version: "2"},
			args: args{
				ctx:       context.Background(),
				favorites: true,
				expand:    []string{"subscriptions"},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/2/filter/my?expand=subscriptions&includeFavourites=true",
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

			applicationService, err := NewFilterService(testCase.fields.c, testCase.fields.version, testCase.fields.share)
			assert.NoError(t, err)

			gotResult, gotResponse, err := applicationService.My(testCase.args.ctx, testCase.args.favorites, testCase.args.expand)

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

func TestFilterService_Search(t *testing.T) {

	type fields struct {
		c       service.Client
		version string
		share   jira.FilterSharingConnector
	}

	type args struct {
		ctx        context.Context
		options    *model.FilterSearchOptionScheme
		startAt    int
		maxResults int
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
			name:   "when the api version is v2",
			fields: fields{version: "2"},
			args: args{
				ctx: context.Background(),
				options: &model.FilterSearchOptionScheme{
					Name:      "filterName",
					AccountID: "owner.accountId",
					Group:     "sharePermissions.group.groupId",
					OrderBy:   "description",
					ProjectID: 100,
					IDs:       []int{10000, 10001},
					Expand:    []string{"description", "viewUrl"},
				},
				startAt:    50,
				maxResults: 100,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/2/filter/search?accountId=owner.accountId&expand=description%2CviewUrl&filterName=filterName&groupname=sharePermissions.group.groupId&id=10000&id=10001&maxResults=100&orderBy=description&projectId=100&startAt=50",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.FilterSearchPageScheme{}).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the api version is v3",
			fields: fields{version: "3"},
			args: args{
				ctx: context.Background(),
				options: &model.FilterSearchOptionScheme{
					Name:      "filterName",
					AccountID: "owner.accountId",
					Group:     "sharePermissions.group.groupId",
					OrderBy:   "description",
					ProjectID: 100,
					IDs:       []int{10000, 10001},
					Expand:    []string{"description", "viewUrl"},
				},
				startAt:    50,
				maxResults: 100,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/3/filter/search?accountId=owner.accountId&expand=description%2CviewUrl&filterName=filterName&groupname=sharePermissions.group.groupId&id=10000&id=10001&maxResults=100&orderBy=description&projectId=100&startAt=50",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.FilterSearchPageScheme{}).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "2"},
			args: args{
				ctx: context.Background(),
				options: &model.FilterSearchOptionScheme{
					Name:      "filterName",
					AccountID: "owner.accountId",
					Group:     "sharePermissions.group.groupId",
					OrderBy:   "description",
					ProjectID: 100,
					IDs:       []int{10000, 10001},
					Expand:    []string{"description", "viewUrl"},
				},
				startAt:    50,
				maxResults: 100,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/2/filter/search?accountId=owner.accountId&expand=description%2CviewUrl&filterName=filterName&groupname=sharePermissions.group.groupId&id=10000&id=10001&maxResults=100&orderBy=description&projectId=100&startAt=50",
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

			applicationService, err := NewFilterService(testCase.fields.c, testCase.fields.version, testCase.fields.share)
			assert.NoError(t, err)

			gotResult, gotResponse, err := applicationService.Search(testCase.args.ctx, testCase.args.options,
				testCase.args.startAt, testCase.args.maxResults)

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

func TestFilterService_Get(t *testing.T) {

	type fields struct {
		c       service.Client
		version string
		share   jira.FilterSharingConnector
	}

	type args struct {
		ctx      context.Context
		filterId int
		expand   []string
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
			name:   "when the api version is v2",
			fields: fields{version: "2"},
			args: args{
				ctx:      context.Background(),
				filterId: 10001,
				expand:   []string{"viewurl"},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/2/filter/10001?expand=viewurl",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.FilterScheme{}).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the api version is v3",
			fields: fields{version: "3"},
			args: args{
				ctx:      context.Background(),
				filterId: 10001,
				expand:   []string{"viewurl"},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/3/filter/10001?expand=viewurl",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.FilterScheme{}).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "2"},
			args: args{
				ctx:      context.Background(),
				filterId: 10001,
				expand:   []string{"viewurl"},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodGet,
					"rest/api/2/filter/10001?expand=viewurl",
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

			applicationService, err := NewFilterService(testCase.fields.c, testCase.fields.version, testCase.fields.share)
			assert.NoError(t, err)

			gotResult, gotResponse, err := applicationService.Get(testCase.args.ctx, testCase.args.filterId, testCase.args.expand)

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

func TestFilterService_Update(t *testing.T) {

	type fields struct {
		c       service.Client
		version string
		share   jira.FilterSharingConnector
	}

	type args struct {
		ctx      context.Context
		filterId int
		payload  *model.FilterPayloadScheme
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
			name:   "when the api version is v2",
			fields: fields{version: "2"},
			args: args{
				ctx:      context.Background(),
				filterId: 10001,
				payload:  &model.FilterPayloadScheme{},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					&model.FilterPayloadScheme{}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPut,
					"rest/api/2/filter/10001",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.FilterScheme{}).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the api version is v3",
			fields: fields{version: "3"},
			args: args{
				ctx:      context.Background(),
				filterId: 10001,
				payload:  &model.FilterPayloadScheme{},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					&model.FilterPayloadScheme{}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPut,
					"rest/api/3/filter/10001",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					&model.FilterScheme{}).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client
			},
		},

		{
			name:   "when the payload is not provided",
			fields: fields{version: "2"},
			args: args{
				ctx:      context.Background(),
				filterId: 10001,
				payload:  nil,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					(*model.FilterPayloadScheme)(nil)).
					Return(nil, model.ErrNilPayloadError)

				fields.c = client

			},
			wantErr: true,
			Err:     model.ErrNilPayloadError,
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "2"},
			args: args{
				ctx:      context.Background(),
				filterId: 10001,
				payload:  &model.FilterPayloadScheme{},
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					&model.FilterPayloadScheme{}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPut,
					"rest/api/2/filter/10001",
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

			applicationService, err := NewFilterService(testCase.fields.c, testCase.fields.version, testCase.fields.share)
			assert.NoError(t, err)

			gotResult, gotResponse, err := applicationService.Update(testCase.args.ctx, testCase.args.filterId, testCase.args.payload)

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

func TestFilterService_Delete(t *testing.T) {

	type fields struct {
		c       service.Client
		version string
		share   jira.FilterSharingConnector
	}

	type args struct {
		ctx      context.Context
		filterId int
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
			name:   "when the api version is v2",
			fields: fields{version: "2"},
			args: args{
				ctx:      context.Background(),
				filterId: 10001,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodDelete,
					"rest/api/2/filter/10001",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					nil).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the api version is v3",
			fields: fields{version: "3"},
			args: args{
				ctx:      context.Background(),
				filterId: 10001,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodDelete,
					"rest/api/3/filter/10001",
					nil).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					nil).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "2"},
			args: args{
				ctx:      context.Background(),
				filterId: 10001,
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("NewRequest",
					context.Background(),
					http.MethodDelete,
					"rest/api/2/filter/10001",
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

			applicationService, err := NewFilterService(testCase.fields.c, testCase.fields.version, testCase.fields.share)
			assert.NoError(t, err)

			gotResponse, err := applicationService.Delete(testCase.args.ctx, testCase.args.filterId)

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

func TestFilterService_Change(t *testing.T) {

	type fields struct {
		c       service.Client
		version string
		share   jira.FilterSharingConnector
	}

	type args struct {
		ctx       context.Context
		filterId  int
		accountId string
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
			name:   "when the api version is v2",
			fields: fields{version: "2"},
			args: args{
				ctx:       context.Background(),
				filterId:  10001,
				accountId: "account-id-sample",
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					struct {
						AccountID string "json:\"accountId\""
					}{AccountID: "account-id-sample"}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPut,
					"rest/api/2/filter/10001/owner",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					nil).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the api version is v3",
			fields: fields{version: "3"},
			args: args{
				ctx:       context.Background(),
				filterId:  10001,
				accountId: "account-id-sample",
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					struct {
						AccountID string "json:\"accountId\""
					}{AccountID: "account-id-sample"}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPut,
					"rest/api/3/filter/10001/owner",
					bytes.NewReader([]byte{})).
					Return(&http.Request{}, nil)

				client.On("Call",
					&http.Request{},
					nil).
					Return(&model.ResponseScheme{}, nil)

				fields.c = client

			},
		},

		{
			name:   "when the http request cannot be created",
			fields: fields{version: "2"},
			args: args{
				ctx:       context.Background(),
				filterId:  10001,
				accountId: "account-id-sample",
			},
			on: func(fields *fields) {

				client := mocks.NewClient(t)

				client.On("TransformStructToReader",
					struct {
						AccountID string "json:\"accountId\""
					}{AccountID: "account-id-sample"}).
					Return(bytes.NewReader([]byte{}), nil)

				client.On("NewRequest",
					context.Background(),
					http.MethodPut,
					"rest/api/2/filter/10001/owner",
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

			applicationService, err := NewFilterService(testCase.fields.c, testCase.fields.version, testCase.fields.share)
			assert.NoError(t, err)

			gotResponse, err := applicationService.Change(testCase.args.ctx, testCase.args.filterId, testCase.args.accountId)

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
