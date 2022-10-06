package jira

import (
	"context"
	model "github.com/chrisccoy/go-atlassian/pkg/infra/models"
)

type CommentRichTextConnector interface {
	CommentSharedConnector

	// Gets returns all comments for an issue.
	//
	// GET /rest/api/{2-3}/issue/{issueIdOrKey}/comment
	//
	// https://docs.go-atlassian.io/jira-software-cloud/issues/comments#get-comments
	Gets(ctx context.Context, issueKeyOrId, orderBy string, expand []string, startAt, maxResults int) (*model.IssueCommentPageSchemeV2, *model.ResponseScheme, error)

	// Get returns a comment.
	//
	// GET /rest/api/{2-3}/issue/{issueIdOrKey}/comment/{id}
	//
	// TODO: The documentation needs to be created, raise a ticket here: https://github.com/chrisccoy/go-atlassian/issues
	Get(ctx context.Context, issueKeyOrId, commentId string) (*model.IssueCommentSchemeV2, *model.ResponseScheme, error)

	// Add adds a comment to an issue.
	//
	// POST /rest/api/{2-3}/issue/{issueIdOrKey}/comment
	//
	//https://docs.go-atlassian.io/jira-software-cloud/issues/comments#add-comment
	Add(ctx context.Context, issueKeyOrId string, payload *model.CommentPayloadSchemeV2, expand []string) (*model.IssueCommentSchemeV2, *model.ResponseScheme, error)
}

type CommentADFConnector interface {
	CommentSharedConnector

	// Gets returns all comments for an issue.
	//
	// GET /rest/api/{2-3}/issue/{issueIdOrKey}/comment
	//
	// https://docs.go-atlassian.io/jira-software-cloud/issues/comments#get-comments
	Gets(ctx context.Context, issueKeyOrId, orderBy string, expand []string, startAt, maxResults int) (*model.IssueCommentPageScheme, *model.ResponseScheme, error)

	// Get returns a comment.
	//
	// GET /rest/api/{2-3}/issue/{issueIdOrKey}/comment/{id}
	//
	// TODO: The documentation needs to be created, raise a ticket here: https://github.com/chrisccoy/go-atlassian/issues
	Get(ctx context.Context, issueKeyOrId, commentId string) (*model.IssueCommentScheme, *model.ResponseScheme, error)

	// Add adds a comment to an issue.
	//
	// POST /rest/api/{2-3}/issue/{issueIdOrKey}/comment
	//
	// https://docs.go-atlassian.io/jira-software-cloud/issues/comments#add-comment
	Add(ctx context.Context, issueKeyOrId string, payload *model.CommentPayloadScheme, expand []string) (*model.IssueCommentScheme, *model.ResponseScheme, error)
}

type CommentSharedConnector interface {

	// Delete deletes a comment.
	//
	// DELETE /rest/api/{2-3}/issue/{issueIdOrKey}/comment/{id}
	//
	// https://docs.go-atlassian.io/jira-software-cloud/issues/comments#delete-comment
	Delete(ctx context.Context, issueKeyOrId, commentId string) (*model.ResponseScheme, error)
}
