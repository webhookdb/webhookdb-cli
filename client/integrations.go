package client

import (
	"context"
	"fmt"
)

type IntegrationsCreateInput struct {
	AuthCookie  string
	OrgKey      string
	ServiceName string `json:"service_name"`
}

type IntegrationsCreateOutput struct {
	Step Step `json:"step"`
}

func IntegrationsCreate(c context.Context, input IntegrationsCreateInput) (step Step, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/organizations/%v/service_integrations/create", input.OrgKey)
	resp, err := resty.R().
		SetBody(&input).
		SetError(&ErrorResponse{}).
		SetResult(&step).
		SetHeader("Cookie", input.AuthCookie).
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
	AuthCookie string
	OrgKey     string
}

type IntegrationsListOutput struct {
	Data []ServiceIntegrationEntity `json:"items"`
}

func IntegrationsList(c context.Context, input IntegrationsListInput) (out IntegrationsListOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/organizations/%v/service_integrations", input.OrgKey)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetResult(&out).
		SetHeader("Cookie", input.AuthCookie).
		Get(url)
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}
