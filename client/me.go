package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/formatting"
	"github.com/lithictech/webhookdb-cli/types"
)

type MeOrgMembershipsInput struct {
	ActiveOrgIdentifier types.OrgIdentifier `json:"-"`
}

type MeOrgMembershipsOutput struct {
	Blocks formatting.Blocks `json:"blocks"`
}

func MeOrgMemberships(c context.Context, auth Auth, input MeOrgMembershipsInput) (out MeOrgMembershipsOutput, err error) {
	err = makeRequest(c, GET, auth, input, &out, "/v1/me/organization_memberships?active_org_identifier=%s", input.ActiveOrgIdentifier)
	return
}
