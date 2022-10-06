package internal

import (
	"github.com/chrisccoy/go-atlassian/service"
	"github.com/chrisccoy/go-atlassian/service/common"
)

func NewAuthenticationService(client service.Client) common.Authentication {
	return &AuthenticationService{c: client}
}

type AuthenticationService struct {
	c service.Client

	basicAuthProvided bool
	mail, token       string

	userAgentProvided bool
	agent             string

	experimentalProvided bool
}

func (a *AuthenticationService) SetExperimentalFlag() {
	a.experimentalProvided = true
}

func (a *AuthenticationService) HasSetExperimentalFlag() bool {
	return a.experimentalProvided
}

func (a *AuthenticationService) SetBasicAuth(mail, token string) {
	a.mail = mail
	a.token = token

	a.basicAuthProvided = true
}

func (a *AuthenticationService) GetBasicAuth() (string, string) {
	return a.mail, a.token
}

func (a *AuthenticationService) HasBasicAuth() bool {
	return a.basicAuthProvided
}

func (a *AuthenticationService) SetUserAgent(agent string) {
	a.agent = agent
	a.userAgentProvided = true
}

func (a *AuthenticationService) GetUserAgent() string {
	return a.agent
}

func (a *AuthenticationService) HasUserAgent() bool {
	return a.userAgentProvided
}
