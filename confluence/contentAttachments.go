package confluence

import (
	"bytes"
	"context"
	"fmt"
	model "github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ContentAttachmentService struct {
	client *Client
}

// Gets returns the attachments for a piece of content.
// By default, the following objects are expanded: metadata.
// Docs: https://docs.go-atlassian.io/confluence-cloud/content/attachments#get-attachments
func (c *ContentAttachmentService) Gets(ctx context.Context, contentID string, startAt, maxResults int,
	options *model.GetContentAttachmentsOptionsScheme) (result *model.ContentPageScheme, response *ResponseScheme, err error) {

	if len(contentID) == 0 {
		return nil, nil, model.ErrNoContentIDError
	}

	query := url.Values{}
	query.Add("start", strconv.Itoa(startAt))
	query.Add("limit", strconv.Itoa(maxResults))

	if options != nil {

		if len(options.Expand) != 0 {
			query.Add("expand", strings.Join(options.Expand, ","))
		}

		if options.FileName != "" {
			query.Add("filename", options.FileName)
		}

		if options.MediaType != "" {
			query.Add("mediaType", options.MediaType)
		}

	}

	var endpoint = fmt.Sprintf("/rest/api/content/%v/child/attachment?%v", contentID, query.Encode())

	request, err := c.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	request.Header.Set("Accept", "application/json")

	response, err = c.client.Call(request, &result)
	if err != nil {
		return nil, response, err
	}

	return
}

// CreateOrUpdate adds an attachment to a piece of content.
// If the attachment already exists for the content,
// then the attachment is updated (i.e. a new version of the attachment is created).
// Docs: https://docs.go-atlassian.io/confluence-cloud/content/attachments#create-or-update-attachment
func (c *ContentAttachmentService) CreateOrUpdate(ctx context.Context, attachmentID, status, fileName string, file io.Reader) (
	result *model.ContentPageScheme, response *ResponseScheme, err error) {

	if len(attachmentID) == 0 {
		return nil, nil, notAttachmentIDError
	}

	if len(fileName) == 0 {
		return nil, nil, notFileNameError
	}

	if file == nil {
		return nil, nil, notReaderError
	}

	query := url.Values{}
	if len(status) != 0 {
		query.Add("status", status)
	}

	var endpoint strings.Builder
	endpoint.WriteString(fmt.Sprintf("/rest/api/content/%v/child/attachment", attachmentID))

	if query.Encode() != "" {
		endpoint.WriteString(fmt.Sprintf("?%v", query.Encode()))
	}

	var (
		body             = &bytes.Buffer{}
		attachmentWriter = multipart.NewWriter(body)
	)

	// create the attachment form row
	part, _ := attachmentWriter.CreateFormFile("file", fileName)

	// add the attachment metadata
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}

	_ = attachmentWriter.WriteField("minorEdit", "true")
	attachmentWriter.Close()

	request, err := c.client.newRequest(ctx, http.MethodPut, endpoint.String(), body)
	if err != nil {
		return nil, nil, err
	}

	request.Header.Add("Content-Type", attachmentWriter.FormDataContentType())
	request.Header.Set("X-Atlassian-Token", "no-check")

	response, err = c.client.Call(request, &result)
	if err != nil {
		return nil, response, err
	}

	return
}

// Create adds an attachment to a piece of content.
// This method only adds a new attachment.
// If you want to update an existing attachment, use Create or update attachments.
// Docs: https://docs.go-atlassian.io/confluence-cloud/content/attachments#create-attachment
func (c *ContentAttachmentService) Create(ctx context.Context, attachmentID, status, fileName string, file io.Reader) (
	result *model.ContentPageScheme, response *ResponseScheme, err error) {

	if len(attachmentID) == 0 {
		return nil, nil, notAttachmentIDError
	}

	if len(fileName) == 0 {
		return nil, nil, notFileNameError
	}

	if file == nil {
		return nil, nil, notReaderError
	}

	query := url.Values{}
	if len(status) != 0 {
		query.Add("status", status)
	}

	var endpoint strings.Builder
	endpoint.WriteString(fmt.Sprintf("/rest/api/content/%v/child/attachment", attachmentID))

	if query.Encode() != "" {
		endpoint.WriteString(fmt.Sprintf("?%v", query.Encode()))
	}

	var (
		body             = &bytes.Buffer{}
		attachmentWriter = multipart.NewWriter(body)
	)

	// create the attachment form row
	part, _ := attachmentWriter.CreateFormFile("file", fileName)

	// add the attachment metadata
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}

	_ = attachmentWriter.WriteField("minorEdit", "true")
	attachmentWriter.Close()

	request, err := c.client.newRequest(ctx, http.MethodPost, endpoint.String(), body)
	if err != nil {
		return nil, nil, err
	}

	request.Header.Add("Content-Type", attachmentWriter.FormDataContentType())
	request.Header.Set("X-Atlassian-Token", "no-check")

	response, err = c.client.Call(request, &result)
	if err != nil {
		return nil, response, err
	}

	return
}

var (
	notAttachmentIDError = fmt.Errorf("error, the attachment ID is required, please provide a valid value")
	notFileNameError     = fmt.Errorf("error, the fileName is required, please provide a valid value")
	notReaderError       = fmt.Errorf("error, the io.Reader cannot be nil, please provide a valid value")
)
