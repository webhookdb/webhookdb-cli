package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/types"
)

type IntegrationsCreateInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	ServiceName   string              `json:"service_name"`
}

func IntegrationsCreate(c context.Context, auth Auth, input IntegrationsCreateInput) (out Step, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/service_integrations/create", input.OrgIdentifier)
	return
}

type IntegrationsListInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

type IntegrationsListOutput struct {
	Data []ServiceIntegrationEntity `json:"items"`
}

func IntegrationsList(c context.Context, auth Auth, input IntegrationsListInput) (out IntegrationsListOutput, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/organizations/%v/service_integrations", input.OrgIdentifier)
	return
}

type IntegrationsResetInput struct {
	OpaqueId string `json:"-"`
}

func IntegrationsReset(c context.Context, auth Auth, input IntegrationsResetInput) (out Step, err error) {
	err = makeRequest(c, POST, auth, nil, &out, "/v1/service_integrations/%v/reset", input.OpaqueId)
	return
}
