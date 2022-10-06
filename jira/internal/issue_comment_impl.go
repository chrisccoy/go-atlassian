package internal

import (
	model "github.com/chrisccoy/go-atlassian/pkg/infra/models"
	"github.com/chrisccoy/go-atlassian/service"
)

func NewCommentService(client service.Client, version string) (*CommentADFService, *CommentRichTextService, error) {

	if version == "" {
		return nil, nil, model.ErrNoVersionProvided
	}

	adfService := &CommentADFService{
		internalClient: &internalAdfCommentImpl{
			c:       client,
			version: version,
		},
	}

	richTextService := &CommentRichTextService{
		internalClient: &internalRichTextCommentImpl{
			c:       client,
			version: version,
		},
	}

	return adfService, richTextService, nil
}
