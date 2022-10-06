package confluence

import (
	"context"
	"fmt"
	model "github.com/chrisccoy/go-atlassian/pkg/infra/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestContentService_Archive(t *testing.T) {

	testCases := []struct {
		name               string
		payload            *model.ContentArchivePayloadScheme
		mockFile           string
		wantHTTPMethod     string
		endpoint           string
		context            context.Context
		wantHTTPCodeReturn int
		wantErr            bool
		expectedError      string
	}{
		{
			name: "when the parameters are correct",
			payload: &model.ContentArchivePayloadScheme{
				Pages: []*model.ContentArchiveIDPayloadScheme{
					{
						ID: 1001,
					},
					{
						ID: 1001,
					},
					{
						ID: 1001,
					},
				}},
			wantHTTPMethod:     http.MethodPost,
			mockFile:           "./mocks/archive-content.json",
			endpoint:           "/rest/api/content/archive",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            false,
		},

		{
			name: "when the response body is empty",
			payload: &model.ContentArchivePayloadScheme{
				Pages: []*model.ContentArchiveIDPayloadScheme{
					{
						ID: 1001,
					},
					{
						ID: 1001,
					},
					{
						ID: 1001,
					},
				}},
			wantHTTPMethod:     http.MethodPost,
			mockFile:           "./mocks/empty-json.json",
			endpoint:           "/rest/api/content/archive",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
			expectedError:      "unexpected end of JSON input",
		},

		{
			name:               "when the payload is not provided",
			payload:            nil,
			wantHTTPMethod:     http.MethodPost,
			mockFile:           "./mocks/archive-content.json",
			endpoint:           "/rest/api/content/archive",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
			expectedError:      "failed to parse the interface pointer, please provide a valid one",
		},

		{
			name: "when the response status is not valid",
			payload: &model.ContentArchivePayloadScheme{
				Pages: []*model.ContentArchiveIDPayloadScheme{
					{
						ID: 1001,
					},
					{
						ID: 1001,
					},
					{
						ID: 1001,
					},
				}},
			wantHTTPMethod:     http.MethodPost,
			endpoint:           "/rest/api/content/archive",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusBadRequest,
			wantErr:            true,
			expectedError:      "unexpected end of JSON input",
		},

		{
			name: "when the context is not provided",
			payload: &model.ContentArchivePayloadScheme{
				Pages: []*model.ContentArchiveIDPayloadScheme{
					{
						ID: 1001,
					},
					{
						ID: 1001,
					},
					{
						ID: 1001,
					},
				}},
			wantHTTPMethod:     http.MethodPost,
			endpoint:           "/rest/api/content/archive",
			context:            nil,
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
			expectedError:      "request creation failed: net/http: nil Context",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			//Init a new HTTP mock server
			mockOptions := mockServerOptions{
				Endpoint:           testCase.endpoint,
				MockFilePath:       testCase.mockFile,
				MethodAccepted:     testCase.wantHTTPMethod,
				ResponseCodeWanted: testCase.wantHTTPCodeReturn,
			}

			mockServer, err := startMockServer(&mockOptions)
			if err != nil {
				t.Fatal(err)
			}

			defer mockServer.Close()

			//Init the library instance
			mockClient, err := startMockClient(mockServer.URL)
			if err != nil {
				t.Fatal(err)
			}

			implementation := &ContentService{client: mockClient}

			gotResult, gotResponse, err := implementation.Archive(testCase.context, testCase.payload)

			if testCase.wantErr {

				if err != nil {
					t.Logf("error returned: %v", err.Error())
				}

				assert.EqualError(t, err, testCase.expectedError)

				if gotResponse != nil {
					t.Logf("HTTP Code Wanted: %v, HTTP Code Returned: %v", testCase.wantHTTPCodeReturn, gotResponse.Code)
				}

			} else {

				assert.NoError(t, err)
				assert.NotEqual(t, gotResponse, nil)
				assert.NotEqual(t, gotResult, nil)

				apiEndpoint, err := url.Parse(gotResponse.Endpoint)
				if err != nil {
					t.Fatal(err)
				}

				var endpointToAssert string

				if apiEndpoint.Query().Encode() != "" {
					endpointToAssert = fmt.Sprintf("%v?%v", apiEndpoint.Path, apiEndpoint.Query().Encode())
				} else {
					endpointToAssert = apiEndpoint.Path
				}

				t.Logf("HTTP Endpoint Wanted: %v, HTTP Endpoint Returned: %v", testCase.endpoint, endpointToAssert)
				assert.Equal(t, testCase.endpoint, endpointToAssert)
			}
		})
	}
}

