package internal

import (
	model "github.com/chrisccoy/go-atlassian/pkg/infra/models"
	"github.com/chrisccoy/go-atlassian/service"
)

func NewLinkService(client service.Client, version string, type_ *LinkTypeService) (*LinkADFService, *LinkRichTextService, error) {

	if version == "" {
		return nil, nil, model.ErrNoVersionProvided
	}

	adfService := &LinkADFService{
		internalClient: &internalLinkADFServiceImpl{
			c:       client,
			version: version,
		},
		Type: type_,
	}

	richTextService := &LinkRichTextService{
		internalClient: &internalLinkRichTextServiceImpl{
			c:       client,
			version: version,
		},
		Type: type_,
	}

	return adfService, richTextService, nil
}
