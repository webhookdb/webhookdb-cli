package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/types"
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
	err = makeRequest(c, GET, auth, nil, &out, "/v1/organizations/%v", input.OrgIdentifier)
	return
}

type OrgInviteInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Email         string              `json:"email"`
}

func OrgInvite(c context.Context, auth Auth, input OrgInviteInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/invite", input.OrgIdentifier)
	return
}

type OrgJoinInput struct {
	InvitationCode string `json:"invitation_code"`
}

func OrgJoin(c context.Context, auth Auth, input OrgJoinInput) (out types.MessageResponse, err error) {
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

type OrgRemoveOutput struct {
	Message string `json:"message"`
}

func OrgRemove(c context.Context, auth Auth, input OrgRemoveInput) (out OrgRemoveOutput, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/remove_member", input.OrgIdentifier)
	return
}

type OrgRenameInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Name          string              `json:"name"`
}

func OrgRename(c context.Context, auth Auth, input OrgRenameInput) (out types.MessageResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/rename", input.OrgIdentifier)
	return
}

type OrgFdwInput struct {
	OrgIdentifier    types.OrgIdentifier `json:"-"`
	MessageFdw       bool                `json:"message_fdw"`
	MessageViews     bool                `json:"message_views"`
	MessageAll       bool                `json:"message_all"`
	RemoteServerName string              `json:"remote_server_name"`
	FetchSize        string              `json:"fetch_size"`
	LocalSchema      string              `json:"local_schema"`
	ViewSchema       string              `json:"view_schema"`
}

type OrgFdwOutput map[string]interface{}

func OrgFdw(c context.Context, auth Auth, input OrgFdwInput) (out OrgFdwOutput, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/fdw", input.OrgIdentifier)
	return
}