func TestContentService_Create(t *testing.T) {

	testCases := []struct {
		name               string
		payload            *model.ContentScheme
		mockFile           string
		wantHTTPMethod     string
		endpoint           string
		context            context.Context
		wantHTTPCodeReturn int
		wantErr            bool
	}{
		{
			name: "CreateContentWhenTheParametersAreCorrect",
			payload: &model.ContentScheme{
				Type:  "page", // Valid values: page, blogpost, comment
				Title: "Confluence Page Title",
				Space: &model.SpaceScheme{Key: "DUMMY"},
				Body: &model.BodyScheme{
					Storage: &model.BodyNodeScheme{
						Value:          "<p>This is <br/> a new page</p>",
						Representation: "storage",
					},
				},
			},
			mockFile:           "./mocks/create-content.json",
			wantHTTPMethod:     http.MethodPost,
			endpoint:           "/rest/api/content",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusCreated,
			wantErr:            false,
		},

		{
			name:               "CreateContentWhenThePayloadIsNotProvided",
			payload:            nil,
			mockFile:           "./mocks/create-content.json",
			wantHTTPMethod:     http.MethodPost,
			endpoint:           "/rest/api/content",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusCreated,
			wantErr:            true,
		},

		{
			name: "CreateContentWhenTheRequestMethodIsIncorrect",
			payload: &model.ContentScheme{
				Type:  "page", // Valid values: page, blogpost, comment
				Title: "Confluence Page Title",
				Space: &model.SpaceScheme{Key: "DUMMY"},
				Body: &model.BodyScheme{
					Storage: &model.BodyNodeScheme{
						Value:          "<p>This is <br/> a new page</p>",
						Representation: "storage",
					},
				},
			},
			mockFile:           "./mocks/create-content.json",
			wantHTTPMethod:     http.MethodPut,
			endpoint:           "/rest/api/content",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusCreated,
			wantErr:            true,
		},

		{
			name: "CreateContentWhenTheContextIsNotProvided",
			payload: &model.ContentScheme{
				Type:  "page", // Valid values: page, blogpost, comment
				Title: "Confluence Page Title",
				Space: &model.SpaceScheme{Key: "DUMMY"},
				Body: &model.BodyScheme{
					Storage: &model.BodyNodeScheme{
						Value:          "<p>This is <br/> a new page</p>",
						Representation: "storage",
					},
				},
			},
			mockFile:           "./mocks/create-content.json",
			wantHTTPMethod:     http.MethodPost,
			endpoint:           "/rest/api/content",
			context:            nil,
			wantHTTPCodeReturn: http.StatusCreated,
			wantErr:            true,
		},

		{
			name: "CreateContentWhenTheRequestBodyIsEmpty",
			payload: &model.ContentScheme{
				Type:  "page", // Valid values: page, blogpost, comment
				Title: "Confluence Page Title",
				Space: &model.SpaceScheme{Key: "DUMMY"},
				Body: &model.BodyScheme{
					Storage: &model.BodyNodeScheme{
						Value:          "<p>This is <br/> a new page</p>",
						Representation: "storage",
					},
				},
			},
			mockFile:           "./mocks/empty-json.json",
			wantHTTPMethod:     http.MethodPost,
			endpoint:           "/rest/api/content",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusCreated,
			wantErr:            true,
		},
	}

	for _, testCase := range testCases {

		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			//Init a new HTTP mock server
			mockOptions := mockServerOptions{
				Endpoint:           testCase.endpoint,
				MockFilePath:       testCase.mockFile,
				MethodAccepted:     testCase.wantHTTPMethod,
				ResponseCodeWanted: testCase.wantHTTPCodeReturn,
			}

			mockServer, err := startMockServer(&mockOptions)
			if err != nil {
				t.Fatal(err)
			}

			defer mockServer.Close()

			//Init the library instance
			mockClient, err := startMockClient(mockServer.URL)
			if err != nil {
				t.Fatal(err)
			}

			service := &ContentService{client: mockClient}

			gotResult, gotResponse, err := service.Create(testCase.context, testCase.payload)

			if testCase.wantErr {

				if err != nil {
					t.Logf("error returned: %v", err.Error())
				}

				assert.Error(t, err)
			} else {

				assert.NoError(t, err)
				assert.NotEqual(t, gotResponse, nil)
				assert.NotEqual(t, gotResult, nil)

				apiEndpoint, err := url.Parse(gotResponse.Endpoint)
				if err != nil {
					t.Fatal(err)
				}

				var endpointToAssert string

				if apiEndpoint.Query().Encode() != "" {
					endpointToAssert = fmt.Sprintf("%v?%v", apiEndpoint.Path, apiEndpoint.Query().Encode())
				} else {
					endpointToAssert = apiEndpoint.Path
				}

				t.Logf("HTTP Endpoint Wanted: %v, HTTP Endpoint Returned: %v", testCase.endpoint, endpointToAssert)
				assert.Equal(t, testCase.endpoint, endpointToAssert)
			}

		})

	}

}

