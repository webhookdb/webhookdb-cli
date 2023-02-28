package client

import (
	"context"
	"github.com/lithictech/webhookdb-cli/types"
)

type SyncTargetCreateInput struct {
	OrgIdentifier         types.OrgIdentifier `json:"-"`
	IntegrationIdentifier string              `json:"service_integration_identifier"`
	ConnectionUrl         string              `json:"connection_url"`
	PageSize              int                 `json:"page_size,omitempty"`
	Period                int                 `json:"period_seconds,omitempty"`
	Schema                string              `json:"schema"`
	Table                 string              `json:"table"`
}

func SyncTargetCreate(c context.Context, auth Auth, input SyncTargetCreateInput) (out types.SingleResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/sync_targets/create", input.OrgIdentifier)
	return
}

type SyncTargetDeleteInput struct {
	OpaqueId      string              `json:"-"`
	OrgIdentifier types.OrgIdentifier `json:"-"`
	Confirm       string              `json:"confirm"`
}

type SyncTargetDeleteOutput struct {
	Message string `json:"message"`
}

func SyncTargetDelete(c context.Context, auth Auth, input SyncTargetDeleteInput) (SyncTargetDeleteOutput, error) {
	out := SyncTargetDeleteOutput{}
	err := makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/sync_targets/%v/delete", input.OrgIdentifier, input.OpaqueId)
	return out, err
}

type SyncTargetListInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
}

func SyncTargetList(c context.Context, auth Auth, input SyncTargetListInput) (out types.CollectionResponse, err error) {
	err = makeRequest(c, GET, auth, nil, &out, "/v1/organizations/%v/sync_targets", input.OrgIdentifier)
	return
}

type SyncTargetUpdateInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	OpaqueId      string              `json:"-"`
	PageSize      int                 `json:"page_size,omitempty"`
	Period        int                 `json:"period_seconds,omitempty"`
	Schema        string              `json:"schema"`
	Table         string              `json:"table"`
}

func SyncTargetUpdate(c context.Context, auth Auth, input SyncTargetUpdateInput) (out types.SingleResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/sync_targets/%v/update", input.OrgIdentifier, input.OpaqueId)
	return
}

type SyncTargetUpdateCredsInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	OpaqueId      string              `json:"-"`
	Username      string              `json:"user"`
	Password      string              `json:"password"`
}

func SyncTargetUpdateCreds(c context.Context, auth Auth, input SyncTargetUpdateCredsInput) (out types.SingleResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/sync_targets/%v/update_credentials", input.OrgIdentifier, input.OpaqueId)
	return
}

type SyncTargetSyncInput struct {
	OrgIdentifier types.OrgIdentifier `json:"-"`
	OpaqueId      string              `json:"-"`
}

func SyncTargetSync(c context.Context, auth Auth, input SyncTargetSyncInput) (out types.SingleResponse, err error) {
	err = makeRequest(c, POST, auth, input, &out, "/v1/organizations/%v/sync_targets/%v/sync", input.OrgIdentifier, input.OpaqueId)
	return
}
