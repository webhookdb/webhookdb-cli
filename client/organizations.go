package client

import (
	"context"
)

type OrganizationsListInput struct {
	AuthCookie string
}

type OrganizationsListOutput struct {
	Data []OrganizationEntity `json:"items"`
}

func OrganizationsList(c context.Context, input OrganizationsListInput) (out OrganizationsListOutput, err error) {
	resty := RestyFromContext(c)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetResult(&out).
		SetHeader("Cookie", input.AuthCookie).
		Get("/v1/organizations/")
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}