func TestContentService_Gets(t *testing.T) {

	testCases := []struct {
		name                string
		options             *model.GetContentOptionsScheme
		startAt, maxResults int
		mockFile            string
		wantHTTPMethod      string
		endpoint            string
		context             context.Context
		wantHTTPCodeReturn  int
		wantErr             bool
	}{
		{
			name: "GetsContentWhenTheAllOptionsAreProvided",
			options: &model.GetContentOptionsScheme{
				ContextType: "page",
				SpaceKey:    "DUMMY",
				Title:       "*page*",
				Trigger:     "trigger-sample",
				OrderBy:     "id",
				Status:      []string{"status-00", "status-01"},
				Expand:      []string{"all"},
				PostingDay:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			startAt:            0,
			maxResults:         50,
			mockFile:           "./mocks/get-contents.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content?expand=all&limit=50&orderby=id&postingDay=2019-11-17&spaceKey=DUMMY&start=0&status=status-00%2Cstatus-01&title=%2Apage%2A&trigger=trigger-sample&type=page",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            false,
		},

		{
			name: "GetsContentWhenTheContextIsNotSet",
			options: &model.GetContentOptionsScheme{
				ContextType: "page",
				SpaceKey:    "DUMMY",
				Title:       "*page*",
				Trigger:     "trigger-sample",
				OrderBy:     "id",
				Status:      []string{"status-00", "status-01"},
				Expand:      []string{"all"},
				PostingDay:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			startAt:            0,
			maxResults:         50,
			mockFile:           "./mocks/get-contents.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content?expand=all&limit=50&orderby=id&postingDay=2019-11-17&spaceKey=DUMMY&start=0&status=status-00%2Cstatus-01&title=%2Apage%2A&trigger=trigger-sample&type=page",
			context:            nil,
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name: "GetsContentWhenTheRequestMethodIsIncorrect",
			options: &model.GetContentOptionsScheme{
				ContextType: "page",
				SpaceKey:    "DUMMY",
				Title:       "*page*",
				Trigger:     "trigger-sample",
				OrderBy:     "id",
				Status:      []string{"status-00", "status-01"},
				Expand:      []string{"all"},
				PostingDay:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			startAt:            0,
			maxResults:         50,
			mockFile:           "./mocks/get-contents.json",
			wantHTTPMethod:     http.MethodPut,
			endpoint:           "/rest/api/content?expand=all&limit=50&orderby=id&postingDay=2019-11-17&spaceKey=DUMMY&start=0&status=status-00%2Cstatus-01&title=%2Apage%2A&trigger=trigger-sample&type=page",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name: "GetsContentWhenTheResponseStatusIsIncorrect",
			options: &model.GetContentOptionsScheme{
				ContextType: "page",
				SpaceKey:    "DUMMY",
				Title:       "*page*",
				Trigger:     "trigger-sample",
				OrderBy:     "id",
				Status:      []string{"status-00", "status-01"},
				Expand:      []string{"all"},
				PostingDay:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			startAt:            0,
			maxResults:         50,
			mockFile:           "./mocks/get-contents.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content?expand=all&limit=50&orderby=id&postingDay=2019-11-17&spaceKey=DUMMY&start=0&status=status-00%2Cstatus-01&title=%2Apage%2A&trigger=trigger-sample&type=page",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusBadRequest,
			wantErr:            true,
		},

		{
			name: "GetsContentWhenTheResponseBodyIsEmpty",
			options: &model.GetContentOptionsScheme{
				ContextType: "page",
				SpaceKey:    "DUMMY",
				Title:       "*page*",
				Trigger:     "trigger-sample",
				OrderBy:     "id",
				Status:      []string{"status-00", "status-01"},
				Expand:      []string{"all"},
				PostingDay:  time.Date(2019, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
			startAt:            0,
			maxResults:         50,
			mockFile:           "./mocks/empty-json.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content?expand=all&limit=50&orderby=id&postingDay=2019-11-17&spaceKey=DUMMY&start=0&status=status-00%2Cstatus-01&title=%2Apage%2A&trigger=trigger-sample&type=page",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},
	}

	for _, testCase := range testCases {

		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			//Init a new HTTP mock server
			mockOptions := mockServerOptions{
				Endpoint:           testCase.endpoint,
				MockFilePath:       testCase.mockFile,
				MethodAccepted:     testCase.wantHTTPMethod,
				ResponseCodeWanted: testCase.wantHTTPCodeReturn,
			}

			mockServer, err := startMockServer(&mockOptions)
			if err != nil {
				t.Fatal(err)
			}

			defer mockServer.Close()

			//Init the library instance
			mockClient, err := startMockClient(mockServer.URL)
			if err != nil {
				t.Fatal(err)
			}

			service := &ContentService{client: mockClient}

			gotResult, gotResponse, err := service.Gets(
				testCase.context,
				testCase.options,
				testCase.startAt,
				testCase.maxResults,
			)

			if testCase.wantErr {

				if err != nil {
					t.Logf("error returned: %v", err.Error())
				}

				assert.Error(t, err)
			} else {

				assert.NoError(t, err)
				assert.NotEqual(t, gotResponse, nil)
				assert.NotEqual(t, gotResult, nil)

				apiEndpoint, err := url.Parse(gotResponse.Endpoint)
				if err != nil {
					t.Fatal(err)
				}

				var endpointToAssert string

				if apiEndpoint.Query().Encode() != "" {
					endpointToAssert = fmt.Sprintf("%v?%v", apiEndpoint.Path, apiEndpoint.Query().Encode())
				} else {
					endpointToAssert = apiEndpoint.Path
				}

				t.Logf("HTTP Endpoint Wanted: %v, HTTP Endpoint Returned: %v", testCase.endpoint, endpointToAssert)
				assert.Equal(t, testCase.endpoint, endpointToAssert)
			}

		})

	}

}

func TestContentService_Search(t *testing.T) {

	testCases := []struct {
		name               string
		cql, cqlContext    string
		expand             []string
		cursor             string
		maxResults         int
		mockFile           string
		wantHTTPMethod     string
		endpoint           string
		context            context.Context
		wantHTTPCodeReturn int
		wantErr            bool
	}{
		{
			name:               "SearchContentsWhenTheParametersAreCorrect",
			cql:                "type=page",
			cqlContext:         "DUMMY",
			expand:             []string{"space", "metadata.labels"},
			cursor:             "ca470fa51eb0",
			maxResults:         50,
			mockFile:           "./mocks/get-contents.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/search?cql=type%3Dpage&cqlcontext=DUMMY&cursor=ca470fa51eb0&expand=space%2Cmetadata.labels&limit=50",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            false,
		},

		{
			name:               "SearchContentsWhenTheCQLIsNotProvided",
			cql:                "",
			cqlContext:         "DUMMY",
			expand:             []string{"space", "metadata.labels"},
			cursor:             "ca470fa51eb0",
			maxResults:         50,
			mockFile:           "./mocks/get-contents.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/search?cql=type%3Dpage&cqlcontext=DUMMY&cursor=ca470fa51eb0&expand=space%2Cmetadata.labels&limit=50",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name:               "SearchContentsWhenTheContextIsNil",
			cql:                "type=page",
			cqlContext:         "DUMMY",
			expand:             []string{"space", "metadata.labels"},
			cursor:             "ca470fa51eb0",
			maxResults:         50,
			mockFile:           "./mocks/get-contents.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/search?cql=type%3Dpage&cqlcontext=DUMMY&cursor=ca470fa51eb0&expand=space%2Cmetadata.labels&limit=50",
			context:            nil,
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name:               "SearchContentsWhenTheRequestMethodIsIncorrect",
			cql:                "type=page",
			cqlContext:         "DUMMY",
			expand:             []string{"space", "metadata.labels"},
			cursor:             "ca470fa51eb0",
			maxResults:         50,
			mockFile:           "./mocks/get-contents.json",
			wantHTTPMethod:     http.MethodPost,
			endpoint:           "/rest/api/content/search?cql=type%3Dpage&cqlcontext=DUMMY&cursor=ca470fa51eb0&expand=space%2Cmetadata.labels&limit=50",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name:               "SearchContentsWhenTheStatusCodeIsIncorrect",
			cql:                "type=page",
			cqlContext:         "DUMMY",
			expand:             []string{"space", "metadata.labels"},
			cursor:             "ca470fa51eb0",
			maxResults:         50,
			mockFile:           "./mocks/get-contents.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/search?cql=type%3Dpage&cqlcontext=DUMMY&cursor=ca470fa51eb0&expand=space%2Cmetadata.labels&limit=50",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusBadRequest,
			wantErr:            true,
		},

		{
			name:               "SearchContentsWhenTheRequestBodyIsEmpty",
			cql:                "type=page",
			cqlContext:         "DUMMY",
			expand:             []string{"space", "metadata.labels"},
			cursor:             "ca470fa51eb0",
			maxResults:         50,
			mockFile:           "./mocks/empty-json.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/search?cql=type%3Dpage&cqlcontext=DUMMY&cursor=ca470fa51eb0&expand=space%2Cmetadata.labels&limit=50",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},
	}

	for _, testCase := range testCases {

		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			//Init a new HTTP mock server
			mockOptions := mockServerOptions{
				Endpoint:           testCase.endpoint,
				MockFilePath:       testCase.mockFile,
				MethodAccepted:     testCase.wantHTTPMethod,
				ResponseCodeWanted: testCase.wantHTTPCodeReturn,
			}

			mockServer, err := startMockServer(&mockOptions)
			if err != nil {
				t.Fatal(err)
			}

			defer mockServer.Close()

			//Init the library instance
			mockClient, err := startMockClient(mockServer.URL)
			if err != nil {
				t.Fatal(err)
			}

			service := &ContentService{client: mockClient}

			gotResult, gotResponse, err := service.Search(testCase.context, testCase.cql, testCase.cqlContext,
				testCase.expand, testCase.cursor, testCase.maxResults)

			if testCase.wantErr {

				if err != nil {
					t.Logf("error returned: %v", err.Error())
				}

				assert.Error(t, err)
			} else {

				assert.NoError(t, err)
				assert.NotEqual(t, gotResponse, nil)
				assert.NotEqual(t, gotResult, nil)

				apiEndpoint, err := url.Parse(gotResponse.Endpoint)
				if err != nil {
					t.Fatal(err)
				}

				var endpointToAssert string

				if apiEndpoint.Query().Encode() != "" {
					endpointToAssert = fmt.Sprintf("%v?%v", apiEndpoint.Path, apiEndpoint.Query().Encode())
				} else {
					endpointToAssert = apiEndpoint.Path
				}

				t.Logf("HTTP Endpoint Wanted: %v, HTTP Endpoint Returned: %v", testCase.endpoint, endpointToAssert)
				assert.Equal(t, testCase.endpoint, endpointToAssert)
			}

		})

	}

}

func TestContentService_Get(t *testing.T) {

	testCases := []struct {
		name               string
		contentID          string
		expand             []string
		version            int
		mockFile           string
		wantHTTPMethod     string
		endpoint           string
		context            context.Context
		wantHTTPCodeReturn int
		wantErr            bool
	}{
		{
			name:               "GetContentWhenTheParametersAreCorrect",
			contentID:          "64290828",
			expand:             []string{"any"},
			version:            1,
			mockFile:           "./mocks/get-content.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/64290828?expand=any&version=1",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            false,
		},

		{
			name:               "GetContentWhenTheContentIDIsNotProvided",
			contentID:          "",
			expand:             []string{"any"},
			version:            1,
			mockFile:           "./mocks/get-content.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/64290828?expand=any&version=1",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name:               "GetContentWhenTheContextIsNotProvided",
			contentID:          "64290828",
			expand:             []string{"any"},
			version:            1,
			mockFile:           "./mocks/get-content.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/64290828?expand=any&version=1",
			context:            nil,
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name:               "GetContentWhenTheRequestMethodIsIncorrect",
			contentID:          "64290828",
			expand:             []string{"any"},
			version:            1,
			mockFile:           "./mocks/get-content.json",
			wantHTTPMethod:     http.MethodHead,
			endpoint:           "/rest/api/content/64290828?expand=any&version=1",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name:               "GetContentWhenTheStatusCodeIsIncorrect",
			contentID:          "64290828",
			expand:             []string{"any"},
			version:            1,
			mockFile:           "./mocks/get-content.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/64290828?expand=any&version=1",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusBadRequest,
			wantErr:            true,
		},

		{
			name:               "GetContentWhenTheResponseBodyIsEmpty",
			contentID:          "64290828",
			expand:             []string{"any"},
			version:            1,
			mockFile:           "./mocks/empty-json.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/64290828?expand=any&version=1",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},
	}

	for _, testCase := range testCases {

		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			//Init a new HTTP mock server
			mockOptions := mockServerOptions{
				Endpoint:           testCase.endpoint,
				MockFilePath:       testCase.mockFile,
				MethodAccepted:     testCase.wantHTTPMethod,
				ResponseCodeWanted: testCase.wantHTTPCodeReturn,
			}

			mockServer, err := startMockServer(&mockOptions)
			if err != nil {
				t.Fatal(err)
			}

			defer mockServer.Close()

			//Init the library instance
			mockClient, err := startMockClient(mockServer.URL)
			if err != nil {
				t.Fatal(err)
			}

			service := &ContentService{client: mockClient}

			gotResult, gotResponse, err := service.Get(testCase.context, testCase.contentID, testCase.expand, testCase.version)

			if testCase.wantErr {

				if err != nil {
					t.Logf("error returned: %v", err.Error())
				}

				assert.Error(t, err)
			} else {

				assert.NoError(t, err)
				assert.NotEqual(t, gotResponse, nil)
				assert.NotEqual(t, gotResult, nil)

				apiEndpoint, err := url.Parse(gotResponse.Endpoint)
				if err != nil {
					t.Fatal(err)
				}

				var endpointToAssert string

				if apiEndpoint.Query().Encode() != "" {
					endpointToAssert = fmt.Sprintf("%v?%v", apiEndpoint.Path, apiEndpoint.Query().Encode())
				} else {
					endpointToAssert = apiEndpoint.Path
				}

				t.Logf("HTTP Endpoint Wanted: %v, HTTP Endpoint Returned: %v", testCase.endpoint, endpointToAssert)
				assert.Equal(t, testCase.endpoint, endpointToAssert)
			}

		})

	}

}

