package client

import (
	"context"
	"fmt"
	"github.com/lithictech/webhookdb-cli/types"
)

type OrgCreateInput struct {
	AuthCookie types.AuthCookie `json:"-"`
	OrgName    string           `json:"name"`
}

type OrgCreateOutput struct {
	Message string `json:"message"`
}

func OrgCreate(c context.Context, input OrgCreateInput) (out OrgCreateOutput, err error) {
	resty := RestyFromContext(c)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetBody(&input).
		SetResult(&out).
		SetHeader("Cookie", string(input.AuthCookie)).
		Post("/v1/organizations/create")
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}

type OrgChangeRolesInput struct {
	AuthCookie types.AuthCookie `json:"-"`
	Emails     []string `json:"emails"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
	RoleName   string `json:"role_name"`
}

type OrgChangeRolesOutput struct {
	Message string `json:"message"`
}

func OrgChangeRoles(c context.Context, input OrgChangeRolesInput) (out []OrgChangeRolesOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/organizations/%v/change_roles", input.OrgIdentifier)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetBody(&input).
		SetResult(&out).
		SetHeader("Cookie", string(input.AuthCookie)).
		Post(url)
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}

type OrgGetInput struct {
	AuthCookie    types.AuthCookie    `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

type OrgGetOutput struct {
	Org     types.Organization `json:"organization"`
	Message string             `json:"message"`
}

func OrgGet(c context.Context, input OrgGetInput) (out OrgGetOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/organizations/%v", input.OrgIdentifier)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetBody(&input).
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

type OrgInviteInput struct {
	AuthCookie    types.AuthCookie    `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Email         string              `json:"email"`
}

type OrgInviteOutput struct {
	Message string `json:"message"`
}

func OrgInvite(c context.Context, input OrgInviteInput) (out OrgInviteOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/organizations/%v/invite", input.OrgIdentifier)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetBody(&input).
		SetResult(&out).
		SetHeader("Cookie", string(input.AuthCookie)).
		Post(url)
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}

type OrgJoinInput struct {
	AuthCookie     types.AuthCookie `json:"-"`
	InvitationCode string           `json:"invitation_code"`
}

type OrgJoinOutput struct {
	Message string `json:"message"`
}

func OrgJoin(c context.Context, input OrgJoinInput) (out OrgJoinOutput, err error) {
	resty := RestyFromContext(c)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetBody(&input).
		SetResult(&out).
		SetHeader("Cookie", string(input.AuthCookie)).
		Post("/v1/organizations/join")
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}

type OrgListInput struct {
	AuthCookie types.AuthCookie `json:"-"`
}

type OrgListOutput struct {
	Items []types.Organization `json:"items"`
}

func OrgList(c context.Context, input OrgListInput) (out OrgListOutput, err error) {
	resty := RestyFromContext(c)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetResult(&out).
		SetHeader("Cookie", string(input.AuthCookie)).
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
	AuthCookie    types.AuthCookie    `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

type OrgMembersOutput struct {
	Data []OrganizationMembershipEntity `json:"items"`
}

func OrgMembers(c context.Context, input OrgMembersInput) (out OrgMembersOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/organizations/%v/members", input.OrgIdentifier)
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

type OrgRemoveInput struct {
	AuthCookie    types.AuthCookie    `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Email         string              `json:"email"`
}

type OrgRemoveOutput struct {
	Message string `json:"message"`
}

func OrgRemove(c context.Context, input OrgRemoveInput) (out OrgRemoveOutput, err error) {
	resty := RestyFromContext(c)
	url := fmt.Sprintf("/v1/organizations/%v/remove", input.OrgIdentifier)
	resp, err := resty.R().
		SetError(&ErrorResponse{}).
		SetBody(&input).
		SetResult(&out).
		SetHeader("Cookie", string(input.AuthCookie)).
		Post(url)
	if err != nil {
		return out, err
	}
	if err := CoerceError(resp); err != nil {
		return out, err
	}
	return out, nil
}
