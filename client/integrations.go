package client

import (
	"context"
	"github.com/webhookdb/webhookdb-cli/formatting"
	"github.com/webhookdb/webhookdb-cli/types"
)

type IntegrationsCreateInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	ServiceName   string              `json:"service_name"`
	GuardConfirm  *string             `json:"guard_confirm,omitempty"`
}

func IntegrationsCreate(c context.Context, auth Auth, input IntegrationsCreateInput) (out Step, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/service_integrations/create", input.OrgIdentifier)
	return
}

type IntegrationsDeleteInput struct {
	IntegrationIdentifier string              `json:"-"`
	OrgIdentifier         types.OrgIdentifier `json:"-"`
	Confirm               string              `json:"confirm"`
}

func IntegrationsDelete(c context.Context, auth Auth, input IntegrationsDeleteInput) (types.MessageResponse, error) {
	out := types.MessageResponse{}
	err := makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/service_integrations/%v/delete", input.OrgIdentifier, input.IntegrationIdentifier)
	return out, err
}

type IntegrationsInfoInput struct {
	IntegrationIdentifier string              `json:"-"`
	OrgIdentifier         types.OrgIdentifier `json:"-"`
	Field                 string              `json:"field"`
}

type IntegrationsInfoOutput struct {
	Blocks formatting.Blocks `json:"blocks"`
}

func IntegrationsInfo(c context.Context, auth Auth, input IntegrationsInfoInput) (out IntegrationsInfoOutput, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/service_integrations/%v/info", input.OrgIdentifier, input.IntegrationIdentifier)
	return out, err
}

type IntegrationsListInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func IntegrationsList(c context.Context, auth Auth, input IntegrationsListInput) (out types.CollectionResponse, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/organizations/%v/service_integrations", input.OrgIdentifier)
	return
}

type IntegrationsSetupInput struct {
	IntegrationIdentifier string              `json:"-"`
	OrgIdentifier         types.OrgIdentifier `json:"-"`
}

func IntegrationsSetup(c context.Context, auth Auth, input IntegrationsSetupInput) (Step, error) {
	return makeStepRequestWithResponse(c, auth, nil, "/v1/organizations/%v/service_integrations/%v/setup", input.OrgIdentifier, input.IntegrationIdentifier)
}

type IntegrationsResetInput struct {
	IntegrationIdentifier string              `json:"-"`
	OrgIdentifier         types.OrgIdentifier `json:"-"`
}

func IntegrationsReset(c context.Context, auth Auth, input IntegrationsResetInput) (Step, error) {
	return makeStepRequestWithResponse(c, auth, nil, "/v1/organizations/%v/service_integrations/%v/reset", input.OrgIdentifier, input.IntegrationIdentifier)
}

type IntegrationsRollKeyInput struct {
	IntegrationIdentifier string              `json:"-"`
	OrgIdentifier         types.OrgIdentifier `json:"-"`
}

func IntegrationsRollKey(c context.Context, auth Auth, input IntegrationsRollKeyInput) (out IntegrationsRollKeyOutput, err error) {
	err = makeRequest(c, POST, auth, nil, &out, "/v1/organizations/%v/service_integrations/%v/roll_api_key", input.OrgIdentifier, input.IntegrationIdentifier)
	return
}

type IntegrationsRollKeyOutput struct {
	WebhookdbApiKey string `json:"webhookdb_api_key"`
}

type IntegrationsStatsInput struct {
	IntegrationIdentifier string              `json:"-"`
	OrgIdentifier         types.OrgIdentifier `json:"-"`
}

func IntegrationsStats(c context.Context, auth Auth, input IntegrationsStatsInput) (out types.SingleResponse, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/organizations/%v/service_integrations/%v/stats", input.OrgIdentifier, input.IntegrationIdentifier)
	return
}
