package client

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/types"
)

type IntegrationsCreateInput struct {
	AuthCookie    types.AuthCookie    `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
	ServiceName   string              `json:"service_name"`
}

func IntegrationsCreate(c context.Context, input IntegrationsCreateInput) (step Step, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/organizations/%v/service_integrations/create", input.OrgIdentifier)
	resp, err := resty.R().
		SetBody(&input).
		SetError(&ErrorResponse{}).
		SetResult(&step).
		SetHeader("Cookie", string(input.AuthCookie)).
		Post(url)
	if err != nil {
		return step, err
	}
	if err := CoerceError(resp); err != nil {
		return step, err
	}
	return step, nil
}

type IntegrationsListInput struct {
	AuthCookie    types.AuthCookie
	OrgIdentifier types.OrgIdentifier
}

type IntegrationsListOutput struct {
	Data []ServiceIntegrationEntity `json:"items"`
}

func IntegrationsList(c context.Context, input IntegrationsListInput) (out IntegrationsListOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/organizations/%v/service_integrations", input.OrgIdentifier)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetResult(&out).
		SetHeader("Cookie", string(input.AuthCookie)).
		Get(url)
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}