func TestContentService_Update(t *testing.T) {

	testCases := []struct {
		name               string
		contentID          string
		payload            *model.ContentScheme
		mockFile           string
		wantHTTPMethod     string
		endpoint           string
		context            context.Context
		wantHTTPCodeReturn int
		wantErr            bool
	}{
		{
			name:      "UpdateContentWhenTheParametersAreCorrect",
			contentID: "2939332",
			payload: &model.ContentScheme{
				Type:  "page", // Valid values: page, blogpost, comment
				Title: "Confluence Page Title",
				Space: &model.SpaceScheme{Key: "DUMMY"},
				Body: &model.BodyScheme{
					Storage: &model.BodyNodeScheme{
						Value:          "<p>This is <br/> a new page</p>",
						Representation: "storage",
					},
				},
				Version: &model.ContentVersionScheme{Number: 2},
			},
			mockFile:           "./mocks/get-content.json",
			wantHTTPMethod:     http.MethodPut,
			endpoint:           "/rest/api/content/2939332",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            false,
		},

		{
			name:      "UpdateContentWhenTheContentIDIsNotProvided",
			contentID: "",
			payload: &model.ContentScheme{
				Type:  "page", // Valid values: page, blogpost, comment
				Title: "Confluence Page Title",
				Space: &model.SpaceScheme{Key: "DUMMY"},
				Body: &model.BodyScheme{
					Storage: &model.BodyNodeScheme{
						Value:          "<p>This is <br/> a new page</p>",
						Representation: "storage",
					},
				},
				Version: &model.ContentVersionScheme{Number: 2},
			},
			mockFile:           "./mocks/get-content.json",
			wantHTTPMethod:     http.MethodPut,
			endpoint:           "/rest/api/content/2939332",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name:               "UpdateContentWhenThePayloadIsNotProvided",
			contentID:          "2939332",
			payload:            nil,
			mockFile:           "./mocks/get-content.json",
			wantHTTPMethod:     http.MethodPut,
			endpoint:           "/rest/api/content/2939332",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name:      "UpdateContentWhenTheContentIsNotProvided",
			contentID: "2939332",
			payload: &model.ContentScheme{
				Type:  "page", // Valid values: page, blogpost, comment
				Title: "Confluence Page Title",
				Space: &model.SpaceScheme{Key: "DUMMY"},
				Body: &model.BodyScheme{
					Storage: &model.BodyNodeScheme{
						Value:          "<p>This is <br/> a new page</p>",
						Representation: "storage",
					},
				},
				Version: &model.ContentVersionScheme{Number: 2},
			},
			mockFile:           "./mocks/get-content.json",
			wantHTTPMethod:     http.MethodPut,
			endpoint:           "/rest/api/content/2939332",
			context:            nil,
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name:      "UpdateContentWhenTheRequestMethodIsIncorrect",
			contentID: "2939332",
			payload: &model.ContentScheme{
				Type:  "page", // Valid values: page, blogpost, comment
				Title: "Confluence Page Title",
				Space: &model.SpaceScheme{Key: "DUMMY"},
				Body: &model.BodyScheme{
					Storage: &model.BodyNodeScheme{
						Value:          "<p>This is <br/> a new page</p>",
						Representation: "storage",
					},
				},
				Version: &model.ContentVersionScheme{Number: 2},
			},
			mockFile:           "./mocks/get-content.json",
			wantHTTPMethod:     http.MethodHead,
			endpoint:           "/rest/api/content/2939332",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name:      "UpdateContentWhenTheStatusCodeIsIncorrect",
			contentID: "2939332",
			payload: &model.ContentScheme{
				Type:  "page", // Valid values: page, blogpost, comment
				Title: "Confluence Page Title",
				Space: &model.SpaceScheme{Key: "DUMMY"},
				Body: &model.BodyScheme{
					Storage: &model.BodyNodeScheme{
						Value:          "<p>This is <br/> a new page</p>",
						Representation: "storage",
					},
				},
				Version: &model.ContentVersionScheme{Number: 2},
			},
			mockFile:           "./mocks/get-content.json",
			wantHTTPMethod:     http.MethodPut,
			endpoint:           "/rest/api/content/2939332",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusBadRequest,
			wantErr:            true,
		},

		{
			name:      "UpdateContentWhenTheResponseBodyIsEmpty",
			contentID: "2939332",
			payload: &model.ContentScheme{
				Type:  "page", // Valid values: page, blogpost, comment
				Title: "Confluence Page Title",
				Space: &model.SpaceScheme{Key: "DUMMY"},
				Body: &model.BodyScheme{
					Storage: &model.BodyNodeScheme{
						Value:          "<p>This is <br/> a new page</p>",
						Representation: "storage",
					},
				},
				Version: &model.ContentVersionScheme{Number: 2},
			},
			mockFile:           "./mocks/empty-json.json",
			wantHTTPMethod:     http.MethodPut,
			endpoint:           "/rest/api/content/2939332",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},
	}

	for _, testCase := range testCases {

		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			//Init a new HTTP mock server
			mockOptions := mockServerOptions{
				Endpoint:           testCase.endpoint,
				MockFilePath:       testCase.mockFile,
				MethodAccepted:     testCase.wantHTTPMethod,
				ResponseCodeWanted: testCase.wantHTTPCodeReturn,
			}

			mockServer, err := startMockServer(&mockOptions)
			if err != nil {
				t.Fatal(err)
			}

			defer mockServer.Close()

			//Init the library instance
			mockClient, err := startMockClient(mockServer.URL)
			if err != nil {
				t.Fatal(err)
			}

			service := &ContentService{client: mockClient}

			gotResult, gotResponse, err := service.Update(testCase.context, testCase.contentID, testCase.payload)

			if testCase.wantErr {

				if err != nil {
					t.Logf("error returned: %v", err.Error())
				}

				assert.Error(t, err)
			} else {

				assert.NoError(t, err)
				assert.NotEqual(t, gotResponse, nil)
				assert.NotEqual(t, gotResult, nil)

				apiEndpoint, err := url.Parse(gotResponse.Endpoint)
				if err != nil {
					t.Fatal(err)
				}

				var endpointToAssert string

				if apiEndpoint.Query().Encode() != "" {
					endpointToAssert = fmt.Sprintf("%v?%v", apiEndpoint.Path, apiEndpoint.Query().Encode())
				} else {
					endpointToAssert = apiEndpoint.Path
				}

				t.Logf("HTTP Endpoint Wanted: %v, HTTP Endpoint Returned: %v", testCase.endpoint, endpointToAssert)
				assert.Equal(t, testCase.endpoint, endpointToAssert)
			}

		})

	}

}

