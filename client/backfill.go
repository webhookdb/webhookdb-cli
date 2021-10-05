package client

import (
	"context"
)

type BackfillInput struct {
	OpaqueId string `json:"-"`
}

func Backfill(c context.Context, auth Auth, input BackfillInput) (step Step, err error) {
	err = makeRequest(c, POST, auth, nil, &step, "/v1/service_integrations/%v/backfill", input.OpaqueId)
	return
}

type BackfillResetInput struct {
	OpaqueId string `json:"-"`
}

func BackfillReset(c context.Context, auth Auth, input BackfillResetInput) (step Step, err error) {
	err = makeRequest(c, POST, auth, nil, &step, "/v1/service_integrations/%v/backfill/reset", input.OpaqueId)
	return
}
