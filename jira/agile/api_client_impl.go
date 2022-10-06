package agile

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/chrisccoy/go-atlassian/jira/agile/internal"
	"github.com/chrisccoy/go-atlassian/pkg/infra/models"
	"github.com/chrisccoy/go-atlassian/service/common"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

func New(httpClient common.HttpClient, site string) (*Client, error) {

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if !strings.HasSuffix(site, "/") {
		site += "/"
	}

	siteAsURL, err := url.Parse(site)
	if err != nil {
		return nil, err
	}

	client := &Client{
		HTTP: httpClient,
		Site: siteAsURL,
	}

	boardService, err := internal.NewBoardService(client, "1.0")
	if err != nil {
		return nil, err
	}

	epicService, err := internal.NewEpicService(client, "1.0")
	if err != nil {
		return nil, err
	}

	sprintService, err := internal.NewSprintService(client, "1.0")
	if err != nil {
		return nil, err
	}

	client.Board = boardService
	client.Epic = epicService
	client.Sprint = sprintService
	client.Auth = internal.NewAuthenticationService(client)

	return client, nil
}

type Client struct {
	HTTP   common.HttpClient
	Site   *url.URL
	Auth   common.Authentication
	Board  *internal.BoardService
	Epic   *internal.EpicService
	Sprint *internal.SprintService
}

func (c *Client) NewFormRequest(ctx context.Context, method, apiEndpoint, contentType string, payload io.Reader) (*http.Request, error) {
	return nil, nil
}

func (c *Client) NewRequest(ctx context.Context, method, apiEndpoint string, payload io.Reader) (*http.Request, error) {

	relativePath, err := url.Parse(apiEndpoint)
	if err != nil {
		return nil, err
	}

	var endpoint = c.Site.ResolveReference(relativePath).String()

	request, err := http.NewRequestWithContext(ctx, method, endpoint, payload)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json")

	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	if c.Auth.HasBasicAuth() {
		request.SetBasicAuth(c.Auth.GetBasicAuth())
	}

	if c.Auth.HasUserAgent() {
		request.Header.Set("User-Agent", c.Auth.GetUserAgent())
	}

	return request, nil
}

func (c *Client) Call(request *http.Request, structure interface{}) (*models.ResponseScheme, error) {

	response, err := c.HTTP.Do(request)
	if err != nil {
		return nil, err
	}

	return c.TransformTheHTTPResponse(response, structure)
}

func (c *Client) TransformTheHTTPResponse(response *http.Response, structure interface{}) (*models.ResponseScheme, error) {

	responseTransformed := &models.ResponseScheme{
		Response: response,
		Code:     response.StatusCode,
		Endpoint: response.Request.URL.String(),
		Method:   response.Request.Method,
	}

	responseAsBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return responseTransformed, err
	}

	responseTransformed.Bytes.Write(responseAsBytes)

	var wasSuccess = response.StatusCode >= 200 && response.StatusCode < 300
	if !wasSuccess {
		return responseTransformed, models.ErrInvalidStatusCodeError
	}

	if structure != nil {
		if err = json.Unmarshal(responseAsBytes, &structure); err != nil {
			return responseTransformed, err
		}
	}

	return responseTransformed, nil
}

func (c *Client) TransformStructToReader(structure interface{}) (io.Reader, error) {

	if structure == nil {
		return nil, models.ErrNilPayloadError
	}

	if reflect.ValueOf(structure).Type().Kind() == reflect.Struct {
		return nil, models.ErrNonPayloadPointerError
	}

	structureAsBodyBytes, err := json.Marshal(structure)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(structureAsBodyBytes), nil
}