func TestContentService_Delete(t *testing.T) {

	testCases := []struct {
		name               string
		contentID          string
		status             string
		mockFile           string
		wantHTTPMethod     string
		endpoint           string
		context            context.Context
		wantHTTPCodeReturn int
		wantErr            bool
	}{
		{
			name:               "DeleteContentWhenTheParametersAreCorrect",
			contentID:          "34325949",
			status:             "trash",
			mockFile:           "",
			wantHTTPMethod:     http.MethodDelete,
			endpoint:           "/rest/api/content/34325949?status=trash",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusNoContent,
			wantErr:            false,
		},

		{
			name:               "DeleteContentWhenTheContentIDIsNotProvided",
			contentID:          "",
			status:             "trash",
			mockFile:           "",
			wantHTTPMethod:     http.MethodDelete,
			endpoint:           "/rest/api/content/34325949?status=trash",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusNoContent,
			wantErr:            true,
		},

		{
			name:               "DeleteContentWhenTheStatusIsNotProvided",
			contentID:          "34325949",
			status:             "",
			mockFile:           "",
			wantHTTPMethod:     http.MethodDelete,
			endpoint:           "/rest/api/content/34325949",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusNoContent,
			wantErr:            false,
		},

		{
			name:               "DeleteContentWhenTheContextIsNotProvided",
			contentID:          "34325949",
			status:             "trash",
			mockFile:           "",
			wantHTTPMethod:     http.MethodDelete,
			endpoint:           "/rest/api/content/34325949?status=trash",
			context:            nil,
			wantHTTPCodeReturn: http.StatusNoContent,
			wantErr:            true,
		},

		{
			name:               "DeleteContentWhenTheRequestMethodIsIncorrect",
			contentID:          "34325949",
			status:             "trash",
			mockFile:           "",
			wantHTTPMethod:     http.MethodPost,
			endpoint:           "/rest/api/content/34325949?status=trash",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusNoContent,
			wantErr:            true,
		},

		{
			name:               "DeleteContentWhenTheStatusCodeIsIncorrect",
			contentID:          "34325949",
			status:             "trash",
			mockFile:           "",
			wantHTTPMethod:     http.MethodDelete,
			endpoint:           "/rest/api/content/34325949?status=trash",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusBadRequest,
			wantErr:            true,
		},
	}

	for _, testCase := range testCases {

		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			//Init a new HTTP mock server
			mockOptions := mockServerOptions{
				Endpoint:           testCase.endpoint,
				MockFilePath:       testCase.mockFile,
				MethodAccepted:     testCase.wantHTTPMethod,
				ResponseCodeWanted: testCase.wantHTTPCodeReturn,
			}

			mockServer, err := startMockServer(&mockOptions)
			if err != nil {
				t.Fatal(err)
			}

			defer mockServer.Close()

			//Init the library instance
			mockClient, err := startMockClient(mockServer.URL)
			if err != nil {
				t.Fatal(err)
			}

			service := &ContentService{client: mockClient}

			gotResponse, err := service.Delete(testCase.context, testCase.contentID, testCase.status)

			if testCase.wantErr {

				if err != nil {
					t.Logf("error returned: %v", err.Error())
				}

				assert.Error(t, err)
			} else {

				assert.NoError(t, err)
				assert.NotEqual(t, gotResponse, nil)

				apiEndpoint, err := url.Parse(gotResponse.Endpoint)
				if err != nil {
					t.Fatal(err)
				}

				var endpointToAssert string

				if apiEndpoint.Query().Encode() != "" {
					endpointToAssert = fmt.Sprintf("%v?%v", apiEndpoint.Path, apiEndpoint.Query().Encode())
				} else {
					endpointToAssert = apiEndpoint.Path
				}

				t.Logf("HTTP Endpoint Wanted: %v, HTTP Endpoint Returned: %v", testCase.endpoint, endpointToAssert)
				assert.Equal(t, testCase.endpoint, endpointToAssert)
			}

		})

	}

}

