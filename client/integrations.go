package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/statemachine"
)

type IntegrationsCreateInput struct {
	ServiceName string `json:"service_name"`
}

type IntegrationsCreateOutput struct {
	Step statemachine.Step `json:"step"`
}

func IntegrationsCreate(c context.Context, input IntegrationsCreateInput) (out IntegrationsCreateOutput, err error) {
	resty := RestyFromContext(c)
	resp, err := resty.R().
		SetBody(&input).
		SetError(&ErrorResponse{}).
		SetResult(&out).
		Post("/v1/integrations/create")
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}
