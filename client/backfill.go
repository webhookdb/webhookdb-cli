package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/types"
)

type BackfillInput struct {
	OpaqueId      string              `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func Backfill(c context.Context, auth Auth, input BackfillInput) (step Step, err error) {
	err = makeRequest(c, POST, auth, nil, &step, "/v1/organizations/%v/service_integrations/%v/backfill", input.OrgIdentifier, input.OpaqueId)
	return
}

type BackfillResetInput struct {
	OpaqueId      string              `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func BackfillReset(c context.Context, auth Auth, input BackfillResetInput) (step Step, err error) {
	err = makeRequest(c, POST, auth, nil, &step, "/v1/organizations/%v/service_integrations/%v/backfill/reset", input.OrgIdentifier, input.OpaqueId)
	return
}
