package client

import (
	"context"
	"github.com/webhookdb/webhookdb-cli/types"
	"net/url"
)

type OrgCloseInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func OrgClose(c context.Context, auth Auth, input OrgCloseInput) (Step, error) {
	return makeStepRequestWithResponse(c, auth, input, "/v1/organizations/%v/close", input.OrgIdentifier)
}

type OrgCreateInput struct {
	OrgName string `json:"name"`
}

type OrgCreateOutput struct {
	types.Organization
	types.MessageResponse
}

func OrgCreate(c context.Context, auth Auth, input OrgCreateInput) (out OrgCreateOutput, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/create")
	return
}

type OrgChangeRolesInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Emails        string              `json:"emails"`
	RoleName      string              `json:"role_name"`
}

func OrgChangeRoles(c context.Context, auth Auth, input OrgChangeRolesInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/change_roles", input.OrgIdentifier)
	return
}

type OrgGetInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

type OrgGetOutput struct {
	types.Organization
	types.MessageResponse
}

func OrgGet(c context.Context, auth Auth, input OrgGetInput) (out OrgGetOutput, err error) {
	values := url.Values{}
	values.Set("org", string(input.OrgIdentifier))
	err = makeRequest(c, GET, auth, values, &out, "/v1/organizations/-")
	return
}

type OrgInviteInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Email         string              `json:"email"`
	Role          string              `json:"role_name,omitempty"`
}

func OrgInvite(c context.Context, auth Auth, input OrgInviteInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/invite", input.OrgIdentifier)
	return
}

type OrgJoinInput struct {
	InvitationCode string `json:"invitation_code"`
}

func OrgJoin(c context.Context, auth Auth, input OrgJoinInput) (out types.HasOrgResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/join")
	return
}

type OrgMembersInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func OrgMembers(c context.Context, auth Auth, input OrgMembersInput) (out types.CollectionResponse, err error) {
	err = makeRequest(c, GET, auth, input, &out, "/v1/organizations/%v/members", input.OrgIdentifier)
	return
}

type OrgRemoveInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Email         string              `json:"email"`
}

func OrgRemove(c context.Context, auth Auth, input OrgRemoveInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/remove_member", input.OrgIdentifier)
	return
}

type OrgUpdateInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Field         string              `json:"field"`
	Value         string              `json:"value"`
}

func OrgUpdate(c context.Context, auth Auth, input OrgUpdateInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/update", input.OrgIdentifier)
	return
}
