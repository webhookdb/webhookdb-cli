package client

import (
	"context"
	"fmt"
)

type OrgListInput struct {
	AuthCookie string
}

type OrgListOutput struct {
	Data []OrganizationEntity `json:"items"`
}

func OrgList(c context.Context, input OrgListInput) (out OrgListOutput, err error) {
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

type OrgMembersInput struct {
	AuthCookie string
	OrgKey string
}

type OrgMembersOutput struct {
	Data []OrganizationMembershipEntity `json:"items"`
}

func OrgMembers(c context.Context, input OrgMembersInput) (out OrgMembersOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/organizations/%v/members", input.OrgKey)
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