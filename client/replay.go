package client

import (
	"context"
	"github.com/webhookdb/webhookdb-cli/types"
)

type ReplayInput struct {
	OrgIdentifier         types.OrgIdentifier `json:"-"`
	IntegrationIdentifier string              `json:"service_integration_identifier,omitempty"`
	Hours                 int                 `json:"hours,omitempty"`
	Before                string              `json:"before,omitempty"`
	After                 string              `json:"after,omitempty"`
}

func Replay(c context.Context, auth Auth, input ReplayInput) (types.MessageResponse, error) {
	out := types.MessageResponse{}
	err := makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/replay", input.OrgIdentifier)
	return out, err
}