func TestContentService_History(t *testing.T) {

	testCases := []struct {
		name               string
		contentID          string
		expand             []string
		mockFile           string
		wantHTTPMethod     string
		endpoint           string
		context            context.Context
		wantHTTPCodeReturn int
		wantErr            bool
	}{
		{
			name:               "GetContentHistoryWhenTheParametersAreCorrect",
			contentID:          "339sa9as9s",
			expand:             []string{"nextVersion", "lastVersion"},
			mockFile:           "./mocks/get-content-history.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/339sa9as9s/history?expand=nextVersion%2ClastVersion",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            false,
		},

		{
			name:               "GetContentHistoryWhenTheContextIsNotProvided",
			contentID:          "339sa9as9s",
			expand:             []string{"nextVersion", "lastVersion"},
			mockFile:           "./mocks/get-content-history.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/339sa9as9s/history?expand=nextVersion%2ClastVersion",
			context:            nil,
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name:               "GetContentHistoryWhenTheContentIDIsNotProvided",
			contentID:          "",
			expand:             []string{"nextVersion", "lastVersion"},
			mockFile:           "./mocks/get-content-history.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/339sa9as9s/history?expand=nextVersion%2ClastVersion",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name:               "GetContentHistoryWhenTheExpandsAreNotProvided",
			contentID:          "339sa9as9s",
			expand:             nil,
			mockFile:           "./mocks/get-content-history.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/339sa9as9s/history",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            false,
		},

		{
			name:               "GetContentHistoryWhenTheRequestMethodIsIncorrect",
			contentID:          "339sa9as9s",
			expand:             []string{"nextVersion", "lastVersion"},
			mockFile:           "./mocks/get-content-history.json",
			wantHTTPMethod:     http.MethodPut,
			endpoint:           "/rest/api/content/339sa9as9s/history?expand=nextVersion%2ClastVersion",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},

		{
			name:               "GetContentHistoryWhenTheStatusCodeIsIncorrect",
			contentID:          "339sa9as9s",
			expand:             []string{"nextVersion", "lastVersion"},
			mockFile:           "./mocks/get-content-history.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/339sa9as9s/history?expand=nextVersion%2ClastVersion",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusBadRequest,
			wantErr:            true,
		},

		{
			name:               "GetContentHistoryWhenTheResponseBodyIsEmpty",
			contentID:          "339sa9as9s",
			expand:             []string{"nextVersion", "lastVersion"},
			mockFile:           "./mocks/empty-json.json",
			wantHTTPMethod:     http.MethodGet,
			endpoint:           "/rest/api/content/339sa9as9s/history?expand=nextVersion%2ClastVersion",
			context:            context.Background(),
			wantHTTPCodeReturn: http.StatusOK,
			wantErr:            true,
		},
	}

	for _, testCase := range testCases {

		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			//Init a new HTTP mock server
			mockOptions := mockServerOptions{
				Endpoint:           testCase.endpoint,
				MockFilePath:       testCase.mockFile,
				MethodAccepted:     testCase.wantHTTPMethod,
				ResponseCodeWanted: testCase.wantHTTPCodeReturn,
			}

			mockServer, err := startMockServer(&mockOptions)
			if err != nil {
				t.Fatal(err)
			}

			defer mockServer.Close()

			//Init the library instance
			mockClient, err := startMockClient(mockServer.URL)
			if err != nil {
				t.Fatal(err)
			}

			service := &ContentService{client: mockClient}

			gotResult, gotResponse, err := service.History(testCase.context, testCase.contentID, testCase.expand)

			if testCase.wantErr {

				if err != nil {
					t.Logf("error returned: %v", err.Error())
				}

				assert.Error(t, err)
			} else {

				assert.NoError(t, err)
				assert.NotEqual(t, gotResponse, nil)
				assert.NotEqual(t, gotResult, nil)

				apiEndpoint, err := url.Parse(gotResponse.Endpoint)
				if err != nil {
					t.Fatal(err)
				}

				var endpointToAssert string

				if apiEndpoint.Query().Encode() != "" {
					endpointToAssert = fmt.Sprintf("%v?%v", apiEndpoint.Path, apiEndpoint.Query().Encode())
				} else {
					endpointToAssert = apiEndpoint.Path
				}

				t.Logf("HTTP Endpoint Wanted: %v, HTTP Endpoint Returned: %v", testCase.endpoint, endpointToAssert)
				assert.Equal(t, testCase.endpoint, endpointToAssert)
			}

		})

	}

}
