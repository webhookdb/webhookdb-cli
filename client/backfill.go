package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/types"
)

type BackfillInput struct {
	OpaqueId      string              `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func Backfill(c context.Context, auth Auth, input BackfillInput) (Step, error) {
	return makeStepRequestWithResponse(c, auth, nil, "/v1/organizations/%v/service_integrations/%v/backfill", input.OrgIdentifier, input.OpaqueId)
}

type BackfillResetInput struct {
	OpaqueId      string              `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func BackfillReset(c context.Context, auth Auth, input BackfillResetInput) (Step, error) {
	return makeStepRequestWithResponse(c, auth, nil, "/v1/organizations/%v/service_integrations/%v/backfill/reset", input.OrgIdentifier, input.OpaqueId)
}
